package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/types"
)

func (q queryServer) GetVerifiedtokenByDenom(ctx context.Context, req *types.QueryGetVerifiedtokenByDenomRequest) (*types.QueryGetVerifiedtokenResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	denom := strings.TrimSpace(req.Denom)
	if denom == "" {
		return nil, status.Error(codes.InvalidArgument, "denom is required")
	}

	return q.GetVerifiedtoken(ctx, &types.QueryGetVerifiedtokenRequest{Denom: denom})
}
