package keeper

import (
	"context"

	"tokenchain/x/loyalty/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k msgServer) FundRewardPool(ctx context.Context, msg *types.MsgFundRewardPool) (*types.MsgFundRewardPoolResponse, error) {
	creatorAddr, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid creator address")
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}
	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount must be greater than zero")
	}

	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, sdkmath.NewIntFromUint64(msg.Amount)))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, creatorAddr, types.ModuleName, coins); err != nil {
		return nil, err
	}

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	newBalance := k.bankKeeper.SpendableCoins(ctx, moduleAddr).AmountOf(msg.Denom)

	return &types.MsgFundRewardPoolResponse{
		ModuleAddress: moduleAddr.String(),
		Denom:         msg.Denom,
		AmountFunded:  msg.Amount,
		NewBalance:    newBalance.String(),
	}, nil
}
