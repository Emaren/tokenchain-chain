package keeper

import (
	"context"
	"errors"
	"strings"
	"time"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RecordMerchantAllocation(ctx context.Context, msg *types.MsgRecordMerchantAllocation) (*types.MsgRecordMerchantAllocationResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}
	if msg.ActivityScore == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "activity score must be greater than zero")
	}
	if msg.BucketCAmount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "bucket C amount must be greater than zero")
	}

	params, err := k.getParams(ctx)
	if err != nil {
		return nil, err
	}

	rollupDate := strings.TrimSpace(msg.Date)
	if rollupDate == "" {
		location, err := loadRollupLocation(params.DailyRollupTimezone)
		if err != nil {
			return nil, err
		}
		rollupDate = sdk.UnwrapSDKContext(ctx).BlockTime().In(location).Format(rollupDateLayout)
	} else if _, err := time.Parse(rollupDateLayout, rollupDate); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "date must be YYYY-MM-DD")
	}

	token, err := k.Verifiedtoken.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrTokenNotFound, "token not found")
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	token = normalizeMerchantIncentiveRouting(token)
	if err := types.ValidateMerchantIncentiveRouting(token.MerchantIncentiveStakersBps, token.MerchantIncentiveTreasuryBps); err != nil {
		return nil, errorsmod.Wrap(types.ErrMerchantRouting, err.Error())
	}

	stakersAmount := (msg.BucketCAmount * token.MerchantIncentiveStakersBps) / types.TotalBPS
	treasuryAmount := msg.BucketCAmount - stakersAmount
	key := merchantAllocationKey(rollupDate, msg.Denom)

	updated, err := k.Merchantallocation.Has(ctx, key)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	record := types.Merchantallocation{
		Creator:                      msg.Creator,
		Key:                          key,
		Date:                         rollupDate,
		Denom:                        msg.Denom,
		ActivityScore:                msg.ActivityScore,
		BucketCAmount:                msg.BucketCAmount,
		StakersAmount:                stakersAmount,
		TreasuryAmount:               treasuryAmount,
		MerchantIncentiveStakersBps:  token.MerchantIncentiveStakersBps,
		MerchantIncentiveTreasuryBps: token.MerchantIncentiveTreasuryBps,
	}

	if err := k.Merchantallocation.Set(ctx, key, record); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgRecordMerchantAllocationResponse{
		Key:                          key,
		Date:                         rollupDate,
		Denom:                        msg.Denom,
		ActivityScore:                msg.ActivityScore,
		BucketCAmount:                msg.BucketCAmount,
		StakersAmount:                stakersAmount,
		TreasuryAmount:               treasuryAmount,
		MerchantIncentiveStakersBps:  token.MerchantIncentiveStakersBps,
		MerchantIncentiveTreasuryBps: token.MerchantIncentiveTreasuryBps,
		Updated:                      updated,
	}, nil
}
