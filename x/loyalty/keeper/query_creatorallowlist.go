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

func (q queryServer) ListCreatorallowlist(ctx context.Context, req *types.QueryAllCreatorallowlistRequest) (*types.QueryAllCreatorallowlistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	creatorallowlists, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Creatorallowlist,
		req.Pagination,
		func(_ string, value types.Creatorallowlist) (types.Creatorallowlist, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCreatorallowlistResponse{Creatorallowlist: creatorallowlists, Pagination: pageRes}, nil
}

func (q queryServer) GetCreatorallowlist(ctx context.Context, req *types.QueryGetCreatorallowlistRequest) (*types.QueryGetCreatorallowlistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Creatorallowlist.Get(ctx, req.Address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetCreatorallowlistResponse{Creatorallowlist: val}, nil
}
