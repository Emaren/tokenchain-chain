package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestRenounceTokenAdminSuccessAndMintLock(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	subdenom := "renounceok"
	denom := factoryDenom(creator, subdenom)

	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
	require.NoError(t, err)

	resp, err := srv.RenounceTokenAdmin(f.ctx, &types.MsgRenounceTokenAdmin{
		Creator: creator,
		Denom:   denom,
	})
	require.NoError(t, err)
	require.Equal(t, denom, resp.Denom)
	require.True(t, resp.AdminRenounced)

	token, err := f.keeper.Verifiedtoken.Get(f.ctx, denom)
	require.NoError(t, err)
	require.True(t, token.AdminRenounced)

	_, err = srv.MintVerifiedToken(f.ctx, &types.MsgMintVerifiedToken{
		Creator:   creator,
		Denom:     denom,
		Recipient: creator,
		Amount:    1,
	})
	require.ErrorIs(t, err, types.ErrAdminRenounced)
}

func TestRenounceTokenAdminUnauthorized(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	subdenom := "renounceauth"
	denom := factoryDenom(creator, subdenom)

	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
	require.NoError(t, err)

	_, err = srv.RenounceTokenAdmin(f.ctx, &types.MsgRenounceTokenAdmin{
		Creator: sample.AccAddress(),
		Denom:   denom,
	})
	require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
}

func TestRenounceTokenAdminAlreadyRenounced(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	subdenom := "renounceonce"
	denom := factoryDenom(creator, subdenom)

	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
	require.NoError(t, err)

	_, err = srv.RenounceTokenAdmin(f.ctx, &types.MsgRenounceTokenAdmin{
		Creator: creator,
		Denom:   denom,
	})
	require.NoError(t, err)

	_, err = srv.RenounceTokenAdmin(f.ctx, &types.MsgRenounceTokenAdmin{
		Creator: creator,
		Denom:   denom,
	})
	require.ErrorIs(t, err, types.ErrAdminRenounced)
}

func TestRenounceTokenAdminBlockedWhenSeizureEnabled(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	subdenom := "renounceseizure"
	denom := factoryDenom(creator, subdenom)
	policy := sample.AccAddress()
	f.groupKeeper.addPolicy(policy)

	msg := baseVerifiedToken(creator, subdenom)
	msg.SeizureOptIn = true
	msg.RecoveryGroupPolicy = policy
	msg.RecoveryTimelockHours = 1
	_, err := srv.CreateVerifiedtoken(f.ctx, msg)
	require.NoError(t, err)

	_, err = srv.RenounceTokenAdmin(f.ctx, &types.MsgRenounceTokenAdmin{
		Creator: creator,
		Denom:   denom,
	})
	require.ErrorIs(t, err, types.ErrAdminRenouncePolicy)
}

func TestUpdateVerifiedTokenBlockedAfterAdminRenounceForCapOrRecoveryChanges(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	subdenom := "renounceupdate"
	denom := factoryDenom(creator, subdenom)

	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
	require.NoError(t, err)

	_, err = srv.RenounceTokenAdmin(f.ctx, &types.MsgRenounceTokenAdmin{
		Creator: creator,
		Denom:   denom,
	})
	require.NoError(t, err)

	_, err = srv.UpdateVerifiedtoken(f.ctx, &types.MsgUpdateVerifiedtoken{
		Creator:               creator,
		Denom:                 denom,
		Issuer:                creator,
		Name:                  "Renounced Token",
		Symbol:                "rnu",
		Description:           "metadata update",
		Website:               "https://tokentap.ca",
		MaxSupply:             2_000_000,
		MintedSupply:          0,
		Verified:              true,
		SeizureOptIn:          false,
		RecoveryGroupPolicy:   "",
		RecoveryTimelockHours: 0,
	})
	require.ErrorIs(t, err, types.ErrAdminRenounced)
}
