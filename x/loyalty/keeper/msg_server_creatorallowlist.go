package keeper

import (
	"context"
	"errors"
	"fmt"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCreatorallowlist(ctx context.Context, msg *types.MsgCreateCreatorallowlist) (*types.MsgCreateCreatorallowlistResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}
	if _, err := k.addressCodec.StringToBytes(msg.Address); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid allowlist address: %s", err))
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}

	// Check if the value already exists
	ok, err := k.Creatorallowlist.Has(ctx, msg.Address)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var creatorallowlist = types.Creatorallowlist{
		Creator: msg.Creator,
		Address: msg.Address,
		Enabled: msg.Enabled,
	}

	if err := k.Creatorallowlist.Set(ctx, creatorallowlist.Address, creatorallowlist); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgCreateCreatorallowlistResponse{}, nil
}

func (k msgServer) UpdateCreatorallowlist(ctx context.Context, msg *types.MsgUpdateCreatorallowlist) (*types.MsgUpdateCreatorallowlistResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	if _, err := k.addressCodec.StringToBytes(msg.Address); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid allowlist address: %s", err))
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}

	// Check if the value exists
	val, err := k.Creatorallowlist.Get(ctx, msg.Address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	var creatorallowlist = types.Creatorallowlist{
		Creator: val.Creator,
		Address: msg.Address,
		Enabled: msg.Enabled,
	}

	if err := k.Creatorallowlist.Set(ctx, creatorallowlist.Address, creatorallowlist); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update creatorallowlist")
	}

	return &types.MsgUpdateCreatorallowlistResponse{}, nil
}

func (k msgServer) DeleteCreatorallowlist(ctx context.Context, msg *types.MsgDeleteCreatorallowlist) (*types.MsgDeleteCreatorallowlistResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	if _, err := k.addressCodec.StringToBytes(msg.Address); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid allowlist address: %s", err))
	}
	if err := k.ensureAuthority(msg.Creator); err != nil {
		return nil, err
	}

	// Check if the value exists
	_, err := k.Creatorallowlist.Get(ctx, msg.Address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	if err := k.Creatorallowlist.Remove(ctx, msg.Address); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove creatorallowlist")
	}

	return &types.MsgDeleteCreatorallowlistResponse{}, nil
}
