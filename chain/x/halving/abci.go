package halving

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Crocodile-ark/gxrchaind/x/halving/keeper"
)

// BeginBlocker checks for halving cycle advancement and distribution status
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Check if we need to advance to next halving cycle (every 5 years)
	if err := k.CheckAndAdvanceHalvingCycle(ctx); err != nil {
		k.Logger(ctx).Error("Failed to check halving cycle advancement", "error", err)
	}

	// Check if distribution period should be updated (2 years active, 3 years inactive)
	if err := k.CheckAndUpdateDistributionStatus(ctx); err != nil {
		k.Logger(ctx).Error("Failed to check distribution status", "error", err)
	}

	// Check if it's time for monthly distribution
	if shouldDistributeMonthly(ctx) {
		if err := k.DistributeHalvingRewards(ctx); err != nil {
			k.Logger(ctx).Error("Failed to distribute monthly rewards", "error", err)
		}
	}
}

// shouldDistributeMonthly checks if it's time for monthly distribution
func shouldDistributeMonthly(ctx sdk.Context) bool {
	// Get the last distribution time from state
	// For simplicity, we'll check if it's a new month (approximately every 30 days)
	
	// This is a simplified check - in production, you might want to store
	// the last distribution time in the state and check against it
	currentTime := ctx.BlockTime()
	
	// Check if it's the first day of a new month (simplified logic)
	// In production, you'd store the last distribution timestamp
	dayOfMonth := currentTime.Day()
	
	// Distribute on the 1st of each month
	return dayOfMonth == 1
}