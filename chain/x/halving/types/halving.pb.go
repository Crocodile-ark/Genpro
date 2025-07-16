package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params defines the parameters for the halving module.
type Params struct {
	HalvingCycleDuration time.Duration `protobuf:"bytes,1,opt,name=halving_cycle_duration,json=halvingCycleDuration,proto3,stdduration" json:"halving_cycle_duration"`
	ValidatorShare       sdk.Dec       `protobuf:"bytes,2,opt,name=validator_share,json=validatorShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"validator_share"`
	DelegatorShare       sdk.Dec       `protobuf:"bytes,3,opt,name=delegator_share,json=delegatorShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegator_share"`
	DexShare            sdk.Dec       `protobuf:"bytes,4,opt,name=dex_share,json=dexShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"dex_share"`
}

// HalvingInfo stores information about the current halving cycle
type HalvingInfo struct {
	CurrentCycle        uint64   `protobuf:"varint,1,opt,name=current_cycle,json=currentCycle,proto3" json:"current_cycle,omitempty"`
	CycleStartTime      int64    `protobuf:"varint,2,opt,name=cycle_start_time,json=cycleStartTime,proto3" json:"cycle_start_time,omitempty"`
	TotalFundsForCycle  sdk.Coin `protobuf:"bytes,3,opt,name=total_funds_for_cycle,json=totalFundsForCycle,proto3" json:"total_funds_for_cycle"`
	DistributedInCycle  sdk.Coin `protobuf:"bytes,4,opt,name=distributed_in_cycle,json=distributedInCycle,proto3" json:"distributed_in_cycle"`
}

// DistributionRecord tracks monthly distributions
type DistributionRecord struct {
	Timestamp int64    `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Amount    sdk.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
	Cycle     uint64   `protobuf:"varint,3,opt,name=cycle,proto3" json:"cycle,omitempty"`
	Month     uint64   `protobuf:"varint,4,opt,name=month,proto3" json:"month,omitempty"`
}

// GenesisState defines the halving module's genesis state.
type GenesisState struct {
	Params              Params               `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	HalvingInfo         HalvingInfo          `protobuf:"bytes,2,opt,name=halving_info,json=halvingInfo,proto3" json:"halving_info"`
	DistributionRecords []DistributionRecord `protobuf:"bytes,3,rep,name=distribution_records,json=distributionRecords,proto3" json:"distribution_records"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, halvingInfo HalvingInfo, distributionRecords []DistributionRecord) *GenesisState {
	return &GenesisState{
		Params:              params,
		HalvingInfo:         halvingInfo,
		DistributionRecords: distributionRecords,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), DefaultHalvingInfo(), []DistributionRecord{})
}

// DefaultHalvingInfo returns default halving info for genesis
func DefaultHalvingInfo() HalvingInfo {
	// GXR Total Supply: 85,000,000 GXR
	// Halving Fund: 21,250,000 GXR (25% of total supply)
	// First cycle allocation: 4,250,000 GXR (20% of halving fund)
	totalFunds := sdk.NewCoin("ugen", sdk.NewInt(425000000000000)) // 4,250,000 GXR in ugen
	
	return HalvingInfo{
		CurrentCycle:       1,
		CycleStartTime:     time.Now().Unix(), // Will be set to genesis time in real deployment
		TotalFundsForCycle: totalFunds,
		DistributedInCycle: sdk.NewCoin("ugen", sdk.ZeroInt()),
	}
}

// Validate performs basic validation of the GenesisState
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	
	// Validate HalvingInfo
	if gs.HalvingInfo.CurrentCycle == 0 || gs.HalvingInfo.CurrentCycle > 5 {
		return fmt.Errorf("invalid current cycle: %d, must be between 1 and 5", gs.HalvingInfo.CurrentCycle)
	}
	
	if gs.HalvingInfo.CycleStartTime <= 0 {
		return fmt.Errorf("invalid cycle start time: %d", gs.HalvingInfo.CycleStartTime)
	}
	
	return nil
}