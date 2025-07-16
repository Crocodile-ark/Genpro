package halving

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Crocodile-ark/gxrchaind/x/halving/keeper"
)

// BeginBlocker checks for halving cycle advancement and monthly reward distribution
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Check if we need to advance to next halving cycle
	if err := k.CheckAndAdvanceCycle(ctx); err != nil {
		k.Logger(ctx).Error("Failed to check halving cycle advancement", "error", err)
	}

	// Check if it's time for monthly distribution
	// For demo purposes, we'll check if it's been more than 30 days since last distribution
	// In production, this would be more sophisticated (checking actual calendar months)
	lastDistribution := getLastDistributionTime(ctx, k)
	
	// If 30 days have passed since last distribution, distribute monthly rewards
	if ctx.BlockTime().Sub(lastDistribution) >= (30 * 24 * time.Hour) {
		if err := k.DistributeMonthlyRewards(ctx); err != nil {
			k.Logger(ctx).Error("Failed to distribute monthly rewards", "error", err)
		}
	}
}

// getLastDistributionTime gets the timestamp of the last distribution
func getLastDistributionTime(ctx sdk.Context, k keeper.Keeper) time.Time {
	records := k.GetAllDistributionRecords(ctx)
	if len(records) == 0 {
		// If no distributions yet, use genesis time
		return ctx.BlockTime().Add(-31 * 24 * time.Hour) // Force first distribution
	}

	// Find the most recent distribution
	var latest int64 = 0
	for _, record := range records {
		if record.Timestamp > latest {
			latest = record.Timestamp
		}
	}

	return time.Unix(latest, 0)
}