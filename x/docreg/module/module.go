// Package docreg wires the x/docreg module (the on-chain document/
// credential registry — SBT/NFT-style tokenized records tied to an
// issuer+owner) into the Cosmos SDK's module.AppModule lifecycle. The
// state-machine logic itself lives in x/docreg/keeper; this file is glue
// between that keeper and the SDK's module system.
package docreg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"pramaan/x/docreg/keeper"
	"pramaan/x/docreg/types"

	authoritytypes "pramaan/x/authority/types"
)

// AppModule implements module.AppModule for x/docreg, holding the live
// keeper plus its cross-module dependencies (auth/bank for standard
// account/balance access, authority/issuer for role and issuer-identity
// checks used by RegisterDocument).
type AppModule struct {
	cdc             codec.Codec
	keeper          keeper.Keeper
	authKeeper      types.AuthKeeper
	bankKeeper      types.BankKeeper
	authorityKeeper authoritytypes.AuthorityKeeper
	issuerKeeper    types.IssuerKeeper
}

// NewAppModule constructs an AppModule bound to the given keeper and its
// dependencies.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	authorityKeeper authoritytypes.AuthorityKeeper,
	issuerKeeper types.IssuerKeeper,
) AppModule {
	return AppModule{
		cdc:             cdc,
		keeper:          keeper,
		authKeeper:      authKeeper,
		bankKeeper:      bankKeeper,
		authorityKeeper: authorityKeeper,
		issuerKeeper:    issuerKeeper,
	}
}

// IsAppModule is a marker method satisfying cosmossdk.io/core/appmodule.AppModule.
func (AppModule) IsAppModule() {}

// Name returns the module's name, used as its key in the module manager and
// as the genesis JSON object key.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers this module's types with the legacy
// Amino codec. No-op here: docreg messages are proto-only.
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterGRPCGatewayRoutes registers this module's REST/gRPC-gateway
// routes so its queries are reachable over plain HTTP.
func (am AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterInterfaces registers this module's Msg types against the
// app-wide interface registry so they can be decoded out of transactions.
func (AppModule) RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registrar)
}

// RegisterServices registers this module's Msg and Query gRPC services,
// wiring in the concrete keeper plus its authority/issuer dependencies via
// NewMsgServerImpl.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(
		cfg.MsgServer(),
		keeper.NewMsgServerImpl(
			am.keeper,
			am.authorityKeeper,
			am.issuerKeeper,
		),
	)

	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServerImpl(am.keeper))
}

// DefaultGenesis returns this module's default genesis state as raw JSON.
func (am AppModule) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	return am.cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis decodes and validates a genesis JSON blob for this
// module, delegating the actual invariant checks to GenesisState.Validate.
func (am AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := am.cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// InitGenesis decodes the module's genesis JSON and loads it into the
// keeper's state (Params + Documents — see keeper/genesis.go) at chain
// start or on genesis import.
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	if err := am.cdc.UnmarshalJSON(gs, &genState); err != nil {
		panic(err)
	}
	if err := am.keeper.InitGenesis(ctx, genState); err != nil {
		panic(err)
	}
}

// ExportGenesis reads the module's current on-chain state back out as
// genesis JSON, for `pramaand export`.
func (am AppModule) ExportGenesis(ctx sdk.Context, _ codec.JSONCodec) json.RawMessage {
	genState, err := am.keeper.ExportGenesis(ctx)
	if err != nil {
		panic(err)
	}
	bz, err := am.cdc.MarshalJSON(genState)
	if err != nil {
		panic(err)
	}
	return bz
}

// ConsensusVersion reports this module's state-machine version, bumped on
// any consensus-breaking change. Still 1: no breaking change has shipped
// for this module yet.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock runs at the start of every block. No-op: no per-block logic.
func (am AppModule) BeginBlock(_ context.Context) error { return nil }

// EndBlock runs at the end of every block. No-op: no per-block logic.
func (am AppModule) EndBlock(_ context.Context) error { return nil }
