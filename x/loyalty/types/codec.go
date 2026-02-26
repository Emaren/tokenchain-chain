package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRecordRewardAccrual{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaimReward{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMintVerifiedToken{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateRewardaccrual{},
		&MsgUpdateRewardaccrual{},
		&MsgDeleteRewardaccrual{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateVerifiedtoken{},
		&MsgUpdateVerifiedtoken{},
		&MsgDeleteVerifiedtoken{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCreatorallowlist{},
		&MsgUpdateCreatorallowlist{},
		&MsgDeleteCreatorallowlist{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
