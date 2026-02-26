package types

func NewMsgRecordRewardAccrual(creator string, address string, denom string, amount uint64, date string) *MsgRecordRewardAccrual {
	return &MsgRecordRewardAccrual{
		Creator: creator,
		Address: address,
		Denom:   denom,
		Amount:  amount,
		Date:    date,
	}
}
