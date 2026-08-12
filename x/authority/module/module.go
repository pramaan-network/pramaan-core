// Package authority wires the x/authority module into the Cosmos SDK's
// module.AppModule / module.AppModuleBasic lifecycle: interface/codec
// registration, genesis (de)serialization, and gRPC service registration.
// The actual state-machine logic lives in x/authority/keeper; this file is
// glue between that keeper and the SDK's module system.
package authority

import (
	"context"
	"encoding/json"

	"pramaan/x/authority/keeper"
	"pramaan/x/authority/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// --------------------
// AppModuleBasic
// --------------------

// AppModuleBasic implements the module.AppModuleBasic interface: the subset
// of module behavior that doesn't need a live keeper (name, codec/interface
// registration, default/validated genesis, CLI/gateway wiring).
type AppModuleBasic struct{}

// Name returns the module's name, used as its key in the module manager and
// as the genesis JSON object key.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterInterfaces registers this module's proto.Message implementations
// (its Msg types) against the app-wide interface registry, so they can be
// decoded out of transactions. Skipping this for a message type makes that
// message permanently undecodable — see SECURITY_CHANGELOG.md bugs #12/#13
// for what happens when this step is missed in a sibling module.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// RegisterLegacyAminoCodec registers this module's types with the legacy
// Amino codec, kept for backwards compatibility with pre-protobuf tooling
// (e.g. the Ledger hardware wallet signing flow).
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis returns this module's default genesis state as raw JSON,
// used when no explicit genesis is supplied (e.g. `pramaand init`).
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis checks that a genesis JSON blob unmarshals cleanly for
// this module. Note: this only confirms the JSON decodes — it does not call
// GenesisState.Validate(), so a structurally valid but semantically invalid
// genesis (e.g. missing ROOT authority) will pass this check and only fail
// later in InitGenesis. Worth tightening if this module's genesis needs are
// ever hand-edited rather than exported from a running chain.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}
	return nil
}

// --------------------
// AppModule
// --------------------

// AppModule implements the full module.AppModule interface, adding the live
// keeper needed for genesis and gRPC service registration on top of
// AppModuleBasic.
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule constructs an AppModule bound to the given keeper.
func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

// Name returns the module's name (see AppModuleBasic.Name).
func (am AppModule) Name() string {
	return types.ModuleName
}

// RegisterServices registers this module's Msg and Query gRPC services with
// the module manager's service configurator, wiring in the concrete keeper
// via NewMsgServerImpl/NewQueryServerImpl.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// InitGenesis decodes the module's genesis JSON and loads it into the
// keeper's state at chain start (or on `pramaand init`/import-genesis).
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, genesisState)
}

// ExportGenesis reads the module's current on-chain state back out as
// genesis JSON — used by `pramaand export` to produce a new genesis file
// (e.g. for a chain upgrade or fork).
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesis := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(genesis)
}

// BeginBlock runs at the start of every block. No-op: this module has no
// per-block logic.
func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// EndBlock runs at the end of every block. No-op: this module has no
// per-block logic (in particular, it does not emit ABCI validator-set
// updates — that responsibility belongs to x/validatorreg/x/staking).
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

// IsAppModule is a marker method satisfying cosmossdk.io/core/appmodule.AppModule.
func (AppModule) IsAppModule() {}

// RegisterGRPCGatewayRoutes registers this module's REST/gRPC-gateway routes
// so its queries are reachable over plain HTTP, not just gRPC.
func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}
