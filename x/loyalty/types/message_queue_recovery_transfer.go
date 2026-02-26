package types

func NewMsgQueueRecoveryTransfer(creator string, denom string, fromAddress string, toAddress string, amount uint64) *MsgQueueRecoveryTransfer {
	return &MsgQueueRecoveryTransfer{
		Creator:     creator,
		Denom:       denom,
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Amount:      amount,
	}
}
