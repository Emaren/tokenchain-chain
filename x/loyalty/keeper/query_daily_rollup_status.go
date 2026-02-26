package keeper

import (
	"context"
	"errors"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"tokenchain/x/loyalty/types"
)

func (q queryServer) DailyRollupStatus(ctx context.Context, req *types.QueryDailyRollupStatusRequest) (*types.QueryDailyRollupStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params, err := q.k.getParams(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	location, err := loadRollupLocation(params.DailyRollupTimezone)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid daily rollup timezone")
	}
	currentDate := sdk.UnwrapSDKContext(ctx).BlockTime().In(location).Format(rollupDateLayout)

	lastRollupDate := ""
	hasRolledToday := false
	storedLastRollupDate, err := q.k.LastDailyRollupDate.Get(ctx)
	if err == nil {
		lastRollupDate = storedLastRollupDate
		hasRolledToday = storedLastRollupDate == currentDate
	} else if !errors.Is(err, collections.ErrNotFound) {
		return nil, status.Error(codes.Internal, "internal error")
	}

	nextRollupDate := ""
	currentTime, err := time.ParseInLocation(rollupDateLayout, currentDate, location)
	if err == nil {
		nextRollupDate = currentTime.AddDate(0, 0, 1).Format(rollupDateLayout)
	}

	return &types.QueryDailyRollupStatusResponse{
		Timezone:            params.DailyRollupTimezone,
		CurrentLocalDate:    currentDate,
		LastDailyRollupDate: lastRollupDate,
		HasRolledToday:      hasRolledToday,
		NextRollupDate:      nextRollupDate,
	}, nil
}
