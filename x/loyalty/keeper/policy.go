package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"tokenchain/x/loyalty/types"
)

func (k Keeper) authorityString() string {
	s, err := k.addressCodec.BytesToString(k.GetAuthority())
	if err != nil {
		panic(err)
	}
	return s
}

func (k Keeper) ensureAuthority(signer string) error {
	signerBz, err := k.addressCodec.StringToBytes(signer)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}

	if !bytes.Equal(signerBz, k.GetAuthority()) {
		return errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authorityString(), signer)
	}

	return nil
}

func (k Keeper) getParams(ctx context.Context) (types.Params, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return types.Params{}, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	return params, nil
}

func (k Keeper) creatorCanCreateToken(ctx context.Context, signer string, mode string) (bool, error) {
	if err := k.ensureAuthority(signer); err == nil {
		return true, nil
	}

	switch mode {
	case types.CreationModeAdminOnly:
		return false, nil
	case types.CreationModeAllowlisted:
		entry, err := k.Creatorallowlist.Get(ctx, signer)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				return false, nil
			}
			return false, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
		}
		return entry.Enabled, nil
	case types.CreationModePermissionless:
		return true, nil
	default:
		return false, errorsmod.Wrap(types.ErrInvalidCreationMode, mode)
	}
}

func rewardAccrualKey(address, denom string) string {
	return fmt.Sprintf("%s|%s", address, denom)
}
