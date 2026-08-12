// Package docreg (this file) wires x/docreg into the app's depinject
// dependency graph: it declares what this module needs from the rest of the
// app (ModuleInputs) and what it provides back (ModuleOutputs), and
// constructs the module's Keeper/AppModule via ProvideModule.
package docreg

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"pramaan/x/docreg/keeper"
	"pramaan/x/docreg/types"

	authoritytypes "pramaan/x/authority/types"
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
// for auth/bank plus the cross-module authority/issuer keepers this module
// needs for role and issuer-identity checks.
type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper types.AuthKeeper
	BankKeeper types.BankKeeper

	AuthorityKeeper authoritytypes.AuthorityKeeper

	IssuerKeeper types.IssuerKeeper
}

// ModuleOutputs lists what this module hands back into the dependency
// graph: its own Keeper and the constructed AppModule.
type ModuleOutputs struct {
	depinject.Out

	DocregKeeper keeper.Keeper
	Module       appmodule.AppModule
}

// ProvideModule is the depinject provider for x/docreg: it resolves the
// module's governance authority address, constructs the Keeper, and wraps
// it in an AppModule.
func ProvideModule(in ModuleInputs) ModuleOutputs {

	// Default UpdateParams authority to the gov module account (matches
	// issuer/pramaan/validatorreg). Previously defaulted to docreg's own
	// module address, which no proposal execution path can ever satisfy —
	// x/gov always signs as the "gov" module account, never as another
	// module's account — leaving MsgUpdateParams permanently unreachable.
	authority := authtypes.NewModuleAddress(types.GovModuleName)
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

	m := NewAppModule(
		in.Cdc,
		k,
		in.AuthKeeper,
		in.BankKeeper,
		in.AuthorityKeeper,
		in.IssuerKeeper, 
	)

	return ModuleOutputs{
		DocregKeeper: k,
		Module:       m,
	}
}
