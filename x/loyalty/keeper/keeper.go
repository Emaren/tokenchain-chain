package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"tokenchain/x/loyalty/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	bankKeeper           types.BankKeeper
	authKeeper           types.AuthKeeper
	stakingKeeper        types.StakingKeeper
	Creatorallowlist     collections.Map[string, types.Creatorallowlist]
	Verifiedtoken        collections.Map[string, types.Verifiedtoken]
	Rewardaccrual        collections.Map[string, types.Rewardaccrual]
	RecoveryoperationSeq collections.Sequence
	Recoveryoperation    collections.Map[uint64, types.Recoveryoperation]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	bankKeeper types.BankKeeper,
	authKeeper types.AuthKeeper,
	stakingKeeper types.StakingKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		bankKeeper:       bankKeeper,
		authKeeper:       authKeeper,
		stakingKeeper:    stakingKeeper,
		Params:           collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Creatorallowlist: collections.NewMap(sb, types.CreatorallowlistKey, "creatorallowlist", collections.StringKey, codec.CollValue[types.Creatorallowlist](cdc)), Verifiedtoken: collections.NewMap(sb, types.VerifiedtokenKey, "verifiedtoken", collections.StringKey, codec.CollValue[types.Verifiedtoken](cdc)), Rewardaccrual: collections.NewMap(sb, types.RewardaccrualKey, "rewardaccrual", collections.StringKey, codec.CollValue[types.Rewardaccrual](cdc)), Recoveryoperation: collections.NewMap(sb, types.RecoveryoperationKey, "recoveryoperation", collections.Uint64Key, codec.CollValue[types.Recoveryoperation](cdc)),
		RecoveryoperationSeq: collections.NewSequence(sb, types.RecoveryoperationCountKey, "recoveryoperationSequence"),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
