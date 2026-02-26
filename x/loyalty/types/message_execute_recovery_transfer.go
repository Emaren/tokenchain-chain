package types

func NewMsgExecuteRecoveryTransfer(creator string, id uint64) *MsgExecuteRecoveryTransfer {
	return &MsgExecuteRecoveryTransfer{
		Creator: creator,
		Id:      id,
	}
}
