package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	grouptypes "github.com/cosmos/cosmos-sdk/x/group"

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

func minimumRecoveryTimelockHours(ctx context.Context, params types.Params) uint64 {
	chainID := strings.ToLower(sdk.UnwrapSDKContext(ctx).ChainID())
	if chainID == "" || strings.Contains(chainID, "testnet") || strings.Contains(chainID, "localnet") {
		return params.TestnetTimelockHours
	}
	return params.MainnetTimelockHours
}

func (k Keeper) validateRecoverySettings(ctx context.Context, seizureOptIn bool, recoveryPolicy string, recoveryTimelock uint64, params types.Params) (string, uint64, error) {
	if !seizureOptIn {
		return "", 0, nil
	}

	recoveryPolicy = strings.TrimSpace(recoveryPolicy)
	if recoveryPolicy == "" {
		return "", 0, errorsmod.Wrap(types.ErrRecoveryPolicy, "recovery group policy is required when seizure is enabled")
	}
	if _, err := k.addressCodec.StringToBytes(recoveryPolicy); err != nil {
		return "", 0, errorsmod.Wrap(types.ErrRecoveryPolicy, "recovery group policy must be a valid account address")
	}
	if err := k.ensureGroupPolicyExists(ctx, recoveryPolicy); err != nil {
		return "", 0, err
	}

	minTimelock := minimumRecoveryTimelockHours(ctx, params)
	if recoveryTimelock < minTimelock {
		return "", 0, errorsmod.Wrapf(types.ErrRecoveryPolicy, "recovery timelock must be at least %d hours for this network", minTimelock)
	}

	return recoveryPolicy, recoveryTimelock, nil
}

func (k Keeper) ensureGroupPolicyExists(ctx context.Context, policyAddress string) error {
	if k.groupKeeper == nil {
		return errorsmod.Wrap(sdkerrors.ErrLogic, "group keeper is not configured")
	}

	resp, err := k.groupKeeper.GroupPolicyInfo(ctx, &grouptypes.QueryGroupPolicyInfoRequest{Address: policyAddress})
	if err != nil || resp == nil || resp.Info == nil {
		return errorsmod.Wrap(types.ErrRecoveryPolicy, "recovery group policy must reference an existing x/group policy")
	}

	return nil
}

func normalizeMerchantIncentiveRouting(token types.Verifiedtoken) types.Verifiedtoken {
	if token.MerchantIncentiveStakersBps == 0 && token.MerchantIncentiveTreasuryBps == 0 {
		token.MerchantIncentiveStakersBps = types.DefaultMerchantIncentiveStakersBps
		token.MerchantIncentiveTreasuryBps = types.DefaultMerchantIncentiveTreasuryBps
	}
	return token
}
