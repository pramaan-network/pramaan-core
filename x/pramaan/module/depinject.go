// Package pramaan (this file) wires x/pramaan into the app's depinject
// dependency graph.
package pramaan

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"pramaan/x/pramaan/keeper"
	"pramaan/x/pramaan/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

// init registers this module's config type and provider function with
// depinject at package-load time.
func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

// ModuleInputs lists the dependencies depinject must supply to construct
// this module: config, storage, codecs, and the expected-keeper interfaces
// for auth/bank. Unlike docreg/issuer/validatorreg, this module takes no
// cross-module keepers — nothing else in the app currently calls into
// x/pramaan's keeper (see the dead-subsystem note in types/keys.go).
type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper types.AuthKeeper
	BankKeeper types.BankKeeper
}

// ModuleOutputs lists what this module hands back into the dependency
// graph: its own Keeper and the constructed AppModule.
type ModuleOutputs struct {
	depinject.Out

	PramaanKeeper keeper.Keeper
	Module        appmodule.AppModule
}

// ProvideModule is the depinject provider for x/pramaan: it resolves the
// module's governance authority address, constructs the Keeper, and wraps
// it in an AppModule.
func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.AddressCodec,
		authority,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{PramaanKeeper: k, Module: m}
}
