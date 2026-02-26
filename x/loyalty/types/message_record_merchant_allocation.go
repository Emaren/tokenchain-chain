package types

func NewMsgRecordMerchantAllocation(creator string, date string, denom string, activityScore uint64, bucketCAmount uint64) *MsgRecordMerchantAllocation {
	return &MsgRecordMerchantAllocation{
		Creator:       creator,
		Date:          date,
		Denom:         denom,
		ActivityScore: activityScore,
		BucketCAmount: bucketCAmount,
	}
}
