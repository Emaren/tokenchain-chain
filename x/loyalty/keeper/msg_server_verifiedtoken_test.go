package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func factoryDenom(issuer, subdenom string) string {
	return fmt.Sprintf("factory/%s/%s", issuer, subdenom)
}

func TestVerifiedtokenCreate_RecoveryPolicyMustExistInGroupModule(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)

	msg := baseVerifiedToken(creator, "recovery-policy")
	msg.SeizureOptIn = true
	msg.RecoveryGroupPolicy = sample.AccAddress()
	msg.RecoveryTimelockHours = 1

	_, err := srv.CreateVerifiedtoken(f.ctx, msg)
	require.ErrorIs(t, err, types.ErrRecoveryPolicy)

	f.groupKeeper.addPolicy(msg.RecoveryGroupPolicy)
	_, err = srv.CreateVerifiedtoken(f.ctx, msg)
	require.NoError(t, err)
}

func TestVerifiedtokenCreate_MainnetUsesMainnetTimelockMinimum(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	mainnetCtx := sdk.UnwrapSDKContext(f.ctx).WithChainID("tokenchain-1")

	msg := baseVerifiedToken(creator, "mainnet-timelock")
	msg.SeizureOptIn = true
	msg.RecoveryGroupPolicy = sample.AccAddress()
	msg.RecoveryTimelockHours = 1
	f.groupKeeper.addPolicy(msg.RecoveryGroupPolicy)

	_, err := srv.CreateVerifiedtoken(mainnetCtx, msg)
	require.ErrorIs(t, err, types.ErrRecoveryPolicy)

	msg.RecoveryTimelockHours = 24
	_, err = srv.CreateVerifiedtoken(mainnetCtx, msg)
	require.NoError(t, err)
}

func baseVerifiedToken(creator, denom string) *types.MsgCreateVerifiedtoken {
	return &types.MsgCreateVerifiedtoken{
		Creator:               creator,
		Denom:                 denom,
		Issuer:                creator,
		Name:                  "Token " + denom,
		Symbol:                "tt" + denom,
		Description:           "test token",
		Website:               "https://tokentap.ca",
		MaxSupply:             1_000_000,
		MintedSupply:          0,
		Verified:              true,
		SeizureOptIn:          false,
		RecoveryGroupPolicy:   "",
		RecoveryTimelockHours: 0,
	}
}

func TestVerifiedtokenMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)

	for i := 0; i < 5; i++ {
		subdenom := fmt.Sprintf("token%d", i)
		resp, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
		require.NoError(t, err)
		require.Equal(t, factoryDenom(creator, subdenom), resp.Denom)

		rst, err := f.keeper.Verifiedtoken.Get(f.ctx, factoryDenom(creator, subdenom))
		require.NoError(t, err)
		require.Equal(t, creator, rst.Creator)
		require.Equal(t, factoryDenom(creator, subdenom), rst.Denom)
		require.EqualValues(t, 0, rst.MintedSupply)
	}
}

func TestVerifiedtokenMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	unauthorizedAddr := sample.AccAddress()

	subdenom := "token0"
	denom := factoryDenom(creator, subdenom)
	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateVerifiedtoken
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgUpdateVerifiedtoken{
				Creator: "invalid",
				Denom:   denom,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgUpdateVerifiedtoken{
				Creator:   unauthorizedAddr,
				Denom:     denom,
				Issuer:    creator,
				Name:      "Updated",
				Symbol:    "updated",
				Website:   "https://tokentap.ca",
				MaxSupply: 1_000_000,
				Verified:  true,
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgUpdateVerifiedtoken{
				Creator:   creator,
				Denom:     factoryDenom(creator, "missingtoken"),
				Issuer:    creator,
				Name:      "Updated",
				Symbol:    "updated",
				Website:   "https://tokentap.ca",
				MaxSupply: 1_000_000,
				Verified:  true,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateVerifiedtoken{
				Creator:               creator,
				Denom:                 denom,
				Issuer:                creator,
				Name:                  "Updated Token",
				Symbol:                "updated",
				Description:           "updated description",
				Website:               "https://tokentap.ca",
				MaxSupply:             2_000_000,
				MintedSupply:          0,
				Verified:              true,
				SeizureOptIn:          false,
				RecoveryGroupPolicy:   "",
				RecoveryTimelockHours: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateVerifiedtoken(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			rst, err := f.keeper.Verifiedtoken.Get(f.ctx, denom)
			require.NoError(t, err)
			require.Equal(t, "Updated Token", rst.Name)
			require.EqualValues(t, 2_000_000, rst.MaxSupply)
		})
	}
}

func TestVerifiedtokenMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	unauthorizedAddr := sample.AccAddress()

	subdenom := "token0"
	denom := factoryDenom(creator, subdenom)
	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(creator, subdenom))
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteVerifiedtoken
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgDeleteVerifiedtoken{
				Creator: "invalid",
				Denom:   subdenom,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgDeleteVerifiedtoken{
				Creator: unauthorizedAddr,
				Denom:   denom,
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "key not found",
			request: &types.MsgDeleteVerifiedtoken{
				Creator: creator,
				Denom:   factoryDenom(creator, "missingtoken"),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgDeleteVerifiedtoken{
				Creator: creator,
				Denom:   denom,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteVerifiedtoken(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			found, err := f.keeper.Verifiedtoken.Has(f.ctx, denom)
			require.NoError(t, err)
			require.False(t, found)
		})
	}
}
