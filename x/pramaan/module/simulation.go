// Package pramaan (this file) implements the module.AppModuleSimulation
// hooks used by the Cosmos SDK's randomized simulation testing framework.
package pramaan

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"pramaan/x/pramaan/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	pramaanGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pramaanGenesis)
}

// RegisterStoreDecoder registers a decoder used by the simulation framework
// to pretty-print raw KV-store diffs for this module. No-op.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the randomized message operations the
// simulator can fire against this module. Empty: this module has no
// simulate-able messages beyond UpdateParams.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	return operations
}

// ProposalMsgs returns the messages this module contributes to randomized
// governance-proposal simulation. Empty: none registered.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
