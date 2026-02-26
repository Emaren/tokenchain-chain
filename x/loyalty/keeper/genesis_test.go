package keeper_test

import (
	"testing"

	"tokenchain/x/loyalty/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	f := initFixture(t)
	creator := authorityAddress(t, f)
	denom0 := factoryDenom(creator, "gen0")
	denom1 := factoryDenom(creator, "gen1")

	genesisState := types.GenesisState{
		Params:              types.DefaultParams(),
		CreatorallowlistMap: []types.Creatorallowlist{{Address: "0"}, {Address: "1"}},
		VerifiedtokenMap: []types.Verifiedtoken{
			{Denom: denom0, Issuer: creator, Name: "Genesis 0", Symbol: "G0"},
			{Denom: denom1, Issuer: creator, Name: "Genesis 1", Symbol: "G1"},
		},
		RewardaccrualMap:       []types.Rewardaccrual{{Key: "0"}, {Key: "1"}},
		RecoveryoperationList:  []types.Recoveryoperation{{Id: 0}, {Id: 1}},
		RecoveryoperationCount: 2,
		LastDailyRollupDate:    "2026-02-26",
	}
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.CreatorallowlistMap, got.CreatorallowlistMap)
	require.EqualExportedValues(t, genesisState.VerifiedtokenMap, got.VerifiedtokenMap)
	require.EqualExportedValues(t, genesisState.RewardaccrualMap, got.RewardaccrualMap)
	require.EqualExportedValues(t, genesisState.RecoveryoperationList, got.RecoveryoperationList)
	require.Equal(t, genesisState.RecoveryoperationCount, got.RecoveryoperationCount)

	metadata, ok := f.bankKeeper.denomMetadata[denom0]
	require.True(t, ok)
	require.Equal(t, denom0, metadata.Base)

}
