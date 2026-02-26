package keeper

import (
	"context"
	"errors"
	"fmt"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ExecuteRecoveryTransfer(ctx context.Context, msg *types.MsgExecuteRecoveryTransfer) (*types.MsgExecuteRecoveryTransferResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	nowUnix := sdkCtx.BlockTime().Unix()
	if nowUnix < 0 {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "invalid block time")
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
	if uint64(nowUnix) < op.ExecuteAfter {
		return nil, errorsmod.Wrapf(types.ErrRecoveryTooEarly, "operation %d unlocks at %d", msg.Id, op.ExecuteAfter)
	}

	token, err := k.Verifiedtoken.Get(ctx, op.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrTokenNotFound, op.Denom)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if !token.SeizureOptIn {
		return nil, errorsmod.Wrap(types.ErrRecoveryPolicy, "token recovery is disabled")
	}

	fromAddr, err := k.addressCodec.StringToBytes(op.FromAddress)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid operation from_address")
	}
	toAddr, err := k.addressCodec.StringToBytes(op.ToAddress)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid operation to_address")
	}
	if op.Amount == 0 {
		return nil, errorsmod.Wrap(types.ErrRecoveryBadRequest, "operation amount must be greater than zero")
	}
	coins := sdk.NewCoins(sdk.NewCoin(op.Denom, sdkmath.NewIntFromUint64(op.Amount)))

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, coins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to collect recovery funds")
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coins); err != nil {
		return nil, errorsmod.Wrap(err, "failed to deliver recovery funds")
	}

	op.Status = types.RecoveryStatusExecuted
	op.ExecutedAt = uint64(nowUnix)
	op.CancelReason = ""
	if err := k.Recoveryoperation.Set(ctx, op.Id, op); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"loyalty.recovery_transfer_executed",
			sdk.NewAttribute("id", fmt.Sprintf("%d", op.Id)),
			sdk.NewAttribute("denom", op.Denom),
			sdk.NewAttribute("from_address", op.FromAddress),
			sdk.NewAttribute("to_address", op.ToAddress),
			sdk.NewAttribute("amount", fmt.Sprintf("%d", op.Amount)),
			sdk.NewAttribute("executed_at", fmt.Sprintf("%d", op.ExecutedAt)),
		),
	)

	return &types.MsgExecuteRecoveryTransferResponse{}, nil
}
