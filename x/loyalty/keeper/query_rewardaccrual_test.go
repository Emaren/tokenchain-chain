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
