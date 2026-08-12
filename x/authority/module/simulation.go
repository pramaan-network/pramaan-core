// Package authority (this file) implements the module.AppModuleSimulation
// hooks used by the Cosmos SDK's randomized simulation testing framework
// (`go test ./app -run TestFullAppSimulation ...`). None of this drives
// normal chain operation.
package authority

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"pramaan/x/authority/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}

	// ✅ correct way
	params := types.DefaultParams()

	authorityGenesis := types.GenesisState{
		Params: &params,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&authorityGenesis)
}

// RegisterStoreDecoder registers a decoder used by the simulation framework
// to pretty-print raw KV-store diffs for this module. No-op: this module has
// no custom decoding needs beyond the default.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the set of randomized message operations the
// simulator can fire against this module. Empty: no simulation operations
// have been implemented for x/authority yet, so the simulator will never
// exercise AddAuthority through this module (it can still be reached
// indirectly if other modules' operations call into it).
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return []simtypes.WeightedOperation{}
}

// ProposalMsgs returns the set of messages this module contributes to
// randomized governance-proposal simulation. Empty: no proposal-eligible
// messages are registered for x/authority.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
