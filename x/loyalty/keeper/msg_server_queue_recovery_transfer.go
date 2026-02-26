package keeper

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) QueueRecoveryTransfer(ctx context.Context, msg *types.MsgQueueRecoveryTransfer) (*types.MsgQueueRecoveryTransferResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "amount must be greater than zero")
	}
	if _, err := k.addressCodec.StringToBytes(msg.FromAddress); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid from_address")
	}
	if _, err := k.addressCodec.StringToBytes(msg.ToAddress); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid to_address")
	}
	if strings.TrimSpace(msg.FromAddress) == strings.TrimSpace(msg.ToAddress) {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "from_address and to_address cannot be equal")
	}
	if err := k.validateTokenFactoryDenom(msg.Denom); err != nil {
		return nil, err
	}

	token, err := k.Verifiedtoken.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrTokenNotFound, msg.Denom)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if !token.SeizureOptIn {
		return nil, errorsmod.Wrap(types.ErrRecoveryPolicy, "token recovery is disabled (no-seizure default)")
	}
	if token.RecoveryGroupPolicy == "" {
		return nil, errorsmod.Wrap(types.ErrRecoveryPolicy, "recovery_group_policy is required for recovery-enabled tokens")
	}
	if err := k.ensureGroupPolicyExists(ctx, token.RecoveryGroupPolicy); err != nil {
		return nil, err
	}
	isAuthority := k.ensureAuthority(msg.Creator) == nil
	if msg.Creator != token.RecoveryGroupPolicy && !isAuthority {
		return nil, errorsmod.Wrap(types.ErrRecoveryUnauthorized, "only recovery group policy or authority can queue recovery")
	}
	params, err := k.getParams(ctx)
	if err != nil {
		return nil, err
	}
	minTimelock := minimumRecoveryTimelockHours(ctx, params)
	if token.RecoveryTimelockHours < minTimelock {
		return nil, errorsmod.Wrapf(types.ErrRecoveryPolicy, "recovery timelock must be at least %d hours for this network", minTimelock)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	nowUnix := sdkCtx.BlockTime().Unix()
	if nowUnix < 0 {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "invalid block time")
	}
	now := uint64(nowUnix)
	if token.RecoveryTimelockHours > math.MaxUint64/3600 {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "timelock overflow")
	}
	timelockSeconds := token.RecoveryTimelockHours * 3600
	if now > math.MaxUint64-timelockSeconds {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "execute_after overflow")
	}

	nextID, err := k.RecoveryoperationSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	op := types.Recoveryoperation{
		Id:           nextID,
		Denom:        msg.Denom,
		FromAddress:  msg.FromAddress,
		ToAddress:    msg.ToAddress,
		Amount:       msg.Amount,
		RequestedBy:  msg.Creator,
		ExecuteAfter: now + timelockSeconds,
		CreatedAt:    now,
		Status:       types.RecoveryStatusQueued,
		ExecutedAt:   0,
		CancelledAt:  0,
		CancelReason: "",
	}
	if err := k.Recoveryoperation.Set(ctx, op.Id, op); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"loyalty.recovery_transfer_queued",
			sdk.NewAttribute("id", fmt.Sprintf("%d", op.Id)),
			sdk.NewAttribute("denom", op.Denom),
			sdk.NewAttribute("from_address", op.FromAddress),
			sdk.NewAttribute("to_address", op.ToAddress),
			sdk.NewAttribute("amount", fmt.Sprintf("%d", op.Amount)),
			sdk.NewAttribute("execute_after", fmt.Sprintf("%d", op.ExecuteAfter)),
		),
	)

	return &types.MsgQueueRecoveryTransferResponse{
		Id:           op.Id,
		Status:       op.Status,
		ExecuteAfter: op.ExecuteAfter,
	}, nil
}
