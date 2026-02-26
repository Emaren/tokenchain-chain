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

func (k msgServer) RecordRewardAccrual(ctx context.Context, msg *types.MsgRecordRewardAccrual) (*types.MsgRecordRewardAccrualResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	if _, err := k.addressCodec.StringToBytes(msg.Address); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid recipient address")
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}

	params, err := k.getParams(ctx)
	if err != nil {
		return nil, err
	}

	rollupDate := msg.Date
	if rollupDate == "" {
		location, err := loadRollupLocation(params.DailyRollupTimezone)
		if err != nil {
			return nil, err
		}
		rollupDate = sdk.UnwrapSDKContext(ctx).BlockTime().In(location).Format(rollupDateLayout)
	} else if _, err := time.Parse(rollupDateLayout, rollupDate); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "date must be YYYY-MM-DD")
	}

	key := rewardAccrualKey(msg.Address, msg.Denom)
	record, err := k.Rewardaccrual.Get(ctx, key)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}
		record = types.Rewardaccrual{
			Creator:        msg.Creator,
			Key:            key,
			Address:        msg.Address,
			Denom:          msg.Denom,
			Amount:         0,
			LastRollupDate: rollupDate,
		}
	}

	record.Amount += msg.Amount
	record.LastRollupDate = rollupDate
	if err := k.Rewardaccrual.Set(ctx, key, record); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgRecordRewardAccrualResponse{
		Key:         key,
		Address:     record.Address,
		Denom:       record.Denom,
		AmountAdded: msg.Amount,
		TotalAmount: record.Amount,
		RollupDate:  rollupDate,
	}, nil
}
