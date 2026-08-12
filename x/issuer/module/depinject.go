// Package issuer (this file) wires x/issuer into the app's depinject
// dependency graph.
package issuer

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"pramaan/x/issuer/keeper"
	"pramaan/x/issuer/types"

	authoritykeeper "pramaan/x/authority/keeper"
	validatorkeeper "pramaan/x/validatorreg/keeper"
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
// this module: config, storage, codecs, and the cross-module
// authority/validator keepers CreateIssuer needs for its role/domain checks
// (note: these are concrete keeper types here, not the narrower
// types.AuthorityKeeper/types.ValidatorKeeper interfaces — both work since
// the concrete keepers satisfy those interfaces, but it does mean this
// module's depinject wiring is coupled to authority/validatorreg's keeper
// packages directly rather than only their types packages).
type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper types.AuthKeeper
	BankKeeper types.BankKeeper

	AuthorityKeeper authoritykeeper.Keeper
	ValidatorKeeper validatorkeeper.Keeper
}

// ModuleOutputs lists what this module hands back into the dependency
// graph: its own Keeper and the constructed AppModule.
type ModuleOutputs struct {
	depinject.Out

	IssuerKeeper keeper.Keeper
	Module       appmodule.AppModule
}

// ProvideModule is the depinject provider for x/issuer: it resolves the
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
		in.AuthorityKeeper,
		in.ValidatorKeeper,
	)

	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{IssuerKeeper: k, Module: m}
}
