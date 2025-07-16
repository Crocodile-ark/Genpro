package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params defines the parameters for the feerouter module.
type Params struct {
	GeneralValidatorShare sdk.Dec `protobuf:"bytes,1,opt,name=general_validator_share,json=generalValidatorShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"general_validator_share"`
	GeneralDexShare       sdk.Dec `protobuf:"bytes,2,opt,name=general_dex_share,json=generalDexShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"general_dex_share"`
	GeneralPosShare       sdk.Dec `protobuf:"bytes,3,opt,name=general_pos_share,json=generalPosShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"general_pos_share"`
	FarmingValidatorShare sdk.Dec `protobuf:"bytes,4,opt,name=farming_validator_share,json=farmingValidatorShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"farming_validator_share"`
	FarmingDexShare       sdk.Dec `protobuf:"bytes,5,opt,name=farming_dex_share,json=farmingDexShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"farming_dex_share"`
	FarmingLPRewardShare  sdk.Dec `protobuf:"bytes,6,opt,name=farming_lp_reward_share,json=farmingLpRewardShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"farming_lp_reward_share"`
	FarmingPosShare       sdk.Dec `protobuf:"bytes,7,opt,name=farming_pos_share,json=farmingPosShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"farming_pos_share"`
}

// FeeStats tracks fee collection and distribution statistics
type FeeStats struct {
	TotalCollected   sdk.Coins `protobuf:"bytes,1,rep,name=total_collected,json=totalCollected,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_collected"`
	TotalToValidators sdk.Coins `protobuf:"bytes,2,rep,name=total_to_validators,json=totalToValidators,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_to_validators"`
	TotalToDex       sdk.Coins `protobuf:"bytes,3,rep,name=total_to_dex,json=totalToDex,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_to_dex"`
	TotalToPos       sdk.Coins `protobuf:"bytes,4,rep,name=total_to_pos,json=totalToPos,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_to_pos"`
	TotalToLPRewards sdk.Coins `protobuf:"bytes,5,rep,name=total_to_lp_rewards,json=totalToLpRewards,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_to_lp_rewards"`
}

// LPPool represents a liquidity pool that can receive farming rewards
type LPPool struct {
	Address      string    `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Name         string    `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Active       bool      `protobuf:"varint,3,opt,name=active,proto3" json:"active,omitempty"`
	TotalRewards sdk.Coins `protobuf:"bytes,4,rep,name=total_rewards,json=totalRewards,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_rewards"`
}

// GenesisState defines the feerouter module's genesis state.
type GenesisState struct {
	Params   Params   `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	FeeStats FeeStats `protobuf:"bytes,2,opt,name=fee_stats,json=feeStats,proto3" json:"fee_stats"`
	LPPools  []LPPool `protobuf:"bytes,3,rep,name=lp_pools,json=lpPools,proto3" json:"lp_pools"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, feeStats FeeStats, lpPools []LPPool) *GenesisState {
	return &GenesisState{
		Params:   params,
		FeeStats: feeStats,
		LPPools:  lpPools,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), DefaultFeeStats(), []LPPool{})
}

// DefaultFeeStats returns default fee stats for genesis
func DefaultFeeStats() FeeStats {
	return FeeStats{
		TotalCollected:   sdk.NewCoins(),
		TotalToValidators: sdk.NewCoins(),
		TotalToDex:       sdk.NewCoins(),
		TotalToPos:       sdk.NewCoins(),
		TotalToLPRewards: sdk.NewCoins(),
	}
}

// Validate performs basic validation of the GenesisState
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate LP pools
	for i, pool := range gs.LPPools {
		if pool.Address == "" {
			return fmt.Errorf("LP pool %d has empty address", i)
		}
		if pool.Name == "" {
			return fmt.Errorf("LP pool %d has empty name", i)
		}
	}

	return nil
}