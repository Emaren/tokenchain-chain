package loyalty

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"tokenchain/testutil/sample"
	loyaltysimulation "tokenchain/x/loyalty/simulation"
	"tokenchain/x/loyalty/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	loyaltyGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		CreatorallowlistMap: []types.Creatorallowlist{{Creator: sample.AccAddress(),
			Address: "0",
		}, {Creator: sample.AccAddress(),
			Address: "1",
		}}, VerifiedtokenMap: []types.Verifiedtoken{{Creator: sample.AccAddress(),
			Denom: "0",
		}, {Creator: sample.AccAddress(),
			Denom: "1",
		}}, RewardaccrualMap: []types.Rewardaccrual{{Creator: sample.AccAddress(),
			Key: "0",
		}, {Creator: sample.AccAddress(),
			Key: "1",
		}}}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&loyaltyGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateCreatorallowlist          = "op_weight_msg_loyalty"
		defaultWeightMsgCreateCreatorallowlist int = 100
	)

	var weightMsgCreateCreatorallowlist int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCreatorallowlist, &weightMsgCreateCreatorallowlist, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCreatorallowlist = defaultWeightMsgCreateCreatorallowlist
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCreatorallowlist,
		loyaltysimulation.SimulateMsgCreateCreatorallowlist(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateCreatorallowlist          = "op_weight_msg_loyalty"
		defaultWeightMsgUpdateCreatorallowlist int = 100
	)

	var weightMsgUpdateCreatorallowlist int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCreatorallowlist, &weightMsgUpdateCreatorallowlist, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCreatorallowlist = defaultWeightMsgUpdateCreatorallowlist
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCreatorallowlist,
		loyaltysimulation.SimulateMsgUpdateCreatorallowlist(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteCreatorallowlist          = "op_weight_msg_loyalty"
		defaultWeightMsgDeleteCreatorallowlist int = 100
	)

	var weightMsgDeleteCreatorallowlist int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCreatorallowlist, &weightMsgDeleteCreatorallowlist, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCreatorallowlist = defaultWeightMsgDeleteCreatorallowlist
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCreatorallowlist,
		loyaltysimulation.SimulateMsgDeleteCreatorallowlist(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateVerifiedtoken          = "op_weight_msg_loyalty"
		defaultWeightMsgCreateVerifiedtoken int = 100
	)

	var weightMsgCreateVerifiedtoken int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateVerifiedtoken, &weightMsgCreateVerifiedtoken, nil,
		func(_ *rand.Rand) {
			weightMsgCreateVerifiedtoken = defaultWeightMsgCreateVerifiedtoken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateVerifiedtoken,
		loyaltysimulation.SimulateMsgCreateVerifiedtoken(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateVerifiedtoken          = "op_weight_msg_loyalty"
		defaultWeightMsgUpdateVerifiedtoken int = 100
	)

	var weightMsgUpdateVerifiedtoken int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateVerifiedtoken, &weightMsgUpdateVerifiedtoken, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateVerifiedtoken = defaultWeightMsgUpdateVerifiedtoken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateVerifiedtoken,
		loyaltysimulation.SimulateMsgUpdateVerifiedtoken(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteVerifiedtoken          = "op_weight_msg_loyalty"
		defaultWeightMsgDeleteVerifiedtoken int = 100
	)

	var weightMsgDeleteVerifiedtoken int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteVerifiedtoken, &weightMsgDeleteVerifiedtoken, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteVerifiedtoken = defaultWeightMsgDeleteVerifiedtoken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteVerifiedtoken,
		loyaltysimulation.SimulateMsgDeleteVerifiedtoken(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCreateRewardaccrual          = "op_weight_msg_loyalty"
		defaultWeightMsgCreateRewardaccrual int = 100
	)

	var weightMsgCreateRewardaccrual int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateRewardaccrual, &weightMsgCreateRewardaccrual, nil,
		func(_ *rand.Rand) {
			weightMsgCreateRewardaccrual = defaultWeightMsgCreateRewardaccrual
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateRewardaccrual,
		loyaltysimulation.SimulateMsgCreateRewardaccrual(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateRewardaccrual          = "op_weight_msg_loyalty"
		defaultWeightMsgUpdateRewardaccrual int = 100
	)

	var weightMsgUpdateRewardaccrual int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateRewardaccrual, &weightMsgUpdateRewardaccrual, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateRewardaccrual = defaultWeightMsgUpdateRewardaccrual
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateRewardaccrual,
		loyaltysimulation.SimulateMsgUpdateRewardaccrual(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteRewardaccrual          = "op_weight_msg_loyalty"
		defaultWeightMsgDeleteRewardaccrual int = 100
	)

	var weightMsgDeleteRewardaccrual int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteRewardaccrual, &weightMsgDeleteRewardaccrual, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteRewardaccrual = defaultWeightMsgDeleteRewardaccrual
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteRewardaccrual,
		loyaltysimulation.SimulateMsgDeleteRewardaccrual(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgMintVerifiedToken          = "op_weight_msg_loyalty"
		defaultWeightMsgMintVerifiedToken int = 100
	)

	var weightMsgMintVerifiedToken int
	simState.AppParams.GetOrGenerate(opWeightMsgMintVerifiedToken, &weightMsgMintVerifiedToken, nil,
		func(_ *rand.Rand) {
			weightMsgMintVerifiedToken = defaultWeightMsgMintVerifiedToken
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgMintVerifiedToken,
		loyaltysimulation.SimulateMsgMintVerifiedToken(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgClaimReward          = "op_weight_msg_loyalty"
		defaultWeightMsgClaimReward int = 100
	)

	var weightMsgClaimReward int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimReward, &weightMsgClaimReward, nil,
		func(_ *rand.Rand) {
			weightMsgClaimReward = defaultWeightMsgClaimReward
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimReward,
		loyaltysimulation.SimulateMsgClaimReward(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRecordRewardAccrual          = "op_weight_msg_loyalty"
		defaultWeightMsgRecordRewardAccrual int = 100
	)

	var weightMsgRecordRewardAccrual int
	simState.AppParams.GetOrGenerate(opWeightMsgRecordRewardAccrual, &weightMsgRecordRewardAccrual, nil,
		func(_ *rand.Rand) {
			weightMsgRecordRewardAccrual = defaultWeightMsgRecordRewardAccrual
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRecordRewardAccrual,
		loyaltysimulation.SimulateMsgRecordRewardAccrual(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgQueueRecoveryTransfer          = "op_weight_msg_loyalty"
		defaultWeightMsgQueueRecoveryTransfer int = 100
	)

	var weightMsgQueueRecoveryTransfer int
	simState.AppParams.GetOrGenerate(opWeightMsgQueueRecoveryTransfer, &weightMsgQueueRecoveryTransfer, nil,
		func(_ *rand.Rand) {
			weightMsgQueueRecoveryTransfer = defaultWeightMsgQueueRecoveryTransfer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgQueueRecoveryTransfer,
		loyaltysimulation.SimulateMsgQueueRecoveryTransfer(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgExecuteRecoveryTransfer          = "op_weight_msg_loyalty"
		defaultWeightMsgExecuteRecoveryTransfer int = 100
	)

	var weightMsgExecuteRecoveryTransfer int
	simState.AppParams.GetOrGenerate(opWeightMsgExecuteRecoveryTransfer, &weightMsgExecuteRecoveryTransfer, nil,
		func(_ *rand.Rand) {
			weightMsgExecuteRecoveryTransfer = defaultWeightMsgExecuteRecoveryTransfer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExecuteRecoveryTransfer,
		loyaltysimulation.SimulateMsgExecuteRecoveryTransfer(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgCancelRecoveryTransfer          = "op_weight_msg_loyalty"
		defaultWeightMsgCancelRecoveryTransfer int = 100
	)

	var weightMsgCancelRecoveryTransfer int
	simState.AppParams.GetOrGenerate(opWeightMsgCancelRecoveryTransfer, &weightMsgCancelRecoveryTransfer, nil,
		func(_ *rand.Rand) {
			weightMsgCancelRecoveryTransfer = defaultWeightMsgCancelRecoveryTransfer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCancelRecoveryTransfer,
		loyaltysimulation.SimulateMsgCancelRecoveryTransfer(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
