package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestRewardaccrualMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)

	for i := 0; i < 5; i++ {
		addr := sample.AccAddress()
		key := addr + "|utoken"
		_, err := srv.CreateRewardaccrual(f.ctx, &types.MsgCreateRewardaccrual{
			Creator:        creator,
			Key:            key,
			Address:        addr,
			Denom:          "utoken",
			Amount:         100,
			LastRollupDate: "2026-02-25",
		})
		require.NoError(t, err)

		rst, err := f.keeper.Rewardaccrual.Get(f.ctx, key)
		require.NoError(t, err)
		require.Equal(t, creator, rst.Creator)
		require.EqualValues(t, 100, rst.Amount)
	}
}

func TestRewardaccrualMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	unauthorizedAddr := sample.AccAddress()
	addr := sample.AccAddress()
	key := addr + "|utoken"
	missingAddr := sample.AccAddress()
	missingKey := missingAddr + "|utoken"

	_, err := srv.CreateRewardaccrual(f.ctx, &types.MsgCreateRewardaccrual{
		Creator:        creator,
		Key:            key,
		Address:        addr,
		Denom:          "utoken",
		Amount:         100,
		LastRollupDate: "2026-02-25",
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateRewardaccrual
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgUpdateRewardaccrual{
				Creator: "invalid",
				Key:     key,
				Address: addr,
				Denom:   "utoken",
				Amount:  50,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgUpdateRewardaccrual{
				Creator: unauthorizedAddr,
				Key:     key,
				Address: addr,
				Denom:   "utoken",
				Amount:  50,
			},
			err: types.ErrInvalidSigner,
		},
		{
			desc: "key not found",
			request: &types.MsgUpdateRewardaccrual{
				Creator: creator,
				Key:     missingKey,
				Address: missingAddr,
				Denom:   "utoken",
				Amount:  50,
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgUpdateRewardaccrual{
				Creator:        creator,
				Key:            key,
				Address:        addr,
				Denom:          "utoken",
				Amount:         200,
				LastRollupDate: "2026-02-26",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateRewardaccrual(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			rst, err := f.keeper.Rewardaccrual.Get(f.ctx, key)
			require.NoError(t, err)
			require.EqualValues(t, 200, rst.Amount)
		})
	}
}

func TestRewardaccrualMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	unauthorizedAddr := sample.AccAddress()
	addr := sample.AccAddress()
	key := addr + "|utoken"

	_, err := srv.CreateRewardaccrual(f.ctx, &types.MsgCreateRewardaccrual{
		Creator:        creator,
		Key:            key,
		Address:        addr,
		Denom:          "utoken",
		Amount:         100,
		LastRollupDate: "2026-02-25",
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteRewardaccrual
		err     error
	}{
		{
			desc: "invalid address",
			request: &types.MsgDeleteRewardaccrual{
				Creator: "invalid",
				Key:     key,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgDeleteRewardaccrual{
				Creator: unauthorizedAddr,
				Key:     key,
			},
			err: types.ErrInvalidSigner,
		},
		{
			desc: "key not found",
			request: &types.MsgDeleteRewardaccrual{
				Creator: creator,
				Key:     sample.AccAddress() + "|utoken",
			},
			err: sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "completed",
			request: &types.MsgDeleteRewardaccrual{
				Creator: creator,
				Key:     key,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteRewardaccrual(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
			found, err := f.keeper.Rewardaccrual.Has(f.ctx, key)
			require.NoError(t, err)
			require.False(t, found)
		})
	}
}
