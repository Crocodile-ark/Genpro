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

// QueryFeeStatsRequest is the request type for the Query/FeeStats RPC method.
type QueryFeeStatsRequest struct{}

// QueryFeeStatsResponse is the response type for the Query/FeeStats RPC method.
type QueryFeeStatsResponse struct {
	FeeStats FeeStats `protobuf:"bytes,1,opt,name=fee_stats,json=feeStats,proto3" json:"fee_stats"`
}

// QueryLPPoolsRequest is the request type for the Query/LPPools RPC method.
type QueryLPPoolsRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryLPPoolsResponse is the response type for the Query/LPPools RPC method.
type QueryLPPoolsResponse struct {
	LPPools    []LPPool            `protobuf:"bytes,1,rep,name=lp_pools,json=lpPools,proto3" json:"lp_pools"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}