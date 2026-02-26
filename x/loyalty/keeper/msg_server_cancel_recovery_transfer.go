package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CancelRecoveryTransfer(ctx context.Context, msg *types.MsgCancelRecoveryTransfer) (*types.MsgCancelRecoveryTransferResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	reason := strings.TrimSpace(msg.Reason)
	if reason == "" {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "cancel reason cannot be empty")
	}
	if len(reason) > 512 {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "cancel reason exceeds 512 characters")
	}

	op, err := k.Recoveryoperation.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "recovery operation %d not found", msg.Id)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if op.Status != types.RecoveryStatusQueued {
		return nil, errorsmod.Wrapf(types.ErrRecoveryNotQueued, "recovery operation %d is in %s state", msg.Id, op.Status)
	}

	token, err := k.Verifiedtoken.Get(ctx, op.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrTokenNotFound, op.Denom)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	isAuthority := k.ensureAuthority(msg.Creator) == nil
	if msg.Creator != token.RecoveryGroupPolicy && !isAuthority {
		return nil, errorsmod.Wrap(types.ErrRecoveryUnauthorized, "only recovery group policy or authority can cancel recovery")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	nowUnix := sdkCtx.BlockTime().Unix()
	if nowUnix < 0 {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "invalid block time")
	}

	op.Status = types.RecoveryStatusCancelled
	op.CancelledAt = uint64(nowUnix)
	op.CancelReason = reason
	if err := k.Recoveryoperation.Set(ctx, op.Id, op); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"loyalty.recovery_transfer_cancelled",
			sdk.NewAttribute("id", fmt.Sprintf("%d", op.Id)),
			sdk.NewAttribute("denom", op.Denom),
			sdk.NewAttribute("cancelled_at", fmt.Sprintf("%d", op.CancelledAt)),
			sdk.NewAttribute("reason", op.CancelReason),
		),
	)

	return &types.MsgCancelRecoveryTransferResponse{
		Id:          op.Id,
		Status:      op.Status,
		CancelledAt: op.CancelledAt,
	}, nil
}
