package loyalty

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"tokenchain/x/loyalty/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListCreatorallowlist",
					Use:       "list-creatorallowlist",
					Short:     "List all creatorallowlist",
				},
				{
					RpcMethod:      "GetCreatorallowlist",
					Use:            "get-creatorallowlist [id]",
					Short:          "Gets a creatorallowlist",
					Alias:          []string{"show-creatorallowlist"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod: "ListVerifiedtoken",
					Use:       "list-verifiedtoken",
					Short:     "List all verifiedtoken",
				},
				{
					RpcMethod:      "GetVerifiedtoken",
					Use:            "get-verifiedtoken [id]",
					Short:          "Gets a verifiedtoken",
					Alias:          []string{"show-verifiedtoken"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "GetVerifiedtokenByDenom",
					Use:            "get-verifiedtoken-by-denom [denom]",
					Short:          "Gets a verifiedtoken by denom query parameter",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod: "ListRewardaccrual",
					Use:       "list-rewardaccrual",
					Short:     "List all rewardaccrual",
				},
				{
					RpcMethod:      "GetRewardaccrual",
					Use:            "get-rewardaccrual [id]",
					Short:          "Gets a rewardaccrual",
					Alias:          []string{"show-rewardaccrual"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "key"}},
				},
				{
					RpcMethod: "FilterRewardaccrual",
					Use:       "filter-rewardaccrual",
					Short:     "Filter reward accrual records by address/denom",
				},
				{
					RpcMethod: "ListRecoveryoperation",
					Use:       "list-recoveryoperation",
					Short:     "List all recoveryoperation",
				},
				{
					RpcMethod:      "GetRecoveryoperation",
					Use:            "get-recoveryoperation [id]",
					Short:          "Gets a recoveryoperation by id",
					Alias:          []string{"show-recoveryoperation"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod: "FilterRecoveryoperation",
					Use:       "filter-recoveryoperation",
					Short:     "Filter recovery operations by status/token/address",
				},
				{
					RpcMethod: "DailyRollupStatus",
					Use:       "daily-rollup-status",
					Short:     "Show daily rollup status for configured timezone",
				},
				{
					RpcMethod: "RewardPoolBalance",
					Use:       "reward-pool-balance [denom]",
					Short:     "Show loyalty reward pool spendable balance for denom",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "denom"},
					},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateCreatorallowlist",
					Use:            "create-creatorallowlist [address] [enabled]",
					Short:          "Create a new creatorallowlist",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "enabled"}},
				},
				{
					RpcMethod:      "UpdateCreatorallowlist",
					Use:            "update-creatorallowlist [address] [enabled]",
					Short:          "Update creatorallowlist",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "enabled"}},
				},
				{
					RpcMethod:      "DeleteCreatorallowlist",
					Use:            "delete-creatorallowlist [address]",
					Short:          "Delete creatorallowlist",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}},
				},
				{
					RpcMethod:      "CreateVerifiedtoken",
					Use:            "create-verifiedtoken [denom] [issuer] [name] [symbol] [description] [website] [max-supply] [minted-supply] [verified] [seizure-opt-in] [recovery-group-policy] [recovery-timelock-hours]",
					Short:          "Create a new verifiedtoken",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "issuer"}, {ProtoField: "name"}, {ProtoField: "symbol"}, {ProtoField: "description"}, {ProtoField: "website"}, {ProtoField: "max_supply"}, {ProtoField: "minted_supply"}, {ProtoField: "verified"}, {ProtoField: "seizure_opt_in"}, {ProtoField: "recovery_group_policy"}, {ProtoField: "recovery_timelock_hours"}},
				},
				{
					RpcMethod:      "UpdateVerifiedtoken",
					Use:            "update-verifiedtoken [denom] [issuer] [name] [symbol] [description] [website] [max-supply] [minted-supply] [verified] [seizure-opt-in] [recovery-group-policy] [recovery-timelock-hours]",
					Short:          "Update verifiedtoken",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "issuer"}, {ProtoField: "name"}, {ProtoField: "symbol"}, {ProtoField: "description"}, {ProtoField: "website"}, {ProtoField: "max_supply"}, {ProtoField: "minted_supply"}, {ProtoField: "verified"}, {ProtoField: "seizure_opt_in"}, {ProtoField: "recovery_group_policy"}, {ProtoField: "recovery_timelock_hours"}},
				},
				{
					RpcMethod:      "RenounceTokenAdmin",
					Use:            "renounce-token-admin [denom]",
					Short:          "Renounce verified token admin powers (locks minting)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "DeleteVerifiedtoken",
					Use:            "delete-verifiedtoken [denom]",
					Short:          "Delete verifiedtoken",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "CreateRewardaccrual",
					Use:            "create-rewardaccrual [key] [address] [denom] [amount] [last-rollup-date]",
					Short:          "Create a new rewardaccrual",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "key"}, {ProtoField: "address"}, {ProtoField: "denom"}, {ProtoField: "amount"}, {ProtoField: "last_rollup_date"}},
				},
				{
					RpcMethod:      "UpdateRewardaccrual",
					Use:            "update-rewardaccrual [key] [address] [denom] [amount] [last-rollup-date]",
					Short:          "Update rewardaccrual",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "key"}, {ProtoField: "address"}, {ProtoField: "denom"}, {ProtoField: "amount"}, {ProtoField: "last_rollup_date"}},
				},
				{
					RpcMethod:      "DeleteRewardaccrual",
					Use:            "delete-rewardaccrual [key]",
					Short:          "Delete rewardaccrual",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "key"}},
				},
				{
					RpcMethod:      "MintVerifiedToken",
					Use:            "mint-verified-token [denom] [recipient] [amount]",
					Short:          "Send a mint-verified-token tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "recipient"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "ClaimReward",
					Use:            "claim-reward [denom]",
					Short:          "Send a claim-reward tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}},
				},
				{
					RpcMethod:      "FundRewardPool",
					Use:            "fund-reward-pool [denom] [amount]",
					Short:          "Fund the loyalty reward pool from signer account",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "RecordRewardAccrual",
					Use:            "record-reward-accrual [address] [denom] [amount] [date]",
					Short:          "Send a record-reward-accrual tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "denom"}, {ProtoField: "amount"}, {ProtoField: "date"}},
				},
				{
					RpcMethod:      "QueueRecoveryTransfer",
					Use:            "queue-recovery-transfer [denom] [from-address] [to-address] [amount]",
					Short:          "Send a queue-recovery-transfer tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "denom"}, {ProtoField: "from_address"}, {ProtoField: "to_address"}, {ProtoField: "amount"}},
				},
				{
					RpcMethod:      "ExecuteRecoveryTransfer",
					Use:            "execute-recovery-transfer [id]",
					Short:          "Send a execute-recovery-transfer tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				{
					RpcMethod:      "CancelRecoveryTransfer",
					Use:            "cancel-recovery-transfer [id] [reason]",
					Short:          "Send a cancel-recovery-transfer tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "reason"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
