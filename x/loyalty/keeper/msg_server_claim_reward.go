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

func (k msgServer) ClaimReward(ctx context.Context, msg *types.MsgClaimReward) (*types.MsgClaimRewardResponse, error) {
	creatorAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}

	key := rewardAccrualKey(msg.Creator, msg.Denom)
	record, err := k.Rewardaccrual.Get(ctx, key)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(types.ErrAccrualNotFound, key)
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if record.Address != msg.Creator || record.Denom != msg.Denom {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "reward accrual record does not match claim parameters")
	}
	if record.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "no reward balance to claim")
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, sdkmath.NewIntFromUint64(record.Amount)))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins); err != nil {
		return nil, err
	}
	if err := k.Rewardaccrual.Remove(ctx, key); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgClaimRewardResponse{}, nil
}
