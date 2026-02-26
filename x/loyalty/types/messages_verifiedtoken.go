package types

func NewMsgCreateVerifiedtoken(
	creator string,
	denom string,
	issuer string,
	name string,
	symbol string,
	description string,
	website string,
	maxSupply uint64,
	mintedSupply uint64,
	verified bool,
	seizureOptIn bool,
	recoveryGroupPolicy string,
	recoveryTimelockHours uint64,

) *MsgCreateVerifiedtoken {
	return &MsgCreateVerifiedtoken{
		Creator:               creator,
		Denom:                 denom,
		Issuer:                issuer,
		Name:                  name,
		Symbol:                symbol,
		Description:           description,
		Website:               website,
		MaxSupply:             maxSupply,
		MintedSupply:          mintedSupply,
		Verified:              verified,
		SeizureOptIn:          seizureOptIn,
		RecoveryGroupPolicy:   recoveryGroupPolicy,
		RecoveryTimelockHours: recoveryTimelockHours,
	}
}

func NewMsgUpdateVerifiedtoken(
	creator string,
	denom string,
	issuer string,
	name string,
	symbol string,
	description string,
	website string,
	maxSupply uint64,
	mintedSupply uint64,
	verified bool,
	seizureOptIn bool,
	recoveryGroupPolicy string,
	recoveryTimelockHours uint64,

) *MsgUpdateVerifiedtoken {
	return &MsgUpdateVerifiedtoken{
		Creator:               creator,
		Denom:                 denom,
		Issuer:                issuer,
		Name:                  name,
		Symbol:                symbol,
		Description:           description,
		Website:               website,
		MaxSupply:             maxSupply,
		MintedSupply:          mintedSupply,
		Verified:              verified,
		SeizureOptIn:          seizureOptIn,
		RecoveryGroupPolicy:   recoveryGroupPolicy,
		RecoveryTimelockHours: recoveryTimelockHours,
	}
}

func NewMsgDeleteVerifiedtoken(
	creator string,
	denom string,

) *MsgDeleteVerifiedtoken {
	return &MsgDeleteVerifiedtoken{
		Creator: creator,
		Denom:   denom,
	}
}
