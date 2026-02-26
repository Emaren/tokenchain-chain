package keeper

import (
	"context"
	"errors"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RenounceTokenAdmin(ctx context.Context, msg *types.MsgRenounceTokenAdmin) (*types.MsgRenounceTokenAdminResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	lookupDenom, err := k.resolveStoredDenom(msg.Denom)
	if err != nil {
		return nil, err
	}

	token, err := k.Verifiedtoken.Get(ctx, lookupDenom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrTokenNotFound, lookupDenom)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	isAuthority := k.ensureAuthority(msg.Creator) == nil
	if msg.Creator != token.Creator && !isAuthority {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only token owner or authority can renounce token admin")
	}
	if token.AdminRenounced {
		return nil, errorsmod.Wrap(types.ErrAdminRenounced, "token admin already renounced")
	}
	if token.SeizureOptIn {
		return nil, errorsmod.Wrap(types.ErrAdminRenouncePolicy, "disable seizure/recovery before renouncing token admin")
	}

	token.AdminRenounced = true
	if err := k.Verifiedtoken.Set(ctx, token.Denom, token); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgRenounceTokenAdminResponse{
		Denom:          token.Denom,
		AdminRenounced: true,
	}, nil
}
