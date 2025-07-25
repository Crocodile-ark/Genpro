// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.21.12
// source: gxr/halving/halving.proto

package types

import (
	"fmt"
	"time"

	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

// Params defines the parameters for the halving module.
type Params struct {
	HalvingCycleDuration time.Duration `protobuf:"bytes,1,opt,name=halving_cycle_duration,json=halvingCycleDuration,proto3,stdduration" json:"halving_cycle_duration"`
	ValidatorShare       types.Dec     `protobuf:"bytes,2,opt,name=validator_share,json=validatorShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"validator_share"`
	DelegatorShare       types.Dec     `protobuf:"bytes,3,opt,name=delegator_share,json=delegatorShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegator_share"`
	DexShare             types.Dec     `protobuf:"bytes,4,opt,name=dex_share,json=dexShare,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"dex_share"`
}

// HalvingInfo stores information about the current halving cycle
type HalvingInfo struct {
	CurrentCycle       uint64     `protobuf:"varint,1,opt,name=current_cycle,json=currentCycle,proto3" json:"current_cycle,omitempty"`
	CycleStartTime     int64      `protobuf:"varint,2,opt,name=cycle_start_time,json=cycleStartTime,proto3" json:"cycle_start_time,omitempty"`
	TotalSupply        types.Coin `protobuf:"bytes,3,opt,name=total_supply,json=totalSupply,proto3" json:"total_supply"`
	HalvingFund        types.Coin `protobuf:"bytes,4,opt,name=halving_fund,json=halvingFund,proto3" json:"halving_fund"`
	DistributionActive bool       `protobuf:"varint,5,opt,name=distribution_active,json=distributionActive,proto3" json:"distribution_active,omitempty"`
	DistributionStart  int64      `protobuf:"varint,6,opt,name=distribution_start,json=distributionStart,proto3" json:"distribution_start,omitempty"`
	DistributedAmount  types.Coin `protobuf:"bytes,7,opt,name=distributed_amount,json=distributedAmount,proto3" json:"distributed_amount"`
	PauseStart         int64      `protobuf:"varint,8,opt,name=pause_start,json=pauseStart,proto3" json:"pause_start,omitempty"`
	LastMonthlyDistrib int64      `protobuf:"varint,9,opt,name=last_monthly_distrib,json=lastMonthlyDistrib,proto3" json:"last_monthly_distrib,omitempty"`
}

// ValidatorUptime tracks validator uptime for reward eligibility
type ValidatorUptime struct {
	ValidatorAddress string `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	CurrentMonth     uint64 `protobuf:"varint,2,opt,name=current_month,json=currentMonth,proto3" json:"current_month,omitempty"`
	InactiveDays     uint64 `protobuf:"varint,3,opt,name=inactive_days,json=inactiveDays,proto3" json:"inactive_days,omitempty"`
	LastCheck        int64  `protobuf:"varint,4,opt,name=last_check,json=lastCheck,proto3" json:"last_check,omitempty"`
}

// DistributionRecord tracks monthly distributions
type DistributionRecord struct {
	Timestamp int64      `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Amount    types.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
	Cycle     uint64     `protobuf:"varint,3,opt,name=cycle,proto3" json:"cycle,omitempty"`
	Month     uint64     `protobuf:"varint,4,opt,name=month,proto3" json:"month,omitempty"`
}

// GenesisState defines the halving module's genesis state.
type GenesisState struct {
	Params              Params               `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	HalvingInfo         HalvingInfo          `protobuf:"bytes,2,opt,name=halving_info,json=halvingInfo,proto3" json:"halving_info"`
	DistributionRecords []DistributionRecord `protobuf:"bytes,3,rep,name=distribution_records,json=distributionRecords,proto3" json:"distribution_records"`
	ValidatorUptimes    []ValidatorUptime    `protobuf:"bytes,4,rep,name=validator_uptimes,json=validatorUptimes,proto3" json:"validator_uptimes"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_halving, []int{0}
}

func (m *HalvingInfo) Reset()         { *m = HalvingInfo{} }
func (m *HalvingInfo) String() string { return proto.CompactTextString(m) }
func (*HalvingInfo) ProtoMessage()    {}
func (*HalvingInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_halving, []int{1}
}

func (m *ValidatorUptime) Reset()         { *m = ValidatorUptime{} }
func (m *ValidatorUptime) String() string { return proto.CompactTextString(m) }
func (*ValidatorUptime) ProtoMessage()    {}
func (*ValidatorUptime) Descriptor() ([]byte, []int) {
	return fileDescriptor_halving, []int{2}
}

func (m *DistributionRecord) Reset()         { *m = DistributionRecord{} }
func (m *DistributionRecord) String() string { return proto.CompactTextString(m) }
func (*DistributionRecord) ProtoMessage()    {}
func (*DistributionRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_halving, []int{3}
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_halving, []int{4}
}

func init() {
	proto.RegisterType((*Params)(nil), "gxr.halving.Params")
	proto.RegisterType((*HalvingInfo)(nil), "gxr.halving.HalvingInfo")
	proto.RegisterType((*ValidatorUptime)(nil), "gxr.halving.ValidatorUptime")
	proto.RegisterType((*DistributionRecord)(nil), "gxr.halving.DistributionRecord")
	proto.RegisterType((*GenesisState)(nil), "gxr.halving.GenesisState")
}

var fileDescriptor_halving = []byte{
	// Binary descriptor would go here in real implementation
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:              DefaultParams(),
		HalvingInfo:         HalvingInfo{},
		DistributionRecords: []DistributionRecord{},
		ValidatorUptimes:    []ValidatorUptime{},
	}
}

// DefaultHalvingInfo returns default halving info for genesis
func DefaultHalvingInfo() HalvingInfo {
	// GXR Total Supply: 85,000,000 GXR
	// Halving Fund: 21,250,000 GXR (25% of total supply)
	// First cycle allocation: 4,250,000 GXR (20% of halving fund)
	totalFunds := types.NewCoin("ugen", types.NewInt(425000000000000)) // 4,250,000 GXR in ugen
	
	return HalvingInfo{
		CurrentCycle:       1,
		CycleStartTime:     time.Now().Unix(), // Will be set to genesis time in real deployment
		TotalSupply:        types.NewCoin("ugen", types.NewInt(850000000000000)), // 85,000,000 GXR in ugen
		HalvingFund:        totalFunds,
		DistributionActive: false,
		DistributionStart:  0,
		DistributedAmount:  types.NewCoin("ugen", types.ZeroInt()),
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