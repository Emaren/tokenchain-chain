package keeper

import (
	"context"
	"errors"
	"fmt"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateRewardaccrual(ctx context.Context, msg *types.MsgCreateRewardaccrual) (*types.MsgCreateRewardaccrualResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}
	if _, err := k.addressCodec.StringToBytes(msg.Address); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid recipient address: %s", err))
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}
	if msg.Key != rewardAccrualKey(msg.Address, msg.Denom) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "key must be <address>|<denom>")
	}
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
	}

	// Check if the value already exists
	ok, err := k.Rewardaccrual.Has(ctx, msg.Key)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var rewardaccrual = types.Rewardaccrual{
		Creator:        msg.Creator,
		Key:            msg.Key,
		Address:        msg.Address,
		Denom:          msg.Denom,
		Amount:         msg.Amount,
		LastRollupDate: msg.LastRollupDate,
	}

	if err := k.Rewardaccrual.Set(ctx, rewardaccrual.Key, rewardaccrual); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateRewardaccrualResponse{}, nil
}

func (k msgServer) UpdateRewardaccrual(ctx context.Context, msg *types.MsgUpdateRewardaccrual) (*types.MsgUpdateRewardaccrualResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}
	if _, err := k.addressCodec.StringToBytes(msg.Address); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid recipient address: %s", err))
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}
	if msg.Key != rewardAccrualKey(msg.Address, msg.Denom) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "key must be <address>|<denom>")
	}
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
	}

	// Check if the value exists
	val, err := k.Rewardaccrual.Get(ctx, msg.Key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	var rewardaccrual = types.Rewardaccrual{
		Creator:        val.Creator,
		Key:            msg.Key,
		Address:        msg.Address,
		Denom:          msg.Denom,
		Amount:         msg.Amount,
		LastRollupDate: msg.LastRollupDate,
	}

	if err := k.Rewardaccrual.Set(ctx, rewardaccrual.Key, rewardaccrual); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update rewardaccrual")
	}

	return &types.MsgUpdateRewardaccrualResponse{}, nil
}

func (k msgServer) DeleteRewardaccrual(ctx context.Context, msg *types.MsgDeleteRewardaccrual) (*types.MsgDeleteRewardaccrualResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}

	// Check if the value exists
	_, err := k.Rewardaccrual.Get(ctx, msg.Key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	if err := k.Rewardaccrual.Remove(ctx, msg.Key); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove rewardaccrual")
	}

	return &types.MsgDeleteRewardaccrualResponse{}, nil
}
