package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestRewardPoolBalanceQuery(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	require.NoError(t, f.bankKeeper.MintCoins(f.ctx, types.ModuleName, sdk.NewCoins(
		sdk.NewCoin("utoken", sdkmath.NewInt(25000)),
		sdk.NewCoin("ustone", sdkmath.NewInt(7)),
	)))

	resp, err := qs.RewardPoolBalance(f.ctx, &types.QueryRewardPoolBalanceRequest{Denom: "utoken"})
	require.NoError(t, err)
	require.Equal(t, moduleAddr.String(), resp.ModuleAddress)
	require.Equal(t, "utoken", resp.Denom)
	require.Equal(t, "25000", resp.Amount)
}

func TestRewardPoolBalanceQueryZeroBalance(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)

	resp, err := qs.RewardPoolBalance(f.ctx, &types.QueryRewardPoolBalanceRequest{Denom: "utoken"})
	require.NoError(t, err)
	require.Equal(t, "0", resp.Amount)
}

func TestRewardPoolBalanceQueryInvalidRequest(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)

	_, err := qs.RewardPoolBalance(f.ctx, nil)
	require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))

	_, err = qs.RewardPoolBalance(f.ctx, &types.QueryRewardPoolBalanceRequest{})
	require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "denom is required"))

	_, err = qs.RewardPoolBalance(f.ctx, &types.QueryRewardPoolBalanceRequest{Denom: "BAD DENOM"})
	require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid denom"))
}
