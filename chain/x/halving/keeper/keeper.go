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

const (
	// MinimumSupplyThreshold is the minimum supply threshold (1,000 GXR)
	MinimumSupplyThreshold = 1000 * 1e8 // 1,000 GXR in ugen
	// HalvingReductionRate is the reduction rate per halving cycle (15%)
	HalvingReductionRate = "0.15"
	// MainDenom is the main denomination
	MainDenom = "ugen"
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

// GetCurrentTotalSupply gets the current total supply of GXR
func (k Keeper) GetCurrentTotalSupply(ctx sdk.Context) sdk.Coin {
	supply := k.bankKeeper.GetSupply(ctx, MainDenom)
	return supply
}

// CalculateMonthlyReward calculates the monthly reward amount based on current total supply
func (k Keeper) CalculateMonthlyReward(ctx sdk.Context) sdk.Coin {
	currentSupply := k.GetCurrentTotalSupply(ctx)
	
	// Calculate 15% of current supply for the 5-year cycle
	reductionRate := sdk.MustNewDecFromStr(HalvingReductionRate)
	cycleReduction := currentSupply.Amount.ToDec().Mul(reductionRate).TruncateInt()
	
	// Distribute the reduction over 60 months (5 years)
	monthlyAmount := cycleReduction.QuoRaw(60)
	
	return sdk.NewCoin(MainDenom, monthlyAmount)
}

// DistributeMonthlyRewards distributes monthly rewards according to GXR specification
func (k Keeper) DistributeMonthlyRewards(ctx sdk.Context) error {
	// Check if total supply is below threshold
	currentSupply := k.GetCurrentTotalSupply(ctx)
	if currentSupply.Amount.LT(sdk.NewInt(MinimumSupplyThreshold)) {
		k.Logger(ctx).Info("Halving stopped: total supply below minimum threshold",
			"current_supply", currentSupply.String(),
			"threshold", fmt.Sprintf("%dugen", MinimumSupplyThreshold))
		return nil
	}

	params := k.GetParams(ctx)
	monthlyReward := k.CalculateMonthlyReward(ctx)

	if monthlyReward.IsZero() {
		return fmt.Errorf("no monthly rewards to distribute")
	}

	// Calculate distribution amounts
	validatorAmount := monthlyReward.Amount.ToDec().Mul(params.ValidatorShare).TruncateInt()
	delegatorAmount := monthlyReward.Amount.ToDec().Mul(params.DelegatorShare).TruncateInt()
	dexAmount := monthlyReward.Amount.ToDec().Mul(params.DexShare).TruncateInt()

	// Burn the monthly reward from total supply (this is the key change)
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(monthlyReward)); err != nil {
		return fmt.Errorf("failed to burn monthly reward: %w", err)
	}

	// Get module account
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	if moduleAddr == nil {
		return fmt.Errorf("module account not found")
	}

	// Create new coins for distribution (minted specifically for rewards)
	rewardCoins := sdk.NewCoins(
		sdk.NewCoin(MainDenom, validatorAmount),
		sdk.NewCoin(MainDenom, delegatorAmount),
		sdk.NewCoin(MainDenom, dexAmount),
	)

	// Mint the distribution amounts to module account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, rewardCoins); err != nil {
		return fmt.Errorf("failed to mint reward coins: %w", err)
	}

	// Distribute to validators (70%)
	if err := k.distributeToValidators(ctx, sdk.NewCoin(MainDenom, validatorAmount)); err != nil {
		return fmt.Errorf("failed to distribute to validators: %w", err)
	}

	// Distribute to delegators via PoS pool (20%)
	if err := k.distributeToDelegators(ctx, sdk.NewCoin(MainDenom, delegatorAmount)); err != nil {
		return fmt.Errorf("failed to distribute to delegators: %w", err)
	}

	// Distribute to DEX pool (10%)
	if err := k.distributeToDEX(ctx, sdk.NewCoin(MainDenom, dexAmount)); err != nil {
		return fmt.Errorf("failed to distribute to DEX: %w", err)
	}

	// Update halving info
	info, found := k.GetHalvingInfo(ctx)
	if !found {
		// Initialize halving info if not found
		info = types.HalvingInfo{
			CurrentCycle:       1,
			CycleStartTime:     ctx.BlockTime().Unix(),
			TotalFundsForCycle: currentSupply,
			DistributedInCycle: sdk.NewCoin(MainDenom, sdk.ZeroInt()),
		}
	}
	
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
		"burned_amount", monthlyReward.String(),
		"distributed_amount", rewardCoins.String(),
		"new_total_supply", k.GetCurrentTotalSupply(ctx).String(),
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

	for _, validator := range validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			k.Logger(ctx).Error("Invalid validator address", "validator", validator.OperatorAddress, "error", err)
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
	if feeCollectorAddr == nil {
		return fmt.Errorf("fee collector account not found")
	}
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
		// Initialize first cycle if not found
		currentSupply := k.GetCurrentTotalSupply(ctx)
		info = types.HalvingInfo{
			CurrentCycle:       1,
			CycleStartTime:     ctx.BlockTime().Unix(),
			TotalFundsForCycle: currentSupply,
			DistributedInCycle: sdk.NewCoin(MainDenom, sdk.ZeroInt()),
		}
		k.SetHalvingInfo(ctx, info)
		return nil
	}

	// Check if total supply is below threshold
	currentSupply := k.GetCurrentTotalSupply(ctx)
	if currentSupply.Amount.LT(sdk.NewInt(MinimumSupplyThreshold)) {
		k.Logger(ctx).Info("Halving cycles stopped: total supply below minimum threshold",
			"current_supply", currentSupply.String(),
			"threshold", fmt.Sprintf("%dugen", MinimumSupplyThreshold))
		return nil
	}

	params := k.GetParams(ctx)
	cycleStart := time.Unix(info.CycleStartTime, 0)
	
	// Check if current cycle duration has passed (5 years)
	if ctx.BlockTime().Sub(cycleStart) >= params.HalvingCycleDuration {
		// Advance to next cycle
		nextCycle := info.CurrentCycle + 1
		
		// The new cycle starts with the current total supply as 100%
		// (This is automatic since we work with actual supply)
		newInfo := types.HalvingInfo{
			CurrentCycle:       nextCycle,
			CycleStartTime:     ctx.BlockTime().Unix(),
			TotalFundsForCycle: currentSupply, // Current supply becomes 100% for next cycle
			DistributedInCycle: sdk.NewCoin(MainDenom, sdk.ZeroInt()),
		}

		k.SetHalvingInfo(ctx, newInfo)
		
		k.Logger(ctx).Info("Advanced to next halving cycle",
			"new_cycle", nextCycle,
			"current_supply", currentSupply.String(),
		)
	}

	return nil
}

// calculateNextCycleFunds is no longer used since we work with actual supply
// Keeping for backwards compatibility
func (k Keeper) calculateNextCycleFunds(currentFunds sdk.Coin) sdk.Coin {
	// This function is deprecated in the new logic
	// We now work directly with the current total supply
	return currentFunds
}