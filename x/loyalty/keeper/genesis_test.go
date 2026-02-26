package keeper_test

import (
	"testing"

	"tokenchain/x/loyalty/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:              types.DefaultParams(),
		CreatorallowlistMap: []types.Creatorallowlist{{Address: "0"}, {Address: "1"}}, VerifiedtokenMap: []types.Verifiedtoken{{Denom: "0"}, {Denom: "1"}}, RewardaccrualMap: []types.Rewardaccrual{{Key: "0"}, {Key: "1"}}, RecoveryoperationList: []types.Recoveryoperation{{Id: 0}, {Id: 1}},
		RecoveryoperationCount: 2,
		LastDailyRollupDate:    "2026-02-26",
	}
	f := initFixture(t)
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

}
