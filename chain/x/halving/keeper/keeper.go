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
	// HalvingCycleDuration is 5 years
	HalvingCycleDuration = 5 * 365 * 24 * time.Hour
	// DistributionPeriod is 2 years (730 days)
	DistributionPeriod = 730 * 24 * time.Hour
	// PausePeriod is 3 years after distribution
	PausePeriod = 3 * 365 * 24 * time.Hour
	// ValidatorInactiveThreshold is 10 days per month
	ValidatorInactiveThreshold = 10
	// MonthDuration is 30 days
	MonthDuration = 30 * 24 * time.Hour
	// DEXDistributionPeriod is 2 years (only years 1-2)
	DEXDistributionPeriod = 2 * 365 * 24 * time.Hour
	// MonthlyDistributionTrigger is 30 days
	MonthlyDistributionTrigger = 30 * 24 * time.Hour
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

// GetCurrentTotalSupply gets the current total supply of GXR
func (k Keeper) GetCurrentTotalSupply(ctx sdk.Context) sdk.Coin {
	supply := k.bankKeeper.GetSupply(ctx, MainDenom)
	return supply
}

// GetValidatorUptime gets validator uptime record
func (k Keeper) GetValidatorUptime(ctx sdk.Context, valAddr sdk.ValAddress) (types.ValidatorUptime, bool) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.ValidatorUptimeKey, valAddr.Bytes()...)
	bz := store.Get(key)
	if bz == nil {
		return types.ValidatorUptime{}, false
	}

	var uptime types.ValidatorUptime
	k.cdc.MustUnmarshal(bz, &uptime)
	return uptime, true
}

// SetValidatorUptime sets validator uptime record
func (k Keeper) SetValidatorUptime(ctx sdk.Context, valAddr sdk.ValAddress, uptime types.ValidatorUptime) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.ValidatorUptimeKey, valAddr.Bytes()...)
	bz := k.cdc.MustMarshal(&uptime)
	store.Set(key, bz)
}

// GetLastDistributionTime gets the last distribution timestamp
func (k Keeper) GetLastDistributionTime(ctx sdk.Context) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LastDistributionKey)
	if bz == nil {
		return 0, false
	}
	
	return sdk.BigEndianToUint64(bz), true
}

// SetLastDistributionTime sets the last distribution timestamp
func (k Keeper) SetLastDistributionTime(ctx sdk.Context, timestamp int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastDistributionKey, sdk.Uint64ToBigEndian(uint64(timestamp)))
}

// CheckAndAdvanceHalvingCycle checks if we should advance to the next halving cycle
func (k Keeper) CheckAndAdvanceHalvingCycle(ctx sdk.Context) error {
	info, found := k.GetHalvingInfo(ctx)
	if !found {
		// Initialize first cycle
		currentSupply := k.GetCurrentTotalSupply(ctx)
		info = types.HalvingInfo{
			CurrentCycle:       1,
			CycleStartTime:     ctx.BlockTime().Unix(),
			TotalSupply:        currentSupply,
			HalvingFund:        sdk.NewCoin(MainDenom, sdk.ZeroInt()),
			DistributionActive: false,
			DistributionStart:  0,
			DistributedAmount:  sdk.NewCoin(MainDenom, sdk.ZeroInt()),
			PauseStart:         0,
			LastMonthlyDistrib: 0,
		}
		k.SetHalvingInfo(ctx, info)
		k.Logger(ctx).Info("Initialized first halving cycle", "cycle", 1, "total_supply", currentSupply.String())
		return nil
	}

	// Check if total supply is below threshold - stop permanently
	currentSupply := k.GetCurrentTotalSupply(ctx)
	if currentSupply.Amount.LT(sdk.NewInt(MinimumSupplyThreshold)) {
		k.Logger(ctx).Info("Halving stopped permanently: total supply below minimum threshold",
			"current_supply", currentSupply.String(),
			"threshold", fmt.Sprintf("%dugen", MinimumSupplyThreshold))
		return nil
	}

	cycleStart := time.Unix(info.CycleStartTime, 0)
	
	// Check if 5 years have passed since cycle start (based on ctx.BlockTime())
	if ctx.BlockTime().Sub(cycleStart) >= HalvingCycleDuration {
		// Advance to next cycle
		return k.advanceToNextCycle(ctx, info)
	}

	return nil
}

// advanceToNextCycle advances to the next halving cycle
func (k Keeper) advanceToNextCycle(ctx sdk.Context, info types.HalvingInfo) error {
	currentSupply := k.GetCurrentTotalSupply(ctx)
	
	// Calculate 15% for halving fund
	reductionRate := sdk.MustNewDecFromStr(HalvingReductionRate)
	halvingAmount := currentSupply.Amount.ToDec().Mul(reductionRate).TruncateInt()
	
	// Create halving fund entry (virtual allocation)
	halvingFund := sdk.NewCoin(MainDenom, halvingAmount)
	
	// Update halving info for next cycle
	newInfo := types.HalvingInfo{
		CurrentCycle:       info.CurrentCycle + 1,
		CycleStartTime:     ctx.BlockTime().Unix(),
		TotalSupply:        currentSupply,
		HalvingFund:        halvingFund,
		DistributionActive: true,
		DistributionStart:  ctx.BlockTime().Unix(),
		DistributedAmount:  sdk.NewCoin(MainDenom, sdk.ZeroInt()),
		PauseStart:         0,
		LastMonthlyDistrib: 0,
	}

	k.SetHalvingInfo(ctx, newInfo)
	
	k.Logger(ctx).Info("Advanced to next halving cycle",
		"new_cycle", newInfo.CurrentCycle,
		"halving_fund", halvingFund.String(),
		"current_supply", currentSupply.String(),
		"distribution_start", ctx.BlockTime().Unix(),
	)

	return nil
}

// CheckAndUpdateDistributionStatus checks and updates distribution status based on timing
func (k Keeper) CheckAndUpdateDistributionStatus(ctx sdk.Context) error {
	info, found := k.GetHalvingInfo(ctx)
	if !found {
		return nil
	}

	// If distribution is active, check if 2 years have passed
	if info.DistributionActive {
		distributionStart := time.Unix(info.DistributionStart, 0)
		if ctx.BlockTime().Sub(distributionStart) >= DistributionPeriod {
			// Stop distribution and start 3-year pause
			info.DistributionActive = false
			info.PauseStart = ctx.BlockTime().Unix()
			k.SetHalvingInfo(ctx, info)
			
			k.Logger(ctx).Info("Distribution period ended, entering 3-year pause",
				"cycle", info.CurrentCycle,
				"distributed_amount", info.DistributedAmount.String(),
				"pause_start", info.PauseStart,
			)
		}
	}

	return nil
}

// ShouldDistribute checks if monthly distribution should occur
func (k Keeper) ShouldDistribute(ctx sdk.Context) bool {
	info, found := k.GetHalvingInfo(ctx)
	if !found || !info.DistributionActive {
		return false
	}

	// Check if 30 days have passed since last distribution
	if info.LastMonthlyDistrib == 0 {
		return true // First distribution
	}

	lastDistrib := time.Unix(info.LastMonthlyDistrib, 0)
	return ctx.BlockTime().Sub(lastDistrib) >= MonthlyDistributionTrigger
}

// DistributeHalvingRewards distributes monthly rewards from halving fund
func (k Keeper) DistributeHalvingRewards(ctx sdk.Context) error {
	info, found := k.GetHalvingInfo(ctx)
	if !found || !info.DistributionActive {
		return nil
	}

	if !k.ShouldDistribute(ctx) {
		return nil
	}

	// Calculate monthly distribution amount (over 24 months)
	monthlyAmount := k.calculateMonthlyDistribution(ctx, info)
	if monthlyAmount.IsZero() {
		return nil
	}

	// Check if we have enough in halving fund
	if info.HalvingFund.Amount.LT(monthlyAmount.Amount) {
		monthlyAmount = info.HalvingFund
	}

	// Burn the monthly amount from total supply
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(monthlyAmount)); err != nil {
		return fmt.Errorf("failed to burn monthly distribution: %w", err)
	}

	// Mint the same amount for distribution
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(monthlyAmount)); err != nil {
		return fmt.Errorf("failed to mint distribution coins: %w", err)
	}

	// Distribute rewards
	if err := k.distributeRewards(ctx, monthlyAmount, info); err != nil {
		return fmt.Errorf("failed to distribute rewards: %w", err)
	}

	// Update halving info
	info.DistributedAmount = info.DistributedAmount.Add(monthlyAmount)
	info.HalvingFund = info.HalvingFund.Sub(monthlyAmount)
	info.LastMonthlyDistrib = ctx.BlockTime().Unix()
	k.SetHalvingInfo(ctx, info)

	k.Logger(ctx).Info("Monthly halving rewards distributed",
		"amount", monthlyAmount.String(),
		"cycle", info.CurrentCycle,
		"remaining_fund", info.HalvingFund.String(),
		"total_distributed", info.DistributedAmount.String(),
	)

	return nil
}

// calculateMonthlyDistribution calculates monthly distribution amount
func (k Keeper) calculateMonthlyDistribution(ctx sdk.Context, info types.HalvingInfo) sdk.Coin {
	// Distribute over 24 months (2 years)
	totalMonths := int64(24)
	monthlyAmount := info.HalvingFund.Amount.QuoRaw(totalMonths)
	
	// Ensure we don't exceed available funds
	if monthlyAmount.GT(info.HalvingFund.Amount) {
		monthlyAmount = info.HalvingFund.Amount
	}

	return sdk.NewCoin(MainDenom, monthlyAmount)
}

// distributeRewards distributes rewards according to the enhanced specifications
func (k Keeper) distributeRewards(ctx sdk.Context, totalAmount sdk.Coin, info types.HalvingInfo) error {
	// Distribution percentages:
	// - 70% to active validators
	// - 20% to delegators (PoS staking pool)
	// - 10% to DEX pools (only years 1-2)
	
	validatorAmount := totalAmount.Amount.ToDec().Mul(sdk.MustNewDecFromStr("0.70")).TruncateInt()
	delegatorAmount := totalAmount.Amount.ToDec().Mul(sdk.MustNewDecFromStr("0.20")).TruncateInt()
	dexAmount := totalAmount.Amount.ToDec().Mul(sdk.MustNewDecFromStr("0.10")).TruncateInt()

	// Distribute to active validators (70%)
	if err := k.distributeToActiveValidators(ctx, sdk.NewCoin(MainDenom, validatorAmount)); err != nil {
		return fmt.Errorf("failed to distribute to validators: %w", err)
	}

	// Distribute to delegators (20%)
	if err := k.distributeToDelegators(ctx, sdk.NewCoin(MainDenom, delegatorAmount)); err != nil {
		return fmt.Errorf("failed to distribute to delegators: %w", err)
	}

	// Distribute to DEX (10%, only in years 1-2)
	if err := k.distributeToDEX(ctx, sdk.NewCoin(MainDenom, dexAmount), info); err != nil {
		return fmt.Errorf("failed to distribute to DEX: %w", err)
	}

	return nil
}

// distributeToActiveValidators distributes rewards to active validators only
func (k Keeper) distributeToActiveValidators(ctx sdk.Context, amount sdk.Coin) error {
	validators := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validators) == 0 {
		k.Logger(ctx).Info("No bonded validators found, forfeiting validator rewards")
		return nil
	}

	// Filter active validators (uptime > 20 days in current month)
	activeValidators := make([]stakingtypes.Validator, 0)
	for _, validator := range validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			k.Logger(ctx).Error("Invalid validator address", "validator", validator.OperatorAddress, "error", err)
			continue
		}

		if k.isValidatorActive(ctx, valAddr) {
			activeValidators = append(activeValidators, validator)
		} else {
			k.Logger(ctx).Info("Validator forfeit rewards due to inactivity",
				"validator", validator.OperatorAddress,
				"month", k.getCurrentMonth(ctx),
			)
		}
	}

	if len(activeValidators) == 0 {
		k.Logger(ctx).Info("No active validators found, forfeiting all validator rewards")
		return nil
	}

	// Distribute equally among active validators
	perValidatorAmount := amount.Amount.QuoRaw(int64(len(activeValidators)))
	if perValidatorAmount.IsZero() {
		return nil
	}

	for _, validator := range activeValidators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			continue
		}

		accAddr := sdk.AccAddress(valAddr)
		reward := sdk.NewCoin(MainDenom, perValidatorAmount)
		
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accAddr, sdk.NewCoins(reward)); err != nil {
			k.Logger(ctx).Error("Failed to send reward to validator", "validator", validator.OperatorAddress, "error", err)
			continue
		}

		k.Logger(ctx).Info("Distributed reward to active validator",
			"validator", validator.OperatorAddress,
			"amount", reward.String(),
		)
	}

	return nil
}

// isValidatorActive checks if validator is active (not inactive >10 days in current month)
func (k Keeper) isValidatorActive(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	uptime, found := k.GetValidatorUptime(ctx, valAddr)
	if !found {
		// First month, initialize uptime record
		uptime = types.ValidatorUptime{
			ValidatorAddress: valAddr.String(),
			CurrentMonth:     k.getCurrentMonth(ctx),
			InactiveDays:     0,
			LastCheck:        ctx.BlockTime().Unix(),
		}
		k.SetValidatorUptime(ctx, valAddr, uptime)
		return true
	}

	currentMonth := k.getCurrentMonth(ctx)
	if uptime.CurrentMonth != currentMonth {
		// New month, reset counters
		uptime.CurrentMonth = currentMonth
		uptime.InactiveDays = 0
		uptime.LastCheck = ctx.BlockTime().Unix()
		k.SetValidatorUptime(ctx, valAddr, uptime)
		return true
	}

	// Check validator status
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		k.Logger(ctx).Error("Validator not found", "validator", valAddr.String())
		return false
	}

	// Update inactive days if validator is not bonded
	if validator.Status != stakingtypes.Bonded {
		lastCheck := time.Unix(uptime.LastCheck, 0)
		if ctx.BlockTime().Sub(lastCheck) >= 24*time.Hour {
			uptime.InactiveDays++
			uptime.LastCheck = ctx.BlockTime().Unix()
			k.SetValidatorUptime(ctx, valAddr, uptime)
		}
	}

	// Validator is active if inactive days <= 10
	return uptime.InactiveDays <= ValidatorInactiveThreshold
}

// getCurrentMonth returns current month identifier
func (k Keeper) getCurrentMonth(ctx sdk.Context) uint64 {
	return uint64(ctx.BlockTime().Unix() / int64(MonthDuration.Seconds()))
}

// distributeToDelegators distributes rewards to delegators via fee pool
func (k Keeper) distributeToDelegators(ctx sdk.Context, amount sdk.Coin) error {
	// Send to fee collector for distribution to delegators
	feeCollectorAddr := k.accountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	if feeCollectorAddr == nil {
		return fmt.Errorf("fee collector account not found")
	}
	
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, feeCollectorAddr, sdk.NewCoins(amount)); err != nil {
		return fmt.Errorf("failed to send to fee collector: %w", err)
	}

	k.Logger(ctx).Info("Distributed rewards to delegators", "amount", amount.String())
	return nil
}

// distributeToDEX distributes rewards to DEX pools (only years 1-2)
func (k Keeper) distributeToDEX(ctx sdk.Context, amount sdk.Coin, info types.HalvingInfo) error {
	// Check if we're in year 1 or 2 of distribution
	distributionStart := time.Unix(info.DistributionStart, 0)
	elapsed := ctx.BlockTime().Sub(distributionStart)
	
	// Only distribute to DEX in first 2 years
	if elapsed >= DEXDistributionPeriod {
		k.Logger(ctx).Info("DEX distribution period ended (after 2 years)", "cycle", info.CurrentCycle)
		return nil
	}

	// Keep DEX allocation in module account for bot to handle
	k.Logger(ctx).Info("DEX rewards allocated for bot distribution", 
		"amount", amount.String(),
		"cycle", info.CurrentCycle,
		"elapsed_days", int(elapsed.Hours()/24),
	)
	
	return nil
}

// GetAllValidatorUptimes returns all validator uptime records
func (k Keeper) GetAllValidatorUptimes(ctx sdk.Context) []types.ValidatorUptime {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorUptimeKey)
	defer iterator.Close()

	var uptimes []types.ValidatorUptime
	for ; iterator.Valid(); iterator.Next() {
		var uptime types.ValidatorUptime
		k.cdc.MustUnmarshal(iterator.Value(), &uptime)
		uptimes = append(uptimes, uptime)
	}

	return uptimes
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

// Helper function to check if validator bot is running (for slashing)
func (k Keeper) IsValidatorBotRunning(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	// This would be implemented with actual bot monitoring logic
	// For now, return true to avoid slashing during development
	return true
}

// SlashInactiveValidators slashes validators without running bots
func (k Keeper) SlashInactiveValidators(ctx sdk.Context) error {
	validators := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	
	for _, validator := range validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			continue
		}

		if !k.IsValidatorBotRunning(ctx, valAddr) {
			// Slash validator for not running mandatory bot
			k.Logger(ctx).Info("Slashing validator for not running mandatory bot",
				"validator", validator.OperatorAddress,
			)
			// Implement slashing logic here
		}
	}

	return nil
}