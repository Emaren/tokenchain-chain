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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

	claimAmount := sdkmath.NewIntFromUint64(record.Amount)
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	moduleBalance := k.bankKeeper.SpendableCoins(ctx, moduleAddr).AmountOf(msg.Denom)
	if moduleBalance.LT(claimAmount) {
		return nil, errorsmod.Wrapf(
			types.ErrRewardPoolInsufficient,
			"module balance %s%s is smaller than claim %s%s",
			moduleBalance.String(),
			msg.Denom,
			claimAmount.String(),
			msg.Denom,
		)
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, sdkmath.NewIntFromUint64(record.Amount)))
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creatorAddr, coins); err != nil {
		return nil, err
	}
	if err := k.Rewardaccrual.Remove(ctx, key); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	return &types.MsgClaimRewardResponse{
		Address:       msg.Creator,
		Denom:         msg.Denom,
		AmountClaimed: record.Amount,
	}, nil
}
