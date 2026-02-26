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

func (q queryServer) ListRewardaccrual(ctx context.Context, req *types.QueryAllRewardaccrualRequest) (*types.QueryAllRewardaccrualResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	rewardaccruals, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Rewardaccrual,
		req.Pagination,
		func(_ string, value types.Rewardaccrual) (types.Rewardaccrual, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRewardaccrualResponse{Rewardaccrual: rewardaccruals, Pagination: pageRes}, nil
}

func (q queryServer) GetRewardaccrual(ctx context.Context, req *types.QueryGetRewardaccrualRequest) (*types.QueryGetRewardaccrualResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Rewardaccrual.Get(ctx, req.Key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetRewardaccrualResponse{Rewardaccrual: val}, nil
}
