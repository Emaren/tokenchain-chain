package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"tokenchain/x/loyalty/types"
)

const businessTokenDisplayExponent uint32 = 6

func (k Keeper) setVerifiedTokenDenomMetadata(ctx context.Context, token types.Verifiedtoken) error {
	if err := sdk.ValidateDenom(token.Denom); err != nil {
		return errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}

	issuer, subdenom, err := splitTokenFactoryDenom(token.Denom)
	if err != nil {
		return err
	}
	if _, err := k.addressCodec.StringToBytes(issuer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid tokenfactory issuer address")
	}
	if err := validateSubdenom(subdenom); err != nil {
		return err
	}

	metadata := banktypes.Metadata{
		Description: token.Description,
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: token.Denom, Exponent: 0},
			{Denom: subdenom, Exponent: businessTokenDisplayExponent},
		},
		Base:    token.Denom,
		Display: subdenom,
		Name:    token.Name,
		Symbol:  token.Symbol,
		URI:     token.Website,
	}
	if err := metadata.Validate(); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	k.bankKeeper.SetDenomMetaData(ctx, metadata)
	return nil
}
