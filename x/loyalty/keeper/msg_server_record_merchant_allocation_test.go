package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestRecordMerchantAllocation(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	authority := authorityAddress(t, f)
	owner := sample.AccAddress()
	subdenom := "alloc1"
	denom := factoryDenom(owner, subdenom)

	params := types.DefaultParams()
	params.CreationMode = types.CreationModePermissionless
	require.NoError(t, f.keeper.Params.Set(f.ctx, params))

	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(owner, subdenom))
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgRecordMerchantAllocation
		err     error
	}{
		{
			desc: "invalid authority address",
			request: &types.MsgRecordMerchantAllocation{
				Creator:       "invalid",
				Date:          "2026-02-26",
				Denom:         denom,
				ActivityScore: 10,
				BucketCAmount: 1000,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgRecordMerchantAllocation{
				Creator:       sample.AccAddress(),
				Date:          "2026-02-26",
				Denom:         denom,
				ActivityScore: 10,
				BucketCAmount: 1000,
			},
			err: types.ErrInvalidSigner,
		},
		{
			desc: "missing token",
			request: &types.MsgRecordMerchantAllocation{
				Creator:       authority,
				Date:          "2026-02-26",
				Denom:         factoryDenom(owner, "missing"),
				ActivityScore: 10,
				BucketCAmount: 1000,
			},
			err: types.ErrTokenNotFound,
		},
		{
			desc: "invalid date",
			request: &types.MsgRecordMerchantAllocation{
				Creator:       authority,
				Date:          "26-02-2026",
				Denom:         denom,
				ActivityScore: 10,
				BucketCAmount: 1000,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			desc: "zero activity score",
			request: &types.MsgRecordMerchantAllocation{
				Creator:       authority,
				Date:          "2026-02-26",
				Denom:         denom,
				ActivityScore: 0,
				BucketCAmount: 1000,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			desc: "zero bucket amount",
			request: &types.MsgRecordMerchantAllocation{
				Creator:       authority,
				Date:          "2026-02-26",
				Denom:         denom,
				ActivityScore: 10,
				BucketCAmount: 0,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.RecordMerchantAllocation(f.ctx, tc.request)
			require.ErrorIs(t, err, tc.err)
		})
	}

	resp, err := srv.RecordMerchantAllocation(f.ctx, &types.MsgRecordMerchantAllocation{
		Creator:       authority,
		Date:          "2026-02-26",
		Denom:         denom,
		ActivityScore: 100,
		BucketCAmount: 1000,
	})
	require.NoError(t, err)
	require.Equal(t, "2026-02-26|"+denom, resp.Key)
	require.False(t, resp.Updated)
	require.EqualValues(t, 500, resp.StakersAmount)
	require.EqualValues(t, 500, resp.TreasuryAmount)
	require.EqualValues(t, types.DefaultMerchantIncentiveStakersBps, resp.MerchantIncentiveStakersBps)
	require.EqualValues(t, types.DefaultMerchantIncentiveTreasuryBps, resp.MerchantIncentiveTreasuryBps)

	_, err = srv.SetMerchantIncentiveRouting(f.ctx, &types.MsgSetMerchantIncentiveRouting{
		Creator:                      owner,
		Denom:                        denom,
		MerchantIncentiveStakersBps:  7000,
		MerchantIncentiveTreasuryBps: 3000,
	})
	require.NoError(t, err)

	resp, err = srv.RecordMerchantAllocation(f.ctx, &types.MsgRecordMerchantAllocation{
		Creator:       authority,
		Date:          "2026-02-26",
		Denom:         denom,
		ActivityScore: 200,
		BucketCAmount: 2000,
	})
	require.NoError(t, err)
	require.True(t, resp.Updated)
	require.EqualValues(t, 1400, resp.StakersAmount)
	require.EqualValues(t, 600, resp.TreasuryAmount)
	require.EqualValues(t, 7000, resp.MerchantIncentiveStakersBps)
	require.EqualValues(t, 3000, resp.MerchantIncentiveTreasuryBps)

	record, err := f.keeper.Merchantallocation.Get(f.ctx, resp.Key)
	require.NoError(t, err)
	require.EqualValues(t, 200, record.ActivityScore)
	require.EqualValues(t, 2000, record.BucketCAmount)
	require.EqualValues(t, 1400, record.StakersAmount)
	require.EqualValues(t, 600, record.TreasuryAmount)
}
