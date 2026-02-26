package types_test

import (
	"testing"

	"tokenchain/x/loyalty/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params:                types.DefaultParams(),
				CreatorallowlistMap:   []types.Creatorallowlist{{Address: "0"}, {Address: "1"}},
				VerifiedtokenMap:      []types.Verifiedtoken{{Denom: "token0"}, {Denom: "token1"}},
				RewardaccrualMap:      []types.Rewardaccrual{{Key: "0"}, {Key: "1"}},
				RecoveryoperationList: []types.Recoveryoperation{{Id: 0}, {Id: 1}}, RecoveryoperationCount: 2,
			}, valid: true,
		}, {
			desc: "duplicated creatorallowlist",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				CreatorallowlistMap: []types.Creatorallowlist{
					{
						Address: "0",
					},
					{
						Address: "0",
					},
				},
				VerifiedtokenMap: []types.Verifiedtoken{{Denom: "token0"}, {Denom: "token1"}}, RewardaccrualMap: []types.Rewardaccrual{{Key: "0"}, {Key: "1"}}, RecoveryoperationList: []types.Recoveryoperation{{Id: 0}, {Id: 1}}, RecoveryoperationCount: 2,
			}, valid: false,
		}, {
			desc: "duplicated verifiedtoken",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				VerifiedtokenMap: []types.Verifiedtoken{
					{
						Denom: "token0",
					},
					{
						Denom: "token0",
					},
				},
				RewardaccrualMap: []types.Rewardaccrual{{Key: "0"}, {Key: "1"}}, RecoveryoperationList: []types.Recoveryoperation{{Id: 0}, {Id: 1}}, RecoveryoperationCount: 2,
			}, valid: false,
		}, {
			desc: "duplicated rewardaccrual",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RewardaccrualMap: []types.Rewardaccrual{
					{
						Key: "0",
					},
					{
						Key: "0",
					},
				},
				RecoveryoperationList: []types.Recoveryoperation{{Id: 0}, {Id: 1}}, RecoveryoperationCount: 2,
			}, valid: false,
		}, {
			desc: "duplicated recoveryoperation",
			genState: &types.GenesisState{
				RecoveryoperationList: []types.Recoveryoperation{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		}, {
			desc: "invalid recoveryoperation count",
			genState: &types.GenesisState{
				RecoveryoperationList: []types.Recoveryoperation{
					{
						Id: 1,
					},
				},
				RecoveryoperationCount: 0,
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
