package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func TestQueryDailyRollupStatus(t *testing.T) {
	f := initFixture(t)
	queryServer := keeper.NewQueryServerImpl(f.keeper)

	ctx := sdk.UnwrapSDKContext(f.ctx).WithBlockTime(time.Date(2026, 2, 26, 8, 0, 0, 0, time.UTC))

	resp, err := queryServer.DailyRollupStatus(ctx, &types.QueryDailyRollupStatusRequest{})
	require.NoError(t, err)
	require.Equal(t, "America/Edmonton", resp.Timezone)
	require.Equal(t, "2026-02-26", resp.CurrentLocalDate)
	require.Equal(t, "", resp.LastDailyRollupDate)
	require.False(t, resp.HasRolledToday)
	require.Equal(t, "2026-02-27", resp.NextRollupDate)

	require.NoError(t, f.keeper.RunDailyRollup(ctx))
	resp, err = queryServer.DailyRollupStatus(ctx, &types.QueryDailyRollupStatusRequest{})
	require.NoError(t, err)
	require.Equal(t, "2026-02-26", resp.LastDailyRollupDate)
	require.True(t, resp.HasRolledToday)

	nextDayCtx := ctx.WithBlockTime(time.Date(2026, 2, 27, 8, 0, 0, 0, time.UTC))
	resp, err = queryServer.DailyRollupStatus(nextDayCtx, &types.QueryDailyRollupStatusRequest{})
	require.NoError(t, err)
	require.Equal(t, "2026-02-26", resp.LastDailyRollupDate)
	require.False(t, resp.HasRolledToday)
}

func TestQueryDailyRollupStatus_InvalidRequest(t *testing.T) {
	f := initFixture(t)
	queryServer := keeper.NewQueryServerImpl(f.keeper)

	_, err := queryServer.DailyRollupStatus(f.ctx, nil)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
