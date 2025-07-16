package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// QueryParamsRequest is the request type for the Query/Params RPC method.
type QueryParamsRequest struct{}

// QueryParamsResponse is the response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

// QueryHalvingInfoRequest is the request type for the Query/HalvingInfo RPC method.
type QueryHalvingInfoRequest struct{}

// QueryHalvingInfoResponse is the response type for the Query/HalvingInfo RPC method.
type QueryHalvingInfoResponse struct {
	HalvingInfo HalvingInfo `protobuf:"bytes,1,opt,name=halving_info,json=halvingInfo,proto3" json:"halving_info"`
}

// QueryDistributionHistoryRequest is the request type for the Query/DistributionHistory RPC method.
type QueryDistributionHistoryRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryDistributionHistoryResponse is the response type for the Query/DistributionHistory RPC method.
type QueryDistributionHistoryResponse struct {
	DistributionRecords []DistributionRecord `protobuf:"bytes,1,rep,name=distribution_records,json=distributionRecords,proto3" json:"distribution_records"`
	Pagination          *query.PageResponse  `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}