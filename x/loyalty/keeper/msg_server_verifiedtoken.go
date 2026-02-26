package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"tokenchain/x/loyalty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateVerifiedtoken(ctx context.Context, msg *types.MsgCreateVerifiedtoken) (*types.MsgCreateVerifiedtokenResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if msg.Issuer == "" {
		msg.Issuer = msg.Creator
	}
	if _, err := k.addressCodec.StringToBytes(msg.Issuer); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid issuer address: %s", err))
	}
	canonicalDenom, err := k.canonicalBusinessDenom(msg.Denom, msg.Issuer)
	if err != nil {
		return nil, err
	}
	msg.Denom = canonicalDenom
	if msg.MaxSupply == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidCap, "max supply must be greater than zero")
	}
	if msg.MintedSupply != 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minted_supply must be zero on create; use mint-verified-token")
	}

	params, err := k.getParams(ctx)
	if err != nil {
		return nil, err
	}
	allowed, err := k.creatorCanCreateToken(ctx, msg.Creator, params.CreationMode)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errorsmod.Wrap(types.ErrCreatorNotAllowed, "creator is not authorized for the current creation mode")
	}

	// Check if the value already exists
	ok, err := k.Verifiedtoken.Has(ctx, msg.Denom)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	} else if ok {
		return nil, errorsmod.Wrap(types.ErrTokenExists, msg.Denom)
	}

	recoveryTimelock := msg.RecoveryTimelockHours
	recoveryPolicy := msg.RecoveryGroupPolicy
	recoveryPolicy, recoveryTimelock, err = k.validateRecoverySettings(ctx, msg.SeizureOptIn, recoveryPolicy, recoveryTimelock, params)
	if err != nil {
		return nil, err
	}
	merchantStakersBps := types.DefaultMerchantIncentiveStakersBps
	merchantTreasuryBps := types.DefaultMerchantIncentiveTreasuryBps
	if err := types.ValidateMerchantIncentiveRouting(merchantStakersBps, merchantTreasuryBps); err != nil {
		return nil, errorsmod.Wrap(types.ErrMerchantRouting, err.Error())
	}

	var verifiedtoken = types.Verifiedtoken{
		Creator:                      msg.Creator,
		Denom:                        msg.Denom,
		Issuer:                       msg.Issuer,
		Name:                         msg.Name,
		Symbol:                       msg.Symbol,
		Description:                  msg.Description,
		Website:                      msg.Website,
		MaxSupply:                    msg.MaxSupply,
		MintedSupply:                 0,
		Verified:                     msg.Verified,
		SeizureOptIn:                 msg.SeizureOptIn,
		RecoveryGroupPolicy:          recoveryPolicy,
		RecoveryTimelockHours:        recoveryTimelock,
		AdminRenounced:               false,
		MerchantIncentiveStakersBps:  merchantStakersBps,
		MerchantIncentiveTreasuryBps: merchantTreasuryBps,
	}

	if err := k.Verifiedtoken.Set(ctx, verifiedtoken.Denom, verifiedtoken); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	if err := k.setVerifiedTokenDenomMetadata(ctx, verifiedtoken); err != nil {
		return nil, err
	}

	return &types.MsgCreateVerifiedtokenResponse{Denom: verifiedtoken.Denom}, nil
}

func (k msgServer) UpdateVerifiedtoken(ctx context.Context, msg *types.MsgUpdateVerifiedtoken) (*types.MsgUpdateVerifiedtokenResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	lookupDenom, err := k.resolveStoredDenom(msg.Denom)
	if err != nil {
		return nil, err
	}
	msg.Denom = lookupDenom

	// Check if the value exists
	val, err := k.Verifiedtoken.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}
	val = normalizeMerchantIncentiveRouting(val)

	if _, err := k.addressCodec.StringToBytes(msg.Issuer); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid issuer address: %s", err))
	}
	if err := k.validateTokenFactoryDenom(msg.Denom); err != nil {
		return nil, err
	}
	denomIssuer, _, err := splitTokenFactoryDenom(msg.Denom)
	if err != nil {
		return nil, err
	}
	if msg.Issuer != val.Issuer {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "issuer cannot be changed after token creation")
	}
	if msg.Issuer != denomIssuer {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "issuer must match tokenfactory denom issuer")
	}
	if msg.MaxSupply < val.MintedSupply {
		return nil, errorsmod.Wrapf(types.ErrInvalidCap, "max supply cannot be lower than minted supply (%d)", val.MintedSupply)
	}
	if !val.SeizureOptIn && msg.SeizureOptIn && val.MintedSupply > 0 {
		return nil, errorsmod.Wrap(types.ErrRecoveryPolicy, "cannot enable seizure/recovery after token minting has started")
	}
	if val.AdminRenounced {
		if msg.MaxSupply != val.MaxSupply ||
			msg.SeizureOptIn != val.SeizureOptIn ||
			strings.TrimSpace(msg.RecoveryGroupPolicy) != val.RecoveryGroupPolicy ||
			msg.RecoveryTimelockHours != val.RecoveryTimelockHours {
			return nil, errorsmod.Wrap(types.ErrAdminRenounced, "admin-renounced token cannot change cap or recovery policy settings")
		}
	}

	// Checks if the msg creator is the same as the current owner or the authority.
	isAuthority := k.ensureAuthority(msg.Creator) == nil
	if msg.Creator != val.Creator && !isAuthority {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	params, err := k.getParams(ctx)
	if err != nil {
		return nil, err
	}
	recoveryTimelock := msg.RecoveryTimelockHours
	recoveryPolicy := msg.RecoveryGroupPolicy
	recoveryPolicy, recoveryTimelock, err = k.validateRecoverySettings(ctx, msg.SeizureOptIn, recoveryPolicy, recoveryTimelock, params)
	if err != nil {
		return nil, err
	}

	var verifiedtoken = types.Verifiedtoken{
		Creator:                      val.Creator,
		Denom:                        msg.Denom,
		Issuer:                       msg.Issuer,
		Name:                         msg.Name,
		Symbol:                       msg.Symbol,
		Description:                  msg.Description,
		Website:                      msg.Website,
		MaxSupply:                    msg.MaxSupply,
		MintedSupply:                 val.MintedSupply,
		Verified:                     msg.Verified,
		SeizureOptIn:                 msg.SeizureOptIn,
		RecoveryGroupPolicy:          recoveryPolicy,
		RecoveryTimelockHours:        recoveryTimelock,
		AdminRenounced:               val.AdminRenounced,
		MerchantIncentiveStakersBps:  val.MerchantIncentiveStakersBps,
		MerchantIncentiveTreasuryBps: val.MerchantIncentiveTreasuryBps,
	}

	if err := k.Verifiedtoken.Set(ctx, verifiedtoken.Denom, verifiedtoken); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update verifiedtoken")
	}
	if err := k.setVerifiedTokenDenomMetadata(ctx, verifiedtoken); err != nil {
		return nil, err
	}

	return &types.MsgUpdateVerifiedtokenResponse{}, nil
}

func (k msgServer) DeleteVerifiedtoken(ctx context.Context, msg *types.MsgDeleteVerifiedtoken) (*types.MsgDeleteVerifiedtokenResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid signer address: %s", err))
	}
	lookupDenom, err := k.resolveStoredDenom(msg.Denom)
	if err != nil {
		return nil, err
	}
	msg.Denom = lookupDenom

	// Check if the value exists
	val, err := k.Verifiedtoken.Get(ctx, msg.Denom)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, err.Error())
	}

	// Checks if the msg creator is the same as the current owner or the authority.
	isAuthority := k.ensureAuthority(msg.Creator) == nil
	if msg.Creator != val.Creator && !isAuthority {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}
	if val.MintedSupply > 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "cannot delete token with non-zero minted supply")
	}

	if err := k.Verifiedtoken.Remove(ctx, msg.Denom); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to remove verifiedtoken")
	}

	return &types.MsgDeleteVerifiedtokenResponse{}, nil
}
