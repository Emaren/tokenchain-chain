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
