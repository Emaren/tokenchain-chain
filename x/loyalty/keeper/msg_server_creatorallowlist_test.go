package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func authorityAddress(t *testing.T, f *fixture) string {
	t.Helper()
	authority, err := f.addressCodec.BytesToString(f.keeper.GetAuthority())
	require.NoError(t, err)
	return authority
}

func TestCreatorallowlistMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)

	for i := 0; i < 5; i++ {
		allowlisted := sample.AccAddress()
		_, err := srv.CreateCreatorallowlist(f.ctx, &types.MsgCreateCreatorallowlist{
			Creator: creator,
			Address: allowlisted,
			Enabled: true,
		})
		require.NoError(t, err)

		rst, err := f.keeper.Creatorallowlist.Get(f.ctx, allowlisted)
		require.NoError(t, err)
		require.Equal(t, creator, rst.Creator)
		require.True(t, rst.Enabled)
	}
}

func TestCreatorallowlistMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	unauthorizedAddr := sample.AccAddress()
	allowlisted := sample.AccAddress()

	_, err := srv.CreateCreatorallowlist(f.ctx, &types.MsgCreateCreatorallowlist{
		Creator: creator,
		Address: allowlisted,
		Enabled: true,
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateCreatorallowlist
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgUpdateCreatorallowlist{
				Creator: "invalid",
				Address: allowlisted,
				Enabled: false,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgUpdateCreatorallowlist{
				Creator: unauthorizedAddr,
				Address: allowlisted,
				Enabled: false,
			},
			err: types.ErrInvalidSigner,
		},
		{
			desc: "key not found",
			request: &types.MsgUpdateCreatorallowlist{
				Creator: creator,
				Address: sample.AccAddress(),
				Enabled: false,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateCreatorallowlist{
				Creator: creator,
				Address: allowlisted,
				Enabled: false,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateCreatorallowlist(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			rst, err := f.keeper.Creatorallowlist.Get(f.ctx, allowlisted)
			require.NoError(t, err)
			require.False(t, rst.Enabled)
		})
	}
}

func TestCreatorallowlistMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	unauthorizedAddr := sample.AccAddress()
	allowlisted := sample.AccAddress()

	_, err := srv.CreateCreatorallowlist(f.ctx, &types.MsgCreateCreatorallowlist{
		Creator: creator,
		Address: allowlisted,
		Enabled: true,
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteCreatorallowlist
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgDeleteCreatorallowlist{
				Creator: "invalid",
				Address: allowlisted,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgDeleteCreatorallowlist{
				Creator: unauthorizedAddr,
				Address: allowlisted,
			},
			err: types.ErrInvalidSigner,
		},
		{
			desc: "key not found",
			request: &types.MsgDeleteCreatorallowlist{
				Creator: creator,
				Address: sample.AccAddress(),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgDeleteCreatorallowlist{
				Creator: creator,
				Address: allowlisted,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteCreatorallowlist(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			found, err := f.keeper.Creatorallowlist.Has(f.ctx, allowlisted)
			require.NoError(t, err)
			require.False(t, found)
		})
	}
}
