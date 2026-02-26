package app

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	wasm "github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	mint "github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestEnforceNoInflationGenesis(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{}, wasm.AppModuleBasic{})

	mintState := minttypes.DefaultGenesisState()
	mintState.Minter.Inflation = sdkmath.LegacyMustNewDecFromStr("0.13")
	mintState.Minter.AnnualProvisions = sdkmath.LegacyMustNewDecFromStr("1.25")
	mintState.Params.InflationRateChange = sdkmath.LegacyMustNewDecFromStr("0.13")
	mintState.Params.InflationMin = sdkmath.LegacyMustNewDecFromStr("0.07")
	mintState.Params.InflationMax = sdkmath.LegacyMustNewDecFromStr("0.20")

	appState := map[string]json.RawMessage{
		minttypes.ModuleName: encCfg.Codec.MustMarshalJSON(mintState),
	}
	appStateBz, err := json.Marshal(appState)
	require.NoError(t, err)

	req := &abci.RequestInitChain{AppStateBytes: appStateBz}
	require.NoError(t, enforceNoInflationGenesis(encCfg.Codec, req))

	var outState map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(req.AppStateBytes, &outState))

	var outMint minttypes.GenesisState
	require.NoError(t, encCfg.Codec.UnmarshalJSON(outState[minttypes.ModuleName], &outMint))
	require.True(t, outMint.Minter.Inflation.IsZero())
	require.True(t, outMint.Minter.AnnualProvisions.IsZero())
	require.True(t, outMint.Params.InflationRateChange.IsZero())
	require.True(t, outMint.Params.InflationMin.IsZero())
	require.True(t, outMint.Params.InflationMax.IsZero())
}

func TestEnforceWasmUploadPolicy(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(wasm.AppModuleBasic{})

	wasmState := wasmtypes.GenesisState{
		Params: wasmtypes.DefaultParams(),
	}
	require.Equal(t, wasmtypes.AccessTypeEverybody, wasmState.Params.CodeUploadAccess.Permission)

	appState := map[string]json.RawMessage{
		wasmtypes.ModuleName: encCfg.Codec.MustMarshalJSON(&wasmState),
	}
	appStateBz, err := json.Marshal(appState)
	require.NoError(t, err)

	req := &abci.RequestInitChain{AppStateBytes: appStateBz}
	require.NoError(t, enforceNoInflationGenesis(encCfg.Codec, req))

	var outState map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(req.AppStateBytes, &outState))

	var outWasm wasmtypes.GenesisState
	require.NoError(t, encCfg.Codec.UnmarshalJSON(outState[wasmtypes.ModuleName], &outWasm))
	require.Equal(t, wasmtypes.AccessTypeNobody, outWasm.Params.CodeUploadAccess.Permission)
	require.Equal(t, wasmtypes.AccessTypeEverybody, outWasm.Params.InstantiateDefaultPermission)
}
