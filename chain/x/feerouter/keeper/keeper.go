package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Crocodile-ark/gxrchaind/x/feerouter/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		paramstore paramtypes.Subspace

		accountKeeper authkeeper.AccountKeeper
		bankKeeper    bankkeeper.Keeper
		stakingKeeper *stakingkeeper.Keeper
		distrKeeper   distrkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
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
		distrKeeper:   distrKeeper,
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

// GetFeeStats gets the fee collection statistics
func (k Keeper) GetFeeStats(ctx sdk.Context) (types.FeeStats, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.FeeStatsKey)
	if bz == nil {
		return types.FeeStats{}, false
	}

	var stats types.FeeStats
	k.cdc.MustUnmarshal(bz, &stats)
	return stats, true
}

// SetFeeStats sets the fee collection statistics
func (k Keeper) SetFeeStats(ctx sdk.Context, stats types.FeeStats) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&stats)
	store.Set(types.FeeStatsKey, bz)
}

// GetLPPool gets a specific LP pool
func (k Keeper) GetLPPool(ctx sdk.Context, address string) (types.LPPool, bool) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.LPPoolsKey, []byte(address)...)
	bz := store.Get(key)
	if bz == nil {
		return types.LPPool{}, false
	}

	var pool types.LPPool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

// SetLPPool sets an LP pool
func (k Keeper) SetLPPool(ctx sdk.Context, pool types.LPPool) {
	store := ctx.KVStore(k.storeKey)
	key := append(types.LPPoolsKey, []byte(pool.Address)...)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
}

// GetAllLPPools gets all LP pools
func (k Keeper) GetAllLPPools(ctx sdk.Context) []types.LPPool {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LPPoolsKey)
	defer iterator.Close()

	var pools []types.LPPool
	for ; iterator.Valid(); iterator.Next() {
		var pool types.LPPool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)
		pools = append(pools, pool)
	}

	return pools
}

// ProcessTransactionFees processes transaction fees according to GXR specification
func (k Keeper) ProcessTransactionFees(ctx sdk.Context, fees sdk.Coins, isFarmingTransaction bool) error {
	if fees.IsZero() {
		return nil
	}

	params := k.GetParams(ctx)
	var validatorShare, dexShare, posShare, lpRewardShare sdk.Dec

	if isFarmingTransaction {
		// Farming transaction: 30/25/25/20
		validatorShare = params.FarmingValidatorShare
		dexShare = params.FarmingDexShare
		lpRewardShare = params.FarmingLPRewardShare
		posShare = params.FarmingPosShare
	} else {
		// General transaction: 40/30/30
		validatorShare = params.GeneralValidatorShare
		dexShare = params.GeneralDexShare
		posShare = params.GeneralPosShare
		lpRewardShare = sdk.ZeroDec()
	}

	// Calculate distribution amounts
	validatorAmount := make(sdk.Coins, len(fees))
	dexAmount := make(sdk.Coins, len(fees))
	posAmount := make(sdk.Coins, len(fees))
	lpRewardAmount := make(sdk.Coins, len(fees))

	for i, fee := range fees {
		validatorAmount[i] = sdk.NewCoin(fee.Denom, fee.Amount.ToDec().Mul(validatorShare).TruncateInt())
		dexAmount[i] = sdk.NewCoin(fee.Denom, fee.Amount.ToDec().Mul(dexShare).TruncateInt())
		posAmount[i] = sdk.NewCoin(fee.Denom, fee.Amount.ToDec().Mul(posShare).TruncateInt())
		if isFarmingTransaction {
			lpRewardAmount[i] = sdk.NewCoin(fee.Denom, fee.Amount.ToDec().Mul(lpRewardShare).TruncateInt())
		}
	}

	// Distribute to validators
	if err := k.distributeToValidators(ctx, validatorAmount); err != nil {
		return fmt.Errorf("failed to distribute to validators: %w", err)
	}

	// Distribute to DEX pools
	if err := k.distributeToDEX(ctx, dexAmount); err != nil {
		return fmt.Errorf("failed to distribute to DEX: %w", err)
	}

	// Distribute to PoS pool
	if err := k.distributeToPoS(ctx, posAmount); err != nil {
		return fmt.Errorf("failed to distribute to PoS: %w", err)
	}

	// Distribute to LP rewards (only for farming transactions)
	if isFarmingTransaction && !lpRewardAmount.IsZero() {
		if err := k.distributeToLPRewards(ctx, lpRewardAmount); err != nil {
			return fmt.Errorf("failed to distribute to LP rewards: %w", err)
		}
	}

	// Update fee stats
	k.updateFeeStats(ctx, fees, validatorAmount, dexAmount, posAmount, lpRewardAmount)

	k.Logger(ctx).Info("Transaction fees processed",
		"total_fees", fees.String(),
		"is_farming", isFarmingTransaction,
		"validator_amount", validatorAmount.String(),
		"dex_amount", dexAmount.String(),
		"pos_amount", posAmount.String(),
		"lp_reward_amount", lpRewardAmount.String(),
	)

	return nil
}

// distributeToValidators distributes fees to active validators
func (k Keeper) distributeToValidators(ctx sdk.Context, amount sdk.Coins) error {
	if amount.IsZero() {
		return nil
	}

	// Get all bonded validators
	validators := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if len(validators) == 0 {
		return fmt.Errorf("no bonded validators found")
	}

	// Distribute equally among active validators
	for _, coin := range amount {
		perValidatorAmount := coin.Amount.QuoRaw(int64(len(validators)))
		if perValidatorAmount.IsZero() {
			continue
		}

		for _, validator := range validators {
			valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
			if err != nil {
				continue
			}

			accAddr := sdk.AccAddress(valAddr)
			reward := sdk.NewCoin(coin.Denom, perValidatorAmount)

			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, accAddr, sdk.NewCoins(reward)); err != nil {
				k.Logger(ctx).Error("Failed to send fee to validator", "validator", validator.OperatorAddress, "error", err)
				continue
			}
		}
	}

	return nil
}

// distributeToDEX distributes fees to DEX pools for auto refill
func (k Keeper) distributeToDEX(ctx sdk.Context, amount sdk.Coins) error {
	if amount.IsZero() {
		return nil
	}

	// For now, keep in fee collector - will be handled by bot validator
	// In production, this would be sent to specific DEX pool addresses
	k.Logger(ctx).Info("DEX fees allocated for auto refill", "amount", amount.String())
	return nil
}

// distributeToPoS distributes fees to PoS pool (delegators)
func (k Keeper) distributeToPoS(ctx sdk.Context, amount sdk.Coins) error {
	if amount.IsZero() {
		return nil
	}

	// Add to distribution module fee pool for delegators
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	k.distrKeeper.SetFeePool(ctx, feePool)

	return nil
}

// distributeToLPRewards distributes fees to LP community rewards
func (k Keeper) distributeToLPRewards(ctx sdk.Context, amount sdk.Coins) error {
	if amount.IsZero() {
		return nil
	}

	// Get active LP pools
	pools := k.GetAllLPPools(ctx)
	activePools := []types.LPPool{}
	for _, pool := range pools {
		if pool.Active {
			activePools = append(activePools, pool)
		}
	}

	if len(activePools) == 0 {
		k.Logger(ctx).Info("No active LP pools found, keeping LP rewards in fee collector")
		return nil
	}

	// Distribute equally among active LP pools
	for _, coin := range amount {
		perPoolAmount := coin.Amount.QuoRaw(int64(len(activePools)))
		if perPoolAmount.IsZero() {
			continue
		}

		for _, pool := range activePools {
			poolAddr, err := sdk.AccAddressFromBech32(pool.Address)
			if err != nil {
				k.Logger(ctx).Error("Invalid LP pool address", "address", pool.Address, "error", err)
				continue
			}

			reward := sdk.NewCoin(coin.Denom, perPoolAmount)
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, poolAddr, sdk.NewCoins(reward)); err != nil {
				k.Logger(ctx).Error("Failed to send reward to LP pool", "pool", pool.Name, "error", err)
				continue
			}

			// Update pool stats
			pool.TotalRewards = pool.TotalRewards.Add(reward)
			k.SetLPPool(ctx, pool)
		}
	}

	return nil
}

// updateFeeStats updates the fee collection statistics
func (k Keeper) updateFeeStats(ctx sdk.Context, totalFees, validatorAmount, dexAmount, posAmount, lpRewardAmount sdk.Coins) {
	stats, found := k.GetFeeStats(ctx)
	if !found {
		stats = types.DefaultFeeStats()
	}

	stats.TotalCollected = stats.TotalCollected.Add(totalFees...)
	stats.TotalToValidators = stats.TotalToValidators.Add(validatorAmount...)
	stats.TotalToDex = stats.TotalToDex.Add(dexAmount...)
	stats.TotalToPos = stats.TotalToPos.Add(posAmount...)
	stats.TotalToLPRewards = stats.TotalToLPRewards.Add(lpRewardAmount...)

	k.SetFeeStats(ctx, stats)
}

// IsFarmingTransaction determines if a transaction is a farming transaction
// This is a simplified implementation - in production this would check
// specific transaction types or message types
func (k Keeper) IsFarmingTransaction(ctx sdk.Context, tx sdk.Tx) bool {
	// For now, return false - this would be implemented based on
	// specific criteria for identifying LP farming transactions
	return false
}