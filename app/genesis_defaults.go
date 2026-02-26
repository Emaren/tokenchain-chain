package app

import (
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// enforceNoInflationGenesis rewrites x/mint genesis values so TOKEN has no inflation.
func enforceNoInflationGenesis(cdc codec.Codec, req *abci.RequestInitChain) error {
	if len(req.AppStateBytes) == 0 {
		return nil
	}

	var appState map[string]json.RawMessage
	if err := json.Unmarshal(req.AppStateBytes, &appState); err != nil {
		return err
	}

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
	updated, err := json.Marshal(appState)
	if err != nil {
		return err
	}
	req.AppStateBytes = updated
	return nil
}
