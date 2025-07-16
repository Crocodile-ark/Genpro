package keeper

import (
	"fmt"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Crocodile-ark/gxrchaind/x/halving/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper authkeeper.AccountKeeper
		bankKeeper    bankkeeper.Keeper
		stakingKeeper *stakingkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetHalvingInfo gets the current halving information
func (k Keeper) GetHalvingInfo(ctx sdk.Context) (types.HalvingInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CurrentHalvingKey)
	if bz == nil {
		return types.HalvingInfo{}, false
	}

	var info types.HalvingInfo
	k.cdc.MustUnmarshal(bz, &info)
	return info, true
}

// SetHalvingInfo sets the current halving information
func (k Keeper) SetHalvingInfo(ctx sdk.Context, info types.HalvingInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&info)
	store.Set(types.CurrentHalvingKey, bz)
}

// GetDistributionRecord gets a specific distribution record
func (k Keeper) GetDistributionRecord(ctx sdk.Context, timestamp int64) (types.DistributionRecord, bool) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.LastDistributionKey, sdk.Uint64ToBigEndian(uint64(timestamp))...)
	bz := store.Get(key)
	if bz == nil {
		return types.DistributionRecord{}, false
	}

	var record types.DistributionRecord
	k.cdc.MustUnmarshal(bz, &record)
	return record, true
}

// SetDistributionRecord sets a distribution record
func (k Keeper) SetDistributionRecord(ctx sdk.Context, record types.DistributionRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.LastDistributionKey, sdk.Uint64ToBigEndian(uint64(record.Timestamp))...)
	bz := k.cdc.MustMarshal(&record)
	store.Set(key, bz)
}

// GetAllDistributionRecords gets all distribution records
func (k Keeper) GetAllDistributionRecords(ctx sdk.Context) []types.DistributionRecord {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LastDistributionKey)
	defer iterator.Close()

	var records []types.DistributionRecord
	for ; iterator.Valid(); iterator.Next() {
		var record types.DistributionRecord
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		records = append(records, record)
	}

	return records
}

// CalculateMonthlyReward calculates the monthly reward amount for current cycle
func (k Keeper) CalculateMonthlyReward(ctx sdk.Context) sdk.Coin {
	info, found := k.GetHalvingInfo(ctx)
	if !found {
		return sdk.NewCoin("ugen", sdk.ZeroInt())
	}

	// Each cycle is 5 years = 60 months
	// Distribute total funds evenly over 60 months
	monthlyAmount := info.TotalFundsForCycle.Amount.QuoRaw(60)
	return sdk.NewCoin(info.TotalFundsForCycle.Denom, monthlyAmount)
}

// DistributeMonthlyRewards distributes monthly rewards according to GXR specification
func (k Keeper) DistributeMonthlyRewards(ctx sdk.Context) error {
	params := k.GetParams(ctx)
	monthlyReward := k.CalculateMonthlyReward(ctx)

	if monthlyReward.IsZero() {
		return fmt.Errorf("no monthly rewards to distribute")
	}

	// Calculate distribution amounts
	validatorAmount := monthlyReward.Amount.ToDec().Mul(params.ValidatorShare).TruncateInt()
	delegatorAmount := monthlyReward.Amount.ToDec().Mul(params.DelegatorShare).TruncateInt()
	dexAmount := monthlyReward.Amount.ToDec().Mul(params.DexShare).TruncateInt()

	// Get module account
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	// Distribute to validators (70%)
	if err := k.distributeToValidators(ctx, sdk.NewCoin(monthlyReward.Denom, validatorAmount)); err != nil {
		return fmt.Errorf("failed to distribute to validators: %w", err)
	}

	// Distribute to delegators via PoS pool (20%)
	if err := k.distributeToDelegators(ctx, sdk.NewCoin(monthlyReward.Denom, delegatorAmount)); err != nil {
		return fmt.Errorf("failed to distribute to delegators: %w", err)
	}

	// Distribute to DEX pool (10%)
	if err := k.distributeToDEX(ctx, sdk.NewCoin(monthlyReward.Denom, dexAmount)); err != nil {
		return fmt.Errorf("failed to distribute to DEX: %w", err)
	}

	// Update halving info
	info, _ := k.GetHalvingInfo(ctx)
	info.DistributedInCycle = info.DistributedInCycle.Add(monthlyReward)
	k.SetHalvingInfo(ctx, info)

	// Record distribution
	record := types.DistributionRecord{
		Timestamp: ctx.BlockTime().Unix(),
		Amount:    monthlyReward,
		Cycle:     info.CurrentCycle,
		Month:     k.getCurrentMonth(ctx, info),
	}
	k.SetDistributionRecord(ctx, record)

	k.Logger(ctx).Info("Monthly rewards distributed",
		"amount", monthlyReward.String(),
		"cycle", info.CurrentCycle,
		"month", record.Month,
	)

	return nil
}

// distributeToValidators distributes rewards to active validators
func (k Keeper) distributeToValidators(ctx sdk.Context, amount sdk.Coin) error {
	// Get all bonded validators
	validators := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validators) == 0 {
		return fmt.Errorf("no bonded validators found")
	}

	// Distribute equally among active validators
	perValidatorAmount := amount.Amount.QuoRaw(int64(len(validators)))
	if perValidatorAmount.IsZero() {
		return nil
	}

	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	
	for _, validator := range validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			continue
		}

		accAddr := sdk.AccAddress(valAddr)
		reward := sdk.NewCoin(amount.Denom, perValidatorAmount)
		
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddr, sdk.NewCoins(reward)); err != nil {
			k.Logger(ctx).Error("Failed to send reward to validator", "validator", validator.OperatorAddress, "error", err)
			continue
		}
	}

	return nil
}

// distributeToDelegators distributes rewards to the PoS pool for delegators
func (k Keeper) distributeToDelegators(ctx sdk.Context, amount sdk.Coin) error {
	// Send to distribution module's fee pool for delegators
	feeCollectorAddr := k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, feeCollectorAddr, sdk.NewCoins(amount))
}

// distributeToDEX distributes rewards to DEX pool
func (k Keeper) distributeToDEX(ctx sdk.Context, amount sdk.Coin) error {
	// For now, keep in module account - will be handled by fee router
	// In production, this would be sent to specific DEX pool addresses
	k.Logger(ctx).Info("DEX rewards allocated", "amount", amount.String())
	return nil
}

// getCurrentMonth calculates current month within the cycle (1-60)
func (k Keeper) getCurrentMonth(ctx sdk.Context, info types.HalvingInfo) uint64 {
	cycleStart := time.Unix(info.CycleStartTime, 0)
	elapsed := ctx.BlockTime().Sub(cycleStart)
	monthsElapsed := int64(elapsed.Hours() / (24 * 30)) // Approximate month as 30 days
	return uint64(monthsElapsed + 1) // Month starts from 1
}

// CheckAndAdvanceCycle checks if current cycle should advance to next
func (k Keeper) CheckAndAdvanceCycle(ctx sdk.Context) error {
	info, found := k.GetHalvingInfo(ctx)
	if !found {
		return fmt.Errorf("halving info not found")
	}

	params := k.GetParams(ctx)
	cycleStart := time.Unix(info.CycleStartTime, 0)
	
	// Check if current cycle duration has passed
	if ctx.BlockTime().Sub(cycleStart) >= params.HalvingCycleDuration {
		// Advance to next cycle
		if info.CurrentCycle >= 5 {
			k.Logger(ctx).Info("All halving cycles completed")
			return nil
		}

		nextCycle := info.CurrentCycle + 1
		nextCycleFunds := k.calculateNextCycleFunds(info.TotalFundsForCycle)

		newInfo := types.HalvingInfo{
			CurrentCycle:       nextCycle,
			CycleStartTime:     ctx.BlockTime().Unix(),
			TotalFundsForCycle: nextCycleFunds,
			DistributedInCycle: sdk.NewCoin("ugen", sdk.ZeroInt()),
		}

		k.SetHalvingInfo(ctx, newInfo)
		
		k.Logger(ctx).Info("Advanced to next halving cycle",
			"new_cycle", nextCycle,
			"new_funds", nextCycleFunds.String(),
		)
	}

	return nil
}

// calculateNextCycleFunds calculates funds for next cycle with 15% reduction
func (k Keeper) calculateNextCycleFunds(currentFunds sdk.Coin) sdk.Coin {
	// Reduce by 15% each cycle according to GXR specification
	nextAmount := currentFunds.Amount.ToDec().Mul(sdk.MustNewDecFromStr("0.85")).TruncateInt()
	return sdk.NewCoin(currentFunds.Denom, nextAmount)
}