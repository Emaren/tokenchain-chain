package types

func NewMsgClaimReward(creator string, denom string) *MsgClaimReward {
	return &MsgClaimReward{
		Creator: creator,
		Denom:   denom,
	}
}
