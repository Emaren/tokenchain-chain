package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/types"
)

func (q queryServer) RewardPoolBalance(ctx context.Context, req *types.QueryRewardPoolBalanceRequest) (*types.QueryRewardPoolBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	denom := strings.TrimSpace(req.Denom)
	if denom == "" {
		return nil, status.Error(codes.InvalidArgument, "denom is required")
	}
	if err := sdk.ValidateDenom(denom); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	balance := q.k.bankKeeper.SpendableCoins(ctx, moduleAddr).AmountOf(denom)

	return &types.QueryRewardPoolBalanceResponse{
		ModuleAddress: moduleAddr.String(),
		Denom:         denom,
		Amount:        balance.String(),
	}, nil
}
