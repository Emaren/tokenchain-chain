package app

import (
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// enforceNoInflationGenesis rewrites selected module defaults before init genesis:
// - x/mint inflation fields are forced to zero
// - x/wasm upload policy defaults to permissioned uploads
func enforceNoInflationGenesis(cdc codec.Codec, req *abci.RequestInitChain) error {
	if len(req.AppStateBytes) == 0 {
		return nil
	}

	var appState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &appState); err != nil {
		return err
	}

	if err := enforceMintNoInflation(cdc, appState); err != nil {
		return err
	}
	if err := enforceWasmUploadPolicy(cdc, appState); err != nil {
		return err
	}

	updated, err := json.Marshal(appState)
	if err != nil {
		return err
	}
	req.AppStateBytes = updated
	return nil
}

func enforceMintNoInflation(cdc codec.Codec, appState map[string]json.RawMessage) error {
	mintStateBz, ok := appState[minttypes.ModuleName]
	if !ok {
		return nil
	}

	var mintState minttypes.GenesisState
	if err := cdc.UnmarshalJSON(mintStateBz, &mintState); err != nil {
		return err
	}

	mintState.Minter.Inflation = sdkmath.LegacyZeroDec()
	mintState.Minter.AnnualProvisions = sdkmath.LegacyZeroDec()
	mintState.Params.InflationRateChange = sdkmath.LegacyZeroDec()
	mintState.Params.InflationMin = sdkmath.LegacyZeroDec()
	mintState.Params.InflationMax = sdkmath.LegacyZeroDec()

	appState[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintState)
	return nil
}

func enforceWasmUploadPolicy(cdc codec.Codec, appState map[string]json.RawMessage) error {
	wasmStateBz, ok := appState[wasmtypes.ModuleName]
	if !ok {
		return nil
	}

	var wasmState wasmtypes.GenesisState
	if err := cdc.UnmarshalJSON(wasmStateBz, &wasmState); err != nil {
		return err
	}

	// Day-1 safety default: keep contract uploads permissioned unless explicitly opened by governance.
	if wasmState.Params.CodeUploadAccess.Permission == wasmtypes.AccessTypeEverybody ||
		wasmState.Params.CodeUploadAccess.Permission == wasmtypes.AccessTypeUnspecified {
		wasmState.Params.CodeUploadAccess = wasmtypes.AllowNobody
	}

	// Keep instantiation open by default so uploaded code can still be used broadly.
	if wasmState.Params.InstantiateDefaultPermission == wasmtypes.AccessTypeUnspecified {
		wasmState.Params.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
	}

	appState[wasmtypes.ModuleName] = cdc.MustMarshalJSON(&wasmState)
	return nil
}
