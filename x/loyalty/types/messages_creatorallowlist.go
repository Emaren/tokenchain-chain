package types

func NewMsgCreateCreatorallowlist(
	creator string,
	address string,
	enabled bool,

) *MsgCreateCreatorallowlist {
	return &MsgCreateCreatorallowlist{
		Creator: creator,
		Address: address,
		Enabled: enabled,
	}
}

func NewMsgUpdateCreatorallowlist(
	creator string,
	address string,
	enabled bool,

) *MsgUpdateCreatorallowlist {
	return &MsgUpdateCreatorallowlist{
		Creator: creator,
		Address: address,
		Enabled: enabled,
	}
}

func NewMsgDeleteCreatorallowlist(
	creator string,
	address string,

) *MsgDeleteCreatorallowlist {
	return &MsgDeleteCreatorallowlist{
		Creator: creator,
		Address: address,
	}
}
