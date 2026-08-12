package validatorreg

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/validatorreg/keeper"
	"pramaan/x/validatorreg/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
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

	// 🔥 YOUR AUTHORITY SYSTEM
	PramaanKeeper   types.PramaanKeeper
	AuthorityKeeper authoritytypes.AuthorityKeeper
}

type ModuleOutputs struct {
	depinject.Out

	ValidatorregKeeper keeper.Keeper
	Module             appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {

	// Default UpdateParams authority to the gov module account (matches every
	// other module in this chain), overridable via app_config.go. Previously
	// hardcoded to an empty byte slice, which made MsgUpdateParams permanently
	// unauthorizable — no address can ever equal an empty slice.
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

	// 🔥 PASS PRAMAAN KEEPER (REAL AUTHORITY)
	m := NewAppModule(
		in.Cdc,
		k,
		in.AuthKeeper,
		in.BankKeeper,
		in.PramaanKeeper,
		in.AuthorityKeeper,
	)

	return ModuleOutputs{
		ValidatorregKeeper: k,
		Module:             m,
	}
}
