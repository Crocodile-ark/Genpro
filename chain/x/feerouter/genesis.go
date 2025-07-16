package feerouter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Crocodile-ark/gxrchaind/x/feerouter/keeper"
	"github.com/Crocodile-ark/gxrchaind/x/feerouter/types"
)

// InitGenesis initializes the feerouter module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module parameters
	k.SetParams(ctx, genState.Params)

	// Set fee stats
	k.SetFeeStats(ctx, genState.FeeStats)

	// Set LP pools
	for _, pool := range genState.LPPools {
		k.SetLPPool(ctx, pool)
	}
}

// ExportGenesis returns the feerouter module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)

	if stats, found := k.GetFeeStats(ctx); found {
		genesis.FeeStats = stats
	}

	genesis.LPPools = k.GetAllLPPools(ctx)

	return genesis
}