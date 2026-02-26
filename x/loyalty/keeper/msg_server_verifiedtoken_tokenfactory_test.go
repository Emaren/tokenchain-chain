package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestVerifiedtokenCreateAcceptsFullTokenFactoryDenom(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)

	fullDenom := factoryDenom(creator, "wheat")
	msg := baseVerifiedToken(creator, fullDenom)
	msg.Name = "Wheat"
	msg.Symbol = "WHEAT"

	_, err := srv.CreateVerifiedtoken(f.ctx, msg)
	require.NoError(t, err)

	stored, err := f.keeper.Verifiedtoken.Get(f.ctx, fullDenom)
	require.NoError(t, err)
	require.Equal(t, fullDenom, stored.Denom)
}

func TestVerifiedtokenCreateRejectsIssuerMismatchForFullDenom(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	otherIssuer := "tokenchain1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq9l7t8v"

	msg := baseVerifiedToken(creator, factoryDenom(otherIssuer, "stone"))
	_, err := srv.CreateVerifiedtoken(f.ctx, msg)
	require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
}

func TestVerifiedtokenCreateRejectsInvalidSubdenom(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)

	msg := baseVerifiedToken(creator, "bad/subdenom")
	_, err := srv.CreateVerifiedtoken(f.ctx, msg)
	require.ErrorIs(t, err, types.ErrInvalidDenom)
}
