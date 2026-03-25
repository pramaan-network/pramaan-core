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

func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

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

type ModuleOutputs struct {
	depinject.Out

	DocregKeeper keeper.Keeper
	Module       appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {

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
