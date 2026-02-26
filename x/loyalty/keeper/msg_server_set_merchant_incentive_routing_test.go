package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestSetMerchantIncentiveRouting(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	authority := authorityAddress(t, f)
	owner := sample.AccAddress()
	subdenom := "routing0"
	denom := factoryDenom(owner, subdenom)

	params := types.DefaultParams()
	params.CreationMode = types.CreationModePermissionless
	require.NoError(t, f.keeper.Params.Set(f.ctx, params))

	_, err := srv.CreateVerifiedtoken(f.ctx, baseVerifiedToken(owner, subdenom))
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgSetMerchantIncentiveRouting
		err     error
	}{
		{
			desc: "invalid creator address",
			request: &types.MsgSetMerchantIncentiveRouting{
				Creator:                      "invalid",
				Denom:                        denom,
				MerchantIncentiveStakersBps:  7000,
				MerchantIncentiveTreasuryBps: 3000,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "unauthorized",
			request: &types.MsgSetMerchantIncentiveRouting{
				Creator:                      sample.AccAddress(),
				Denom:                        denom,
				MerchantIncentiveStakersBps:  7000,
				MerchantIncentiveTreasuryBps: 3000,
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "token not found",
			request: &types.MsgSetMerchantIncentiveRouting{
				Creator:                      owner,
				Denom:                        factoryDenom(owner, "missing"),
				MerchantIncentiveStakersBps:  7000,
				MerchantIncentiveTreasuryBps: 3000,
			},
			err: types.ErrTokenNotFound,
		},
		{
			desc: "invalid routing sum",
			request: &types.MsgSetMerchantIncentiveRouting{
				Creator:                      owner,
				Denom:                        denom,
				MerchantIncentiveStakersBps:  6000,
				MerchantIncentiveTreasuryBps: 3000,
			},
			err: types.ErrMerchantRouting,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.SetMerchantIncentiveRouting(f.ctx, tc.request)
			require.ErrorIs(t, err, tc.err)
		})
	}

	resp, err := srv.SetMerchantIncentiveRouting(f.ctx, &types.MsgSetMerchantIncentiveRouting{
		Creator:                      owner,
		Denom:                        denom,
		MerchantIncentiveStakersBps:  7000,
		MerchantIncentiveTreasuryBps: 3000,
	})
	require.NoError(t, err)
	require.Equal(t, denom, resp.Denom)
	require.EqualValues(t, 7000, resp.MerchantIncentiveStakersBps)
	require.EqualValues(t, 3000, resp.MerchantIncentiveTreasuryBps)

	updated, err := f.keeper.Verifiedtoken.Get(f.ctx, denom)
	require.NoError(t, err)
	require.EqualValues(t, 7000, updated.MerchantIncentiveStakersBps)
	require.EqualValues(t, 3000, updated.MerchantIncentiveTreasuryBps)

	_, err = srv.SetMerchantIncentiveRouting(f.ctx, &types.MsgSetMerchantIncentiveRouting{
		Creator:                      authority,
		Denom:                        denom,
		MerchantIncentiveStakersBps:  4000,
		MerchantIncentiveTreasuryBps: 6000,
	})
	require.NoError(t, err)

	updated, err = f.keeper.Verifiedtoken.Get(f.ctx, denom)
	require.NoError(t, err)
	require.EqualValues(t, 4000, updated.MerchantIncentiveStakersBps)
	require.EqualValues(t, 6000, updated.MerchantIncentiveTreasuryBps)
}
