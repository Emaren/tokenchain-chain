package keeper

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"tokenchain/x/loyalty/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	tokenFactoryPrefix = "factory"
)

var subdenomPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9._-]{2,63}$`)

func (k msgServer) canonicalBusinessDenom(denomOrSubdenom, issuer string) (string, error) {
	denomOrSubdenom = strings.TrimSpace(denomOrSubdenom)
	if denomOrSubdenom == "" {
		return "", errorsmod.Wrap(types.ErrInvalidDenom, "denom cannot be empty")
	}

	if strings.HasPrefix(denomOrSubdenom, tokenFactoryPrefix+"/") {
		tfIssuer, _, err := splitTokenFactoryDenom(denomOrSubdenom)
		if err != nil {
			return "", err
		}
		if tfIssuer != issuer {
			return "", errorsmod.Wrapf(
				sdkerrors.ErrInvalidRequest,
				"tokenfactory denom issuer %s must match message issuer %s",
				tfIssuer,
				issuer,
			)
		}
		return denomOrSubdenom, k.validateTokenFactoryDenom(denomOrSubdenom)
	}

	if err := validateSubdenom(denomOrSubdenom); err != nil {
		return "", err
	}

	fullDenom := fmt.Sprintf("%s/%s/%s", tokenFactoryPrefix, issuer, denomOrSubdenom)
	return fullDenom, k.validateTokenFactoryDenom(fullDenom)
}

func (k msgServer) validateTokenFactoryDenom(denom string) error {
	if err := sdk.ValidateDenom(denom); err != nil {
		return errorsmod.Wrap(types.ErrInvalidDenom, err.Error())
	}

	tfIssuer, subdenom, err := splitTokenFactoryDenom(denom)
	if err != nil {
		return err
	}
	if _, err := k.addressCodec.StringToBytes(tfIssuer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "invalid tokenfactory issuer address")
	}
	if err := validateSubdenom(subdenom); err != nil {
		return err
	}

	return nil
}

func (k msgServer) resolveStoredDenom(ctx context.Context, creator, denom string) (string, error) {
	exists, err := k.Verifiedtoken.Has(ctx, denom)
	if err != nil {
		return "", errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if exists {
		return denom, nil
	}

	// Backward-compatible lookup for calls that pass a subdenom instead of the full tokenfactory denom.
	if strings.HasPrefix(denom, tokenFactoryPrefix+"/") {
		return denom, nil
	}
	candidate := fmt.Sprintf("%s/%s/%s", tokenFactoryPrefix, creator, strings.TrimSpace(denom))
	exists, err = k.Verifiedtoken.Has(ctx, candidate)
	if err != nil {
		return "", errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if exists {
		return candidate, nil
	}

	return denom, nil
}

func splitTokenFactoryDenom(denom string) (issuer string, subdenom string, err error) {
	parts := strings.Split(denom, "/")
	if len(parts) != 3 || parts[0] != tokenFactoryPrefix {
		return "", "", errorsmod.Wrap(types.ErrInvalidDenom, "expected tokenfactory denom format: factory/{issuer}/{subdenom}")
	}
	if parts[1] == "" || parts[2] == "" {
		return "", "", errorsmod.Wrap(types.ErrInvalidDenom, "tokenfactory denom issuer/subdenom cannot be empty")
	}
	return parts[1], parts[2], nil
}

func validateSubdenom(subdenom string) error {
	if !subdenomPattern.MatchString(subdenom) {
		return errorsmod.Wrap(
			types.ErrInvalidDenom,
			"subdenom must match ^[a-zA-Z][a-zA-Z0-9._-]{2,63}$",
		)
	}
	return nil
}
