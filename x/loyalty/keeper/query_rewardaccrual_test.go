package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func createNRewardaccrual(keeper keeper.Keeper, ctx context.Context, n int) []types.Rewardaccrual {
	items := make([]types.Rewardaccrual, n)
	for i := range items {
		items[i].Key = strconv.Itoa(i)
		items[i].Address = strconv.Itoa(i)
		items[i].Denom = strconv.Itoa(i)
		items[i].Amount = uint64(i)
		items[i].LastRollupDate = strconv.Itoa(i)
		_ = keeper.Rewardaccrual.Set(ctx, items[i].Key, items[i])
	}
	return items
}

func TestRewardaccrualQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRewardaccrual(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetRewardaccrualRequest
		response *types.QueryGetRewardaccrualResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRewardaccrualRequest{
				Key: msgs[0].Key,
			},
			response: &types.QueryGetRewardaccrualResponse{Rewardaccrual: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetRewardaccrualRequest{
				Key: msgs[1].Key,
			},
			response: &types.QueryGetRewardaccrualResponse{Rewardaccrual: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRewardaccrualRequest{
				Key: strconv.Itoa(100000),
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
			response, err := qs.GetRewardaccrual(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestRewardaccrualQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRewardaccrual(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRewardaccrualRequest {
		return &types.QueryAllRewardaccrualRequest{
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
			resp, err := qs.ListRewardaccrual(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Rewardaccrual), step)
			require.Subset(t, msgs, resp.Rewardaccrual)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListRewardaccrual(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Rewardaccrual), step)
			require.Subset(t, msgs, resp.Rewardaccrual)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListRewardaccrual(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Rewardaccrual)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListRewardaccrual(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestRewardaccrualQueryFiltered(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)

	addrA := sample.AccAddress()
	addrB := sample.AccAddress()
	denomA := "utoken"
	denomB := "factory/" + addrA + "/wheat"
	denomC := "factory/" + addrB + "/stone"

	records := []types.Rewardaccrual{
		{
			Key:            addrA + "|" + denomA,
			Address:        addrA,
			Denom:          denomA,
			Amount:         11,
			LastRollupDate: "2026-02-25",
		},
		{
			Key:            addrA + "|" + denomB,
			Address:        addrA,
			Denom:          denomB,
			Amount:         22,
			LastRollupDate: "2026-02-25",
		},
		{
			Key:            addrB + "|" + denomA,
			Address:        addrB,
			Denom:          denomA,
			Amount:         33,
			LastRollupDate: "2026-02-25",
		},
		{
			Key:            addrB + "|" + denomC,
			Address:        addrB,
			Denom:          denomC,
			Amount:         44,
			LastRollupDate: "2026-02-25",
		},
	}

	for _, record := range records {
		require.NoError(t, f.keeper.Rewardaccrual.Set(f.ctx, record.Key, record))
	}

	t.Run("by_address", func(t *testing.T) {
		resp, err := qs.FilterRewardaccrual(f.ctx, &types.QueryFilterRewardaccrualRequest{
			Address: addrA,
		})
		require.NoError(t, err)
		require.Len(t, resp.Rewardaccrual, 2)
		for _, record := range resp.Rewardaccrual {
			require.Equal(t, addrA, record.Address)
		}
	})

	t.Run("by_denom", func(t *testing.T) {
		resp, err := qs.FilterRewardaccrual(f.ctx, &types.QueryFilterRewardaccrualRequest{
			Denom: denomA,
		})
		require.NoError(t, err)
		require.Len(t, resp.Rewardaccrual, 2)
		for _, record := range resp.Rewardaccrual {
			require.Equal(t, denomA, record.Denom)
		}
	})

	t.Run("by_address_and_denom", func(t *testing.T) {
		resp, err := qs.FilterRewardaccrual(f.ctx, &types.QueryFilterRewardaccrualRequest{
			Address: addrB,
			Denom:   denomC,
		})
		require.NoError(t, err)
		require.Len(t, resp.Rewardaccrual, 1)
		require.Equal(t, records[3], resp.Rewardaccrual[0])
	})

	t.Run("pagination", func(t *testing.T) {
		req := &types.QueryFilterRewardaccrualRequest{
			Address: addrB,
			Pagination: &query.PageRequest{
				Limit: 1,
			},
		}

		page1, err := qs.FilterRewardaccrual(f.ctx, req)
		require.NoError(t, err)
		require.Len(t, page1.Rewardaccrual, 1)
		require.NotEmpty(t, page1.Pagination.NextKey)

		page2, err := qs.FilterRewardaccrual(f.ctx, &types.QueryFilterRewardaccrualRequest{
			Address: addrB,
			Pagination: &query.PageRequest{
				Key:   page1.Pagination.NextKey,
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.Len(t, page2.Rewardaccrual, 1)
		require.Empty(t, page2.Pagination.NextKey)
	})

	t.Run("invalid_address", func(t *testing.T) {
		_, err := qs.FilterRewardaccrual(f.ctx, &types.QueryFilterRewardaccrualRequest{
			Address: "not-an-address",
		})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid address filter"))
	})

	t.Run("invalid_denom", func(t *testing.T) {
		_, err := qs.FilterRewardaccrual(f.ctx, &types.QueryFilterRewardaccrualRequest{
			Denom: "BAD DENOM",
		})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid denom filter"))
	})

	t.Run("invalid_pagination_key", func(t *testing.T) {
		_, err := qs.FilterRewardaccrual(f.ctx, &types.QueryFilterRewardaccrualRequest{
			Address: addrA,
			Pagination: &query.PageRequest{
				Key: []byte("oops"),
			},
		})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid pagination key"))
	})

	t.Run("nil_request", func(t *testing.T) {
		_, err := qs.FilterRewardaccrual(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
