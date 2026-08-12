// Package authority (this file) wires x/authority into the app's depinject
// dependency graph: it declares what this module needs from the rest of the
// app (ModuleInputs) and what it provides back (ModuleOutputs), and
// constructs the module's Keeper/AppModule via ProvideModule.
package authority

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"pramaan/x/authority/keeper"
	"pramaan/x/authority/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

// init registers this module's config type and provider function with
// depinject at package-load time, so app_config.go's module list can
// reference &types.Module{} without any further wiring here.
func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

// ModuleInputs lists the dependencies depinject must supply to construct
// this module: config, storage, codecs, and the expected-keeper interfaces
// for auth/bank.
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
// graph: its own Keeper (so other modules can depend on
// authoritytypes.AuthorityKeeper) and the constructed AppModule.
type ModuleOutputs struct {
	depinject.Out

	AuthorityKeeper keeper.Keeper
	Module          appmodule.AppModule
}

// ProvideModule is the depinject provider for x/authority: it resolves the
// module's governance authority address, constructs the Keeper, and wraps
// it in an AppModule.
func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to the standard module-derived governance authority (matches every
	// other module in this app); override via Config.Authority if ever needed.
	authority := authtypes.NewModuleAddress(types.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.StoreService,
		in.Cdc,
		in.AddressCodec,
		authority,
		in.BankKeeper,
	)
	m := NewAppModule(k)

	return ModuleOutputs{AuthorityKeeper: k, Module: m}
}
