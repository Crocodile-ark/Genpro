package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Crocodile-ark/gxrchaind/x/halving/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the total set of halving parameters.
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// HalvingInfo returns the current halving cycle information.
func (k Keeper) HalvingInfo(goCtx context.Context, req *types.QueryHalvingInfoRequest) (*types.QueryHalvingInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	info, found := k.GetHalvingInfo(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "halving info not found")
	}

	return &types.QueryHalvingInfoResponse{HalvingInfo: info}, nil
}

// DistributionHistory returns the distribution history with pagination.
func (k Keeper) DistributionHistory(goCtx context.Context, req *types.QueryDistributionHistoryRequest) (*types.QueryDistributionHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	distributionStore := prefix.NewStore(store, types.LastDistributionKey)

	var records []types.DistributionRecord
	pageRes, err := query.Paginate(distributionStore, req.Pagination, func(key []byte, value []byte) error {
		var record types.DistributionRecord
		if err := k.cdc.Unmarshal(value, &record); err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDistributionHistoryResponse{
		DistributionRecords: records,
		Pagination:         pageRes,
	}, nil
}