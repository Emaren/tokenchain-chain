package types

func NewMsgCancelRecoveryTransfer(creator string, id uint64, reason string) *MsgCancelRecoveryTransfer {
	return &MsgCancelRecoveryTransfer{
		Creator: creator,
		Id:      id,
		Reason:  reason,
	}
}
