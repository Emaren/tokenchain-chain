package keeper

import (
	"context"
	"errors"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListRecoveryoperation(ctx context.Context, req *types.QueryAllRecoveryoperationRequest) (*types.QueryAllRecoveryoperationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	recoveryoperations, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Recoveryoperation,
		req.Pagination,
		func(_ uint64, value types.Recoveryoperation) (types.Recoveryoperation, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRecoveryoperationResponse{Recoveryoperation: recoveryoperations, Pagination: pageRes}, nil
}

func (q queryServer) GetRecoveryoperation(ctx context.Context, req *types.QueryGetRecoveryoperationRequest) (*types.QueryGetRecoveryoperationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	recoveryoperation, err := q.k.Recoveryoperation.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetRecoveryoperationResponse{Recoveryoperation: recoveryoperation}, nil
}
