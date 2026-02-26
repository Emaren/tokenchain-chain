package app

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// registerWasmModules registers the wasm keeper and module.
func (app *App) registerWasmModules(appOpts types.AppOptions) error {
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(wasmtypes.StoreKey),
	); err != nil {
		return err
	}

	homePath := DefaultNodeHome
	if h, ok := appOpts.Get("home").(string); ok && h != "" {
		homePath = h
	}

	nodeConfig, err := wasm.ReadNodeConfig(appOpts)
	if err != nil {
		return fmt.Errorf("read wasm node config: %w", err)
	}

	govModuleAddr, err := app.AuthKeeper.AddressCodec().BytesToString(authtypes.NewModuleAddress(govtypes.ModuleName))
	if err != nil {
		return fmt.Errorf("resolve gov module address: %w", err)
	}

	app.WasmKeeper = wasmkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(wasmtypes.StoreKey)),
		app.AuthKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		distrkeeper.NewQuerier(app.DistrKeeper),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeperV2,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		homePath,
		nodeConfig,
		wasmtypes.VMConfig{},
		wasmkeeper.BuiltInCapabilities(),
		govModuleAddr,
	)

	return app.RegisterModules(
		wasm.NewAppModule(app.appCodec, &app.WasmKeeper, app.StakingKeeper, app.AuthKeeper, app.BankKeeper, app.MsgServiceRouter(), nil),
	)
}

// RegisterWasm manually registers wasm interfaces for client-side commands.
func RegisterWasm(cdc codec.Codec) module.AppModuleBasic {
	wasmBasic := wasm.AppModuleBasic{}
	wasmBasic.RegisterInterfaces(cdc.InterfaceRegistry())
	return module.CoreAppModuleBasicAdaptor(wasmtypes.ModuleName, wasmBasic)
}

// RegisterWasmModule returns a wasm app module for autocli wiring.
func RegisterWasmModule(cdc codec.Codec) appmodule.AppModule {
	wasmBasic := wasm.AppModuleBasic{}
	wasmBasic.RegisterInterfaces(cdc.InterfaceRegistry())
	return wasm.NewAppModule(cdc, nil, nil, nil, nil, nil, nil)
}
