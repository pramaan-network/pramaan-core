// Package validatorreg wires the x/validatorreg module (the
// validator-admission workflow: apply -> AUTHORITY-quorum approve ->
// applicant-triggered activate, plus a direct ROOT-only add/remove path)
// into the Cosmos SDK's module.AppModule lifecycle. The state-machine logic
// itself lives in x/validatorreg/keeper; this file is glue between that
// keeper and the SDK's module system.
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

	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/validatorreg/keeper"
	"pramaan/x/validatorreg/types"
)

var (
	_ module.AppModuleBasic = (*AppModule)(nil)
	_ module.AppModule      = (*AppModule)(nil)
	_ module.HasGenesis     = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
)

// AppModule implements module.AppModule for x/validatorreg. Holds the live
// keeper plus its cross-module dependencies — authorityKeeper for the real
// role checks used throughout this module's message handlers, and
// pramaanKeeper, which is wired in but currently unused (see
// types/expected_keepers.go's note on PramaanKeeper).
type AppModule struct {
	cdc             codec.Codec
	keeper          keeper.Keeper
	authKeeper      types.AuthKeeper
	bankKeeper      types.BankKeeper
	pramaanKeeper   types.PramaanKeeper
	authorityKeeper authoritytypes.AuthorityKeeper
}

// NewAppModule constructs an AppModule bound to the given keeper and its
// dependencies.
func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	authKeeper types.AuthKeeper,
	bankKeeper types.BankKeeper,
	pramaanKeeper types.PramaanKeeper,
	authorityKeeper authoritytypes.AuthorityKeeper,
) AppModule {
	return AppModule{
		cdc:             cdc,
		keeper:          keeper,
		authKeeper:      authKeeper,
		bankKeeper:      bankKeeper,
		pramaanKeeper:   pramaanKeeper,
		authorityKeeper: authorityKeeper,
	}
}

// IsAppModule is a marker method satisfying cosmossdk.io/core/appmodule.AppModule.
func (AppModule) IsAppModule() {}

// Name returns the module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers this module's types with the legacy
// Amino codec. No-op: this module's messages are proto-only.
func (AppModule) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterGRPCGatewayRoutes registers this module's REST/gRPC-gateway
// routes so its queries are reachable over plain HTTP.
func (AppModule) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(clientCtx.CmdContext, mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// RegisterInterfaces registers this module's Msg types against the
// app-wide interface registry. See SECURITY_CHANGELOG.md #13:
// MsgApplyValidator/MsgApproveValidator/MsgActivateValidator were missing
// from types.RegisterInterfaces until that fix, making them permanently
// undecodable from any transaction.
func (AppModule) RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registrar)
}

// RegisterServices registers this module's Msg and Query gRPC services,
// wiring in the concrete keeper plus its pramaan/authority dependencies via
// NewMsgServerImpl.
func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {

	types.RegisterMsgServer(
		registrar,
		keeper.NewMsgServerImpl(am.keeper, am.pramaanKeeper, am.authorityKeeper),
	)

	types.RegisterQueryServer(registrar, keeper.NewQueryServerImpl(am.keeper))

	return nil
}

// DefaultGenesis returns this module's default genesis state as raw JSON.
func (am AppModule) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	return am.cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis decodes and validates a genesis JSON blob for this
// module.
func (am AppModule) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := am.cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genState.Validate()
}

// InitGenesis decodes the module's genesis JSON and loads it into the
// keeper's state (Params + Validators + Proposals — see keeper/genesis.go)
// at chain start or on genesis import.
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState

	if err := am.cdc.UnmarshalJSON(gs, &genState); err != nil {
		panic(fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err))
	}

	if err := am.keeper.InitGenesis(ctx, genState); err != nil {
		panic(fmt.Errorf("failed to initialize %s genesis state: %w", types.ModuleName, err))
	}
}

// ExportGenesis reads the module's current on-chain state back out as
// genesis JSON, for `pramaand export`.
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

// ConsensusVersion reports this module's state-machine version.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock runs at the start of every block. No-op.
func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// EndBlock runs at the end of every block. No-op: in particular, this
// module does NOT emit ABCI validator-set updates here — Active validators
// tracked by this module are not (yet) wired into CometBFT's actual
// consensus validator set. Consensus participation today is still driven
// by x/staking/gentx, independent of this module's Validators collection.
func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}
