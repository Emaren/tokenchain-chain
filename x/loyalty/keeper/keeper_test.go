package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"tokenchain/x/loyalty/keeper"
	module "tokenchain/x/loyalty/module"
	"tokenchain/x/loyalty/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	bankKeeper   *mockBankKeeper
}

type mockBankKeeper struct {
	accountBalances map[string]sdk.Coins
	moduleBalances  map[string]sdk.Coins
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		accountBalances: make(map[string]sdk.Coins),
		moduleBalances:  make(map[string]sdk.Coins),
	}
}

func cloneCoins(in sdk.Coins) sdk.Coins {
	if len(in) == 0 {
		return sdk.NewCoins()
	}
	out := make(sdk.Coins, len(in))
	copy(out, in)
	return out.Sort()
}

func (m *mockBankKeeper) SpendableCoins(_ context.Context, addr sdk.AccAddress) sdk.Coins {
	return cloneCoins(m.accountBalances[addr.String()])
}

func (m *mockBankKeeper) MintCoins(_ context.Context, moduleName string, amt sdk.Coins) error {
	m.moduleBalances[moduleName] = m.moduleBalances[moduleName].Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(_ context.Context, moduleName string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	moduleBal := m.moduleBalances[moduleName]
	if !moduleBal.IsAllGTE(amt) {
		return sdkerrors.ErrInsufficientFunds
	}
	m.moduleBalances[moduleName] = moduleBal.Sub(amt...)
	m.accountBalances[recipientAddr.String()] = m.accountBalances[recipientAddr.String()].Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(_ context.Context, senderAddr sdk.AccAddress, moduleName string, amt sdk.Coins) error {
	accountBal := m.accountBalances[senderAddr.String()]
	if !accountBal.IsAllGTE(amt) {
		return sdkerrors.ErrInsufficientFunds
	}
	m.accountBalances[senderAddr.String()] = accountBal.Sub(amt...)
	m.moduleBalances[moduleName] = m.moduleBalances[moduleName].Add(amt...)
	return nil
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)
	bankKeeper := newMockBankKeeper()

	k := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		bankKeeper,
		nil,
		nil,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
		bankKeeper:   bankKeeper,
	}
}
