package keeper

import (
	"context"
	"errors"
	"time"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const rollupDateLayout = "2006-01-02"

// RunDailyRollup records the first block observed for a new local calendar day
// (according to params.daily_rollup_timezone) and emits a rollup event.
func (k Keeper) RunDailyRollup(ctx context.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	location, err := loadRollupLocation(params.DailyRollupTimezone)
	if err != nil {
		return err
	}
	today := sdk.UnwrapSDKContext(ctx).BlockTime().In(location).Format(rollupDateLayout)

	lastDate, err := k.LastDailyRollupDate.Get(ctx)
	if err == nil && lastDate == today {
		return nil
	}
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	if err := k.LastDailyRollupDate.Set(ctx, today); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDailyRollup,
			sdk.NewAttribute(types.AttributeKeyDate, today),
			sdk.NewAttribute(types.AttributeKeyTimezone, params.DailyRollupTimezone),
		),
	)

	return nil
}

func loadRollupLocation(name string) (*time.Location, error) {
	location, err := time.LoadLocation(name)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return location, nil
}
