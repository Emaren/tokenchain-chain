package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestFundRewardPoolMsgServer(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := sample.AccAddress()
	creatorAddr := sdk.MustAccAddressFromBech32(creator)

	f.bankKeeper.accountBalances[creator] = sdk.NewCoins(sdk.NewCoin("utoken", sdkmath.NewInt(500)))

	resp, err := srv.FundRewardPool(f.ctx, &types.MsgFundRewardPool{
		Creator: creator,
		Denom:   "utoken",
		Amount:  200,
	})
	require.NoError(t, err)
	require.Equal(t, "utoken", resp.Denom)
	require.EqualValues(t, 200, resp.AmountFunded)
	require.Equal(t, "200", resp.NewBalance)
	require.NotEmpty(t, resp.ModuleAddress)

	creatorBalance := f.bankKeeper.SpendableCoins(f.ctx, creatorAddr).AmountOf("utoken")
	require.Equal(t, sdkmath.NewInt(300), creatorBalance)
}

func TestFundRewardPoolMsgServerValidation(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	_, err := srv.FundRewardPool(f.ctx, &types.MsgFundRewardPool{
		Creator: "not-an-address",
		Denom:   "utoken",
		Amount:  1,
	})
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	creator := sample.AccAddress()
	_, err = srv.FundRewardPool(f.ctx, &types.MsgFundRewardPool{
		Creator: creator,
		Denom:   "BAD DENOM",
		Amount:  1,
	})
	require.ErrorIs(t, err, types.ErrInvalidDenom)

	_, err = srv.FundRewardPool(f.ctx, &types.MsgFundRewardPool{
		Creator: creator,
		Denom:   "utoken",
		Amount:  0,
	})
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
}

func TestFundRewardPoolMsgServerInsufficientFunds(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := sample.AccAddress()

	_, err := srv.FundRewardPool(f.ctx, &types.MsgFundRewardPool{
		Creator: creator,
		Denom:   "utoken",
		Amount:  1,
	})
	require.ErrorIs(t, err, sdkerrors.ErrInsufficientFunds)
}
