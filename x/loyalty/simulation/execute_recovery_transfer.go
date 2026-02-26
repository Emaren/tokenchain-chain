package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func SimulateMsgExecuteRecoveryTransfer(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgExecuteRecoveryTransfer{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the ExecuteRecoveryTransfer simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "ExecuteRecoveryTransfer simulation not implemented"), nil, nil
	}
}
