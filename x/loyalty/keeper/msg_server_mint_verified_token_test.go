package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestMintVerifiedTokenRequiresFullDenom(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)

	subdenom := "mintfull"
	denom := factoryDenom(creator, subdenom)
	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
	require.NoError(t, err)

	_, err = srv.MintVerifiedToken(f.ctx, &types.MsgMintVerifiedToken{
		Creator:   creator,
		Denom:     subdenom,
		Recipient: creator,
		Amount:    1,
	})
	require.ErrorIs(t, err, types.ErrInvalidDenom)

	resp, err := srv.MintVerifiedToken(f.ctx, &types.MsgMintVerifiedToken{
		Creator:   creator,
		Denom:     denom,
		Recipient: creator,
		Amount:    7,
	})
	require.NoError(t, err)
	require.Equal(t, denom, resp.Denom)
	require.EqualValues(t, 7, resp.MintedSupply)
}
