package keeper

import (
	"context"
	"errors"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) MintVerifiedToken(ctx context.Context, msg *types.MsgMintVerifiedToken) (*types.MsgMintVerifiedTokenResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	recipientAddr, err := k.addressCodec.StringToBytes(msg.Recipient)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid recipient address")
	}
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}
	if msg.Denom == sdk.DefaultBondDenom {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, "cannot mint chain base denom via loyalty module")
	}

	token, err := k.Verifiedtoken.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrTokenNotFound, msg.Denom)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	isAuthority := k.ensureAuthority(msg.Creator) == nil
	if msg.Creator != token.Creator && !isAuthority {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only token owner or authority can mint")
	}

	if token.MintedSupply > token.MaxSupply-msg.Amount {
		return nil, errorsmod.Wrap(types.ErrCapExceeded, "mint amount exceeds configured cap")
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, sdkmath.NewIntFromUint64(msg.Amount)))
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return nil, err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipientAddr, coins); err != nil {
		return nil, err
	}

	token.MintedSupply += msg.Amount
	if err := k.Verifiedtoken.Set(ctx, token.Denom, token); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgMintVerifiedTokenResponse{}, nil
}
