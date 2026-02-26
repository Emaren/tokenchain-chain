package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func createNRecoveryoperation(keeper keeper.Keeper, ctx context.Context, n int) []types.Recoveryoperation {
	items := make([]types.Recoveryoperation, n)
	for i := range items {
		iu := uint64(i)
		items[i].Id = iu
		items[i].Denom = strconv.Itoa(i)
		items[i].FromAddress = strconv.Itoa(i)
		items[i].ToAddress = strconv.Itoa(i)
		items[i].Amount = uint64(i)
		items[i].RequestedBy = strconv.Itoa(i)
		items[i].ExecuteAfter = uint64(i)
		items[i].CreatedAt = uint64(i)
		items[i].Status = strconv.Itoa(i)
		items[i].ExecutedAt = uint64(i)
		items[i].CancelledAt = uint64(i)
		items[i].CancelReason = strconv.Itoa(i)
		_ = keeper.Recoveryoperation.Set(ctx, iu, items[i])
		_ = keeper.RecoveryoperationSeq.Set(ctx, iu)
	}
	return items
}

func TestRecoveryoperationQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRecoveryoperation(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetRecoveryoperationRequest
		response *types.QueryGetRecoveryoperationResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetRecoveryoperationRequest{Id: msgs[0].Id},
			response: &types.QueryGetRecoveryoperationResponse{Recoveryoperation: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetRecoveryoperationRequest{Id: msgs[1].Id},
			response: &types.QueryGetRecoveryoperationResponse{Recoveryoperation: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetRecoveryoperationRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetRecoveryoperation(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestRecoveryoperationQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNRecoveryoperation(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRecoveryoperationRequest {
		return &types.QueryAllRecoveryoperationRequest{
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
			resp, err := qs.ListRecoveryoperation(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Recoveryoperation), step)
			require.Subset(t, msgs, resp.Recoveryoperation)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListRecoveryoperation(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Recoveryoperation), step)
			require.Subset(t, msgs, resp.Recoveryoperation)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListRecoveryoperation(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Recoveryoperation)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListRecoveryoperation(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestRecoveryoperationQueryFiltered(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)

	ops := []types.Recoveryoperation{
		{
			Id:           0,
			Denom:        "factory/a/alpha",
			FromAddress:  "addr1",
			ToAddress:    "addr2",
			Amount:       1,
			RequestedBy:  "policyA",
			ExecuteAfter: 1,
			CreatedAt:    1,
			Status:       types.RecoveryStatusQueued,
		},
		{
			Id:           1,
			Denom:        "factory/a/alpha",
			FromAddress:  "addr1",
			ToAddress:    "addr3",
			Amount:       2,
			RequestedBy:  "policyA",
			ExecuteAfter: 2,
			CreatedAt:    2,
			Status:       types.RecoveryStatusCancelled,
		},
		{
			Id:           2,
			Denom:        "factory/b/beta",
			FromAddress:  "addr4",
			ToAddress:    "addr5",
			Amount:       3,
			RequestedBy:  "policyB",
			ExecuteAfter: 3,
			CreatedAt:    3,
			Status:       types.RecoveryStatusQueued,
		},
		{
			Id:           3,
			Denom:        "factory/c/gamma",
			FromAddress:  "addr6",
			ToAddress:    "addr7",
			Amount:       4,
			RequestedBy:  "policyC",
			ExecuteAfter: 4,
			CreatedAt:    4,
			Status:       types.RecoveryStatusExecuted,
		},
	}
	for _, op := range ops {
		require.NoError(t, f.keeper.Recoveryoperation.Set(f.ctx, op.Id, op))
	}
	require.NoError(t, f.keeper.RecoveryoperationSeq.Set(f.ctx, uint64(len(ops))))

	t.Run("by_status", func(t *testing.T) {
		resp, err := qs.FilterRecoveryoperation(f.ctx, &types.QueryFilterRecoveryoperationRequest{
			Status: types.RecoveryStatusQueued,
		})
		require.NoError(t, err)
		require.Len(t, resp.Recoveryoperation, 2)
		for _, op := range resp.Recoveryoperation {
			require.Equal(t, types.RecoveryStatusQueued, op.Status)
		}
		require.EqualValues(t, 2, resp.Pagination.Total)
	})

	t.Run("by_compound_filters", func(t *testing.T) {
		resp, err := qs.FilterRecoveryoperation(f.ctx, &types.QueryFilterRecoveryoperationRequest{
			Denom:       "factory/a/alpha",
			RequestedBy: "policyA",
			FromAddress: "addr1",
		})
		require.NoError(t, err)
		require.Len(t, resp.Recoveryoperation, 2)
	})

	t.Run("pagination_by_key", func(t *testing.T) {
		req := &types.QueryFilterRecoveryoperationRequest{
			Status: types.RecoveryStatusQueued,
			Pagination: &query.PageRequest{
				Limit:      1,
				CountTotal: true,
			},
		}
		page1, err := qs.FilterRecoveryoperation(f.ctx, req)
		require.NoError(t, err)
		require.Len(t, page1.Recoveryoperation, 1)
		require.NotEmpty(t, page1.Pagination.NextKey)
		require.EqualValues(t, 2, page1.Pagination.Total)

		page2, err := qs.FilterRecoveryoperation(f.ctx, &types.QueryFilterRecoveryoperationRequest{
			Status: types.RecoveryStatusQueued,
			Pagination: &query.PageRequest{
				Key:        page1.Pagination.NextKey,
				Limit:      1,
				CountTotal: true,
			},
		})
		require.NoError(t, err)
		require.Len(t, page2.Recoveryoperation, 1)
		require.EqualValues(t, 2, page2.Pagination.Total)
	})

	t.Run("invalid_status", func(t *testing.T) {
		_, err := qs.FilterRecoveryoperation(f.ctx, &types.QueryFilterRecoveryoperationRequest{Status: "not-a-status"})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid status filter"))
	})

	t.Run("invalid_key", func(t *testing.T) {
		_, err := qs.FilterRecoveryoperation(f.ctx, &types.QueryFilterRecoveryoperationRequest{
			Pagination: &query.PageRequest{Key: []byte("bad-key")},
		})
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid pagination key"))
	})

	t.Run("nil_request", func(t *testing.T) {
		_, err := qs.FilterRecoveryoperation(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
