package types

func NewMsgSetMerchantIncentiveRouting(
	creator string,
	denom string,
	merchantIncentiveStakersBps uint64,
	merchantIncentiveTreasuryBps uint64,
) *MsgSetMerchantIncentiveRouting {
	return &MsgSetMerchantIncentiveRouting{
		Creator:                      creator,
		Denom:                        denom,
		MerchantIncentiveStakersBps:  merchantIncentiveStakersBps,
		MerchantIncentiveTreasuryBps: merchantIncentiveTreasuryBps,
	}
}
