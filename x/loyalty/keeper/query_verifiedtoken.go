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

func (q queryServer) ListVerifiedtoken(ctx context.Context, req *types.QueryAllVerifiedtokenRequest) (*types.QueryAllVerifiedtokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	verifiedtokens, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Verifiedtoken,
		req.Pagination,
		func(_ string, value types.Verifiedtoken) (types.Verifiedtoken, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVerifiedtokenResponse{Verifiedtoken: verifiedtokens, Pagination: pageRes}, nil
}

func (q queryServer) GetVerifiedtoken(ctx context.Context, req *types.QueryGetVerifiedtokenRequest) (*types.QueryGetVerifiedtokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.Verifiedtoken.Get(ctx, req.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetVerifiedtokenResponse{Verifiedtoken: val}, nil
}
