// Package docreg (this file) implements the module.AppModuleSimulation
// hooks used by the Cosmos SDK's randomized simulation testing framework.
package docreg

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	docregsimulation "pramaan/x/docreg/simulation"
	"pramaan/x/docreg/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	docregGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&docregGenesis)
}

// RegisterStoreDecoder registers a decoder used by the simulation framework
// to pretty-print raw KV-store diffs for this module. No-op.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the randomized message operations (and their
// relative weights) the simulator can fire against this module: currently
// RegisterDocument and TransferDocument, both still stub NoOp
// implementations (see simulation/register_document.go and
// simulation/transfer_document.go).
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgRegisterDocument          = "op_weight_msg_docreg"
		defaultWeightMsgRegisterDocument int = 100
	)

	var weightMsgRegisterDocument int
	simState.AppParams.GetOrGenerate(opWeightMsgRegisterDocument, &weightMsgRegisterDocument, nil,
		func(_ *rand.Rand) {
			weightMsgRegisterDocument = defaultWeightMsgRegisterDocument
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRegisterDocument,
		docregsimulation.SimulateMsgRegisterDocument(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgTransferDocument          = "op_weight_msg_docreg"
		defaultWeightMsgTransferDocument int = 100
	)

	var weightMsgTransferDocument int
	simState.AppParams.GetOrGenerate(opWeightMsgTransferDocument, &weightMsgTransferDocument, nil,
		func(_ *rand.Rand) {
			weightMsgTransferDocument = defaultWeightMsgTransferDocument
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTransferDocument,
		docregsimulation.SimulateMsgTransferDocument(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns the messages this module contributes to randomized
// governance-proposal simulation. Empty: none registered.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
