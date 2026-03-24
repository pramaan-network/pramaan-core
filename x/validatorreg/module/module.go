package validatorreg

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"pramaan/x/validatorreg/keeper"
	"pramaan/x/validatorreg/types"
	authoritytypes "pramaan/x/authority/types"
)

var (
	_ module.AppModuleBasic = (*AppModule)(nil)
	_ module.AppModule      = (*AppModule)(nil)
	_ module.HasGenesis     = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
)

// 🔥 UPDATED STRUCT
type AppModule struct {
	cdc            codec.Codec
	keeper         keeper.Keeper
	authKeeper     types.AuthKeeper
	bankKeeper     types.BankKeeper
	pramaanKeeper  types.PramaanKeeper
	authorityKeeper authoritytypes.AuthorityKeeper
}

// 🔥 UPDATED CONSTRUCTOR
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	pramaanKeeper types.PramaanKeeper, 
	authorityKeeper authoritytypes.AuthorityKeeper,
) AppModule {
	return AppModule{
		cdc:            cdc,
		keeper:         keeper,
		authKeeper:     authKeeper,
		bankKeeper:     bankKeeper,
		pramaanKeeper:  pramaanKeeper,
		authorityKeeper: authorityKeeper,
	}
}

func (AppModule) IsAppModule() {}

func (AppModule) Name() string {
	return types.ModuleName
}

func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(clientCtx.CmdContext, mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModule) RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registrar)
}

// 🔥 FIXED REGISTER SERVICES
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {

	types.RegisterMsgServer(
		registrar,
		keeper.NewMsgServerImpl(am.keeper, am.pramaanKeeper, am.authorityKeeper),
	)

	types.RegisterQueryServer(registrar, keeper.NewQueryServerImpl(am.keeper))

	return nil
}

func (am AppModule) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	return am.cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (am AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := am.cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genState.Validate()
}

func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState

	if err := am.cdc.UnmarshalJSON(gs, &genState); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err))
	}

	if err := am.keeper.InitGenesis(ctx, genState); err != nil {
		panic(fmt.Errorf("failed to initialize %s genesis state: %w", types.ModuleName, err))
	}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	genState, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to export %s genesis state: %w", types.ModuleName, err))
	}

	bz, err := am.cdc.MarshalJSON(genState)
	if err != nil {
		panic(fmt.Errorf("failed to marshal %s genesis state: %w", types.ModuleName, err))
	}

	return bz
}

func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}