package keeper

import (
	"context"
	"errors"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListMerchantallocation(ctx context.Context, req *types.QueryAllMerchantallocationRequest) (*types.QueryAllMerchantallocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	records, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Merchantallocation,
		req.Pagination,
		func(_ string, value types.Merchantallocation) (types.Merchantallocation, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllMerchantallocationResponse{Merchantallocation: records, Pagination: pageRes}, nil
}

func (q queryServer) GetMerchantallocation(ctx context.Context, req *types.QueryGetMerchantallocationRequest) (*types.QueryGetMerchantallocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Merchantallocation.Get(ctx, req.Key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetMerchantallocationResponse{Merchantallocation: val}, nil
}
