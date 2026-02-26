package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"tokenchain/testutil/sample"
	"tokenchain/x/loyalty/keeper"
	"tokenchain/x/loyalty/types"
)

func createRecoveryEnabledToken(t *testing.T, f *fixture, srv types.MsgServer, ctx sdk.Context, creator string, subdenom string) string {
	t.Helper()

	f.groupKeeper.addPolicy(creator)
	msg := baseVerifiedToken(creator, subdenom)
	msg.SeizureOptIn = true
	msg.RecoveryGroupPolicy = creator
	msg.RecoveryTimelockHours = 1
	_, err := srv.CreateVerifiedtoken(ctx, msg)
	require.NoError(t, err)

	return factoryDenom(creator, subdenom)
}

func TestQueueRecoveryTransfer(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	ctx := sdk.UnwrapSDKContext(f.ctx).WithBlockTime(time.Unix(1700000000, 0))

	denom := createRecoveryEnabledToken(t, f, srv, ctx, creator, "recoverqueue")
	from := sample.AccAddress()
	to := sample.AccAddress()

	_, err := srv.QueueRecoveryTransfer(ctx, &types.MsgQueueRecoveryTransfer{
		Creator:     sample.AccAddress(),
		Denom:       denom,
		FromAddress: from,
		ToAddress:   to,
		Amount:      25,
	})
	require.ErrorIs(t, err, types.ErrRecoveryUnauthorized)

	queueResp, err := srv.QueueRecoveryTransfer(ctx, &types.MsgQueueRecoveryTransfer{
		Creator:     creator,
		Denom:       denom,
		FromAddress: from,
		ToAddress:   to,
		Amount:      25,
	})
	require.NoError(t, err)

	lastID, err := f.keeper.RecoveryoperationSeq.Peek(ctx)
	require.NoError(t, err)
	require.NotZero(t, lastID)
	opID := lastID - 1
	require.EqualValues(t, opID, queueResp.Id)
	require.Equal(t, types.RecoveryStatusQueued, queueResp.Status)

	op, err := f.keeper.Recoveryoperation.Get(ctx, opID)
	require.NoError(t, err)
	require.Equal(t, types.RecoveryStatusQueued, op.Status)
	require.Equal(t, creator, op.RequestedBy)
	require.EqualValues(t, 1700000000+3600, op.ExecuteAfter)
}

func TestExecuteRecoveryTransfer(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	baseCtx := sdk.UnwrapSDKContext(f.ctx).WithBlockTime(time.Unix(1700001000, 0))

	denom := createRecoveryEnabledToken(t, f, srv, baseCtx, creator, "recoverexec")
	from := sample.AccAddress()
	to := sample.AccAddress()

	mintResp, err := srv.MintVerifiedToken(baseCtx, &types.MsgMintVerifiedToken{
		Creator:   creator,
		Denom:     denom,
		Recipient: from,
		Amount:    50,
	})
	require.NoError(t, err)
	require.Equal(t, denom, mintResp.Denom)
	require.EqualValues(t, 50, mintResp.MintedSupply)

	queueResp, err := srv.QueueRecoveryTransfer(baseCtx, &types.MsgQueueRecoveryTransfer{
		Creator:     creator,
		Denom:       denom,
		FromAddress: from,
		ToAddress:   to,
		Amount:      50,
	})
	require.NoError(t, err)

	opID, err := f.keeper.RecoveryoperationSeq.Peek(baseCtx)
	require.NoError(t, err)
	opID--
	require.EqualValues(t, opID, queueResp.Id)

	_, err = srv.ExecuteRecoveryTransfer(baseCtx, &types.MsgExecuteRecoveryTransfer{
		Creator: creator,
		Id:      opID,
	})
	require.ErrorIs(t, err, types.ErrRecoveryTooEarly)

	execCtx := baseCtx.WithBlockTime(time.Unix(1700001000+3601, 0))
	_, err = srv.ExecuteRecoveryTransfer(execCtx, &types.MsgExecuteRecoveryTransfer{
		Creator: sample.AccAddress(),
		Id:      opID,
	})
	require.ErrorIs(t, err, types.ErrRecoveryUnauthorized)

	executeResp, err := srv.ExecuteRecoveryTransfer(execCtx, &types.MsgExecuteRecoveryTransfer{
		Creator: creator,
		Id:      opID,
	})
	require.NoError(t, err)
	require.EqualValues(t, opID, executeResp.Id)
	require.Equal(t, types.RecoveryStatusExecuted, executeResp.Status)
	require.EqualValues(t, 1700001000+3601, executeResp.ExecutedAt)

	op, err := f.keeper.Recoveryoperation.Get(execCtx, opID)
	require.NoError(t, err)
	require.Equal(t, types.RecoveryStatusExecuted, op.Status)
	require.EqualValues(t, 1700001000+3601, op.ExecutedAt)

	fromAddrBz, err := f.addressCodec.StringToBytes(from)
	require.NoError(t, err)
	toAddrBz, err := f.addressCodec.StringToBytes(to)
	require.NoError(t, err)

	fromBal := f.bankKeeper.SpendableCoins(execCtx, fromAddrBz).AmountOf(denom)
	toBal := f.bankKeeper.SpendableCoins(execCtx, toAddrBz).AmountOf(denom)
	require.True(t, fromBal.IsZero())
	require.EqualValues(t, 50, toBal.Uint64())
}

func TestCancelRecoveryTransfer(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	creator := authorityAddress(t, f)
	baseCtx := sdk.UnwrapSDKContext(f.ctx).WithBlockTime(time.Unix(1700002000, 0))

	denom := createRecoveryEnabledToken(t, f, srv, baseCtx, creator, "recovercancel")
	from := sample.AccAddress()
	to := sample.AccAddress()

	queueResp, err := srv.QueueRecoveryTransfer(baseCtx, &types.MsgQueueRecoveryTransfer{
		Creator:     creator,
		Denom:       denom,
		FromAddress: from,
		ToAddress:   to,
		Amount:      12,
	})
	require.NoError(t, err)

	opID, err := f.keeper.RecoveryoperationSeq.Peek(baseCtx)
	require.NoError(t, err)
	opID--
	require.EqualValues(t, opID, queueResp.Id)

	_, err = srv.CancelRecoveryTransfer(baseCtx, &types.MsgCancelRecoveryTransfer{
		Creator: sample.AccAddress(),
		Id:      opID,
		Reason:  "unauthorized",
	})
	require.ErrorIs(t, err, types.ErrRecoveryUnauthorized)

	cancelResp, err := srv.CancelRecoveryTransfer(baseCtx, &types.MsgCancelRecoveryTransfer{
		Creator: creator,
		Id:      opID,
		Reason:  "customer support request",
	})
	require.NoError(t, err)
	require.EqualValues(t, opID, cancelResp.Id)
	require.Equal(t, types.RecoveryStatusCancelled, cancelResp.Status)
	require.EqualValues(t, 1700002000, cancelResp.CancelledAt)

	op, err := f.keeper.Recoveryoperation.Get(baseCtx, opID)
	require.NoError(t, err)
	require.Equal(t, types.RecoveryStatusCancelled, op.Status)
	require.Equal(t, "customer support request", op.CancelReason)

	_, err = srv.ExecuteRecoveryTransfer(baseCtx.WithBlockTime(time.Unix(1700002000+3601, 0)), &types.MsgExecuteRecoveryTransfer{
		Creator: creator,
		Id:      opID,
	})
	require.ErrorIs(t, err, types.ErrRecoveryNotQueued)
	require.NotErrorIs(t, err, sdkerrors.ErrUnauthorized)
}
