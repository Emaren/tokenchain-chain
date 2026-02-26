package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func createNVerifiedtoken(keeper keeper.Keeper, ctx context.Context, n int) []types.Verifiedtoken {
	items := make([]types.Verifiedtoken, n)
	for i := range items {
		items[i].Denom = strconv.Itoa(i)
		items[i].Issuer = strconv.Itoa(i)
		items[i].Name = strconv.Itoa(i)
		items[i].Symbol = strconv.Itoa(i)
		items[i].Description = strconv.Itoa(i)
		items[i].Website = strconv.Itoa(i)
		items[i].MaxSupply = uint64(i)
		items[i].MintedSupply = uint64(i)
		items[i].Verified = true
		items[i].SeizureOptIn = true
		items[i].RecoveryGroupPolicy = strconv.Itoa(i)
		items[i].RecoveryTimelockHours = uint64(i)
		items[i].MerchantIncentiveStakersBps = types.DefaultMerchantIncentiveStakersBps
		items[i].MerchantIncentiveTreasuryBps = types.DefaultMerchantIncentiveTreasuryBps
		_ = keeper.Verifiedtoken.Set(ctx, items[i].Denom, items[i])
	}
	return items
}

func TestVerifiedtokenQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNVerifiedtoken(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetVerifiedtokenRequest
		response *types.QueryGetVerifiedtokenResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetVerifiedtokenRequest{
				Denom: msgs[0].Denom,
			},
			response: &types.QueryGetVerifiedtokenResponse{Verifiedtoken: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetVerifiedtokenRequest{
				Denom: msgs[1].Denom,
			},
			response: &types.QueryGetVerifiedtokenResponse{Verifiedtoken: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetVerifiedtokenRequest{
				Denom: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetVerifiedtoken(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestVerifiedtokenQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNVerifiedtoken(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllVerifiedtokenRequest {
		return &types.QueryAllVerifiedtokenRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListVerifiedtoken(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Verifiedtoken), step)
			require.Subset(t, msgs, resp.Verifiedtoken)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListVerifiedtoken(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Verifiedtoken), step)
			require.Subset(t, msgs, resp.Verifiedtoken)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListVerifiedtoken(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Verifiedtoken)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListVerifiedtoken(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestVerifiedtokenQueryByDenomParam(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNVerifiedtoken(f.keeper, f.ctx, 2)

	t.Run("found", func(t *testing.T) {
		resp, err := qs.GetVerifiedtokenByDenom(f.ctx, &types.QueryGetVerifiedtokenByDenomRequest{
			Denom: msgs[0].Denom,
		})
		require.NoError(t, err)
		require.EqualExportedValues(t, &types.QueryGetVerifiedtokenResponse{Verifiedtoken: msgs[0]}, resp)
	})

	t.Run("missing", func(t *testing.T) {
		_, err := qs.GetVerifiedtokenByDenom(f.ctx, &types.QueryGetVerifiedtokenByDenomRequest{
			Denom: "does-not-exist",
		})
		require.ErrorIs(t, err, status.Error(codes.NotFound, "not found"))
	})

	t.Run("empty", func(t *testing.T) {
		_, err := qs.GetVerifiedtokenByDenom(f.ctx, &types.QueryGetVerifiedtokenByDenomRequest{})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "denom is required"))
	})

	t.Run("invalid_request", func(t *testing.T) {
		_, err := qs.GetVerifiedtokenByDenom(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestVerifiedtokenQueryLegacyRoutingDefaults(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	denom := "legacy-denom"
	legacy := types.Verifiedtoken{
		Denom:                        denom,
		Issuer:                       "legacy-issuer",
		Name:                         "legacy",
		Symbol:                       "LEG",
		MerchantIncentiveStakersBps:  0,
		MerchantIncentiveTreasuryBps: 0,
	}
	require.NoError(t, f.keeper.Verifiedtoken.Set(f.ctx, denom, legacy))

	resp, err := qs.GetVerifiedtoken(f.ctx, &types.QueryGetVerifiedtokenRequest{Denom: denom})
	require.NoError(t, err)
	require.EqualValues(t, types.DefaultMerchantIncentiveStakersBps, resp.Verifiedtoken.MerchantIncentiveStakersBps)
	require.EqualValues(t, types.DefaultMerchantIncentiveTreasuryBps, resp.Verifiedtoken.MerchantIncentiveTreasuryBps)
}
