package types

func NewMsgCreateRewardaccrual(
	creator string,
	key string,
	address string,
	denom string,
	amount uint64,
	lastRollupDate string,

) *MsgCreateRewardaccrual {
	return &MsgCreateRewardaccrual{
		Creator:        creator,
		Key:            key,
		Address:        address,
		Denom:          denom,
		Amount:         amount,
		LastRollupDate: lastRollupDate,
	}
}

func NewMsgUpdateRewardaccrual(
	creator string,
	key string,
	address string,
	denom string,
	amount uint64,
	lastRollupDate string,

) *MsgUpdateRewardaccrual {
	return &MsgUpdateRewardaccrual{
		Creator:        creator,
		Key:            key,
		Address:        address,
		Denom:          denom,
		Amount:         amount,
		LastRollupDate: lastRollupDate,
	}
}

func NewMsgDeleteRewardaccrual(
	creator string,
	key string,

) *MsgDeleteRewardaccrual {
	return &MsgDeleteRewardaccrual{
		Creator: creator,
		Key:     key,
	}
}
