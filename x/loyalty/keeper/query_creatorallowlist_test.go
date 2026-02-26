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

func createNCreatorallowlist(keeper keeper.Keeper, ctx context.Context, n int) []types.Creatorallowlist {
	items := make([]types.Creatorallowlist, n)
	for i := range items {
		items[i].Address = strconv.Itoa(i)
		items[i].Enabled = true
		_ = keeper.Creatorallowlist.Set(ctx, items[i].Address, items[i])
	}
	return items
}

func TestCreatorallowlistQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNCreatorallowlist(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetCreatorallowlistRequest
		response *types.QueryGetCreatorallowlistResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetCreatorallowlistRequest{
				Address: msgs[0].Address,
			},
			response: &types.QueryGetCreatorallowlistResponse{Creatorallowlist: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetCreatorallowlistRequest{
				Address: msgs[1].Address,
			},
			response: &types.QueryGetCreatorallowlistResponse{Creatorallowlist: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetCreatorallowlistRequest{
				Address: strconv.Itoa(100000),
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
			response, err := qs.GetCreatorallowlist(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestCreatorallowlistQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNCreatorallowlist(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllCreatorallowlistRequest {
		return &types.QueryAllCreatorallowlistRequest{
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
			resp, err := qs.ListCreatorallowlist(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Creatorallowlist), step)
			require.Subset(t, msgs, resp.Creatorallowlist)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListCreatorallowlist(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Creatorallowlist), step)
			require.Subset(t, msgs, resp.Creatorallowlist)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListCreatorallowlist(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Creatorallowlist)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListCreatorallowlist(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
