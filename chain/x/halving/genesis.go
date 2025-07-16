package halving

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Crocodile-ark/gxrchaind/x/halving/keeper"
	"github.com/Crocodile-ark/gxrchaind/x/halving/types"
)

// InitGenesis initializes the halving module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set module parameters
	k.SetParams(ctx, genState.Params)

	// Set halving info
	k.SetHalvingInfo(ctx, genState.HalvingInfo)

	// Set distribution records
	for _, record := range genState.DistributionRecords {
		k.SetDistributionRecord(ctx, record)
	}
}

// ExportGenesis returns the halving module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesisState()
	genesis.Params = k.GetParams(ctx)

	if info, found := k.GetHalvingInfo(ctx); found {
		genesis.HalvingInfo = info
	}

	genesis.DistributionRecords = k.GetAllDistributionRecords(ctx)

	return genesis
}