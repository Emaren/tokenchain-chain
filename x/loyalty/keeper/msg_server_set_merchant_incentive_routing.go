package keeper

import (
	"context"
	"errors"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SetMerchantIncentiveRouting(ctx context.Context, msg *types.MsgSetMerchantIncentiveRouting) (*types.MsgSetMerchantIncentiveRoutingResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}

	lookupDenom, err := k.resolveStoredDenom(msg.Denom)
	if err != nil {
		return nil, err
	}
	msg.Denom = lookupDenom

	token, err := k.Verifiedtoken.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrTokenNotFound, msg.Denom)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	isAuthority := k.ensureAuthority(msg.Creator) == nil
	if msg.Creator != token.Creator && !isAuthority {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only token owner or authority can set merchant incentive routing")
	}

	if err := types.ValidateMerchantIncentiveRouting(msg.MerchantIncentiveStakersBps, msg.MerchantIncentiveTreasuryBps); err != nil {
		return nil, errorsmod.Wrap(types.ErrMerchantRouting, err.Error())
	}

	token.MerchantIncentiveStakersBps = msg.MerchantIncentiveStakersBps
	token.MerchantIncentiveTreasuryBps = msg.MerchantIncentiveTreasuryBps
	if err := k.Verifiedtoken.Set(ctx, token.Denom, token); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgSetMerchantIncentiveRoutingResponse{
		Denom:                        token.Denom,
		MerchantIncentiveStakersBps:  token.MerchantIncentiveStakersBps,
		MerchantIncentiveTreasuryBps: token.MerchantIncentiveTreasuryBps,
	}, nil
}
