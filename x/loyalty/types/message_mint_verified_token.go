package types

func NewMsgMintVerifiedToken(creator string, denom string, recipient string, amount uint64) *MsgMintVerifiedToken {
	return &MsgMintVerifiedToken{
		Creator:   creator,
		Denom:     denom,
		Recipient: recipient,
		Amount:    amount,
	}
}
