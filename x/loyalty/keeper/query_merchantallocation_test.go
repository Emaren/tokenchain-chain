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

func createNMerchantallocation(keeper keeper.Keeper, ctx context.Context, n int) []types.Merchantallocation {
	items := make([]types.Merchantallocation, n)
	for i := range items {
		denom := factoryDenom(sample.AccAddress(), "token"+strconv.Itoa(i))
		date := "2026-02-2" + strconv.Itoa(i)
		items[i].Key = date + "|" + denom
		items[i].Date = date
		items[i].Denom = denom
		items[i].ActivityScore = uint64(100 + i)
		items[i].BucketCAmount = uint64(1000 + i)
		items[i].StakersAmount = uint64(500 + i)
		items[i].TreasuryAmount = uint64(500)
		_ = keeper.Merchantallocation.Set(ctx, items[i].Key, items[i])
	}
	return items
}

func TestMerchantallocationQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNMerchantallocation(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetMerchantallocationRequest
		response *types.QueryGetMerchantallocationResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetMerchantallocationRequest{
				Key: msgs[0].Key,
			},
			response: &types.QueryGetMerchantallocationResponse{Merchantallocation: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetMerchantallocationRequest{
				Key: msgs[1].Key,
			},
			response: &types.QueryGetMerchantallocationResponse{Merchantallocation: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetMerchantallocationRequest{
				Key: "2026-02-26|factory/missing/missing",
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
			response, err := qs.GetMerchantallocation(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestMerchantallocationQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNMerchantallocation(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllMerchantallocationRequest {
		return &types.QueryAllMerchantallocationRequest{
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
			resp, err := qs.ListMerchantallocation(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Merchantallocation), step)
			require.Subset(t, msgs, resp.Merchantallocation)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListMerchantallocation(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Merchantallocation), step)
			require.Subset(t, msgs, resp.Merchantallocation)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListMerchantallocation(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Merchantallocation)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListMerchantallocation(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestMerchantallocationQueryFiltered(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)

	denomA := factoryDenom(sample.AccAddress(), "wheat")
	denomB := factoryDenom(sample.AccAddress(), "stone")

	records := []types.Merchantallocation{
		{
			Key:           "2026-02-25|" + denomA,
			Date:          "2026-02-25",
			Denom:         denomA,
			ActivityScore: 11,
			BucketCAmount: 1000,
		},
		{
			Key:           "2026-02-26|" + denomA,
			Date:          "2026-02-26",
			Denom:         denomA,
			ActivityScore: 22,
			BucketCAmount: 2000,
		},
		{
			Key:           "2026-02-26|" + denomB,
			Date:          "2026-02-26",
			Denom:         denomB,
			ActivityScore: 33,
			BucketCAmount: 3000,
		},
	}
	for _, record := range records {
		require.NoError(t, f.keeper.Merchantallocation.Set(f.ctx, record.Key, record))
	}

	t.Run("by_date", func(t *testing.T) {
		resp, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Date: "2026-02-26",
		})
		require.NoError(t, err)
		require.Len(t, resp.Merchantallocation, 2)
		for _, record := range resp.Merchantallocation {
			require.Equal(t, "2026-02-26", record.Date)
		}
	})

	t.Run("by_denom", func(t *testing.T) {
		resp, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Denom: denomA,
		})
		require.NoError(t, err)
		require.Len(t, resp.Merchantallocation, 2)
		for _, record := range resp.Merchantallocation {
			require.Equal(t, denomA, record.Denom)
		}
	})

	t.Run("by_date_and_denom", func(t *testing.T) {
		resp, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Date:  "2026-02-26",
			Denom: denomB,
		})
		require.NoError(t, err)
		require.Len(t, resp.Merchantallocation, 1)
		require.Equal(t, records[2], resp.Merchantallocation[0])
	})

	t.Run("pagination", func(t *testing.T) {
		page1, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Date: "2026-02-26",
			Pagination: &query.PageRequest{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.Len(t, page1.Merchantallocation, 1)
		require.NotEmpty(t, page1.Pagination.NextKey)

		page2, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Date: "2026-02-26",
			Pagination: &query.PageRequest{
				Key:   page1.Pagination.NextKey,
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.Len(t, page2.Merchantallocation, 1)
		require.Empty(t, page2.Pagination.NextKey)
	})

	t.Run("invalid_date", func(t *testing.T) {
		_, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Date: "2026/02/26",
		})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid date filter"))
	})

	t.Run("invalid_denom", func(t *testing.T) {
		_, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Denom: "BAD DENOM",
		})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid denom filter"))
	})

	t.Run("invalid_pagination_key", func(t *testing.T) {
		_, err := qs.FilterMerchantallocation(f.ctx, &types.QueryFilterMerchantallocationRequest{
			Pagination: &query.PageRequest{
				Key: []byte("oops"),
			},
		})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid pagination key"))
	})

	t.Run("nil_request", func(t *testing.T) {
		_, err := qs.FilterMerchantallocation(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
