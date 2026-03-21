package validatorreg

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	validatorregsimulation "pramaan/x/validatorreg/simulation"
	"pramaan/x/validatorreg/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	validatorregGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&validatorregGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgAddValidator          = "op_weight_msg_validatorreg"
		defaultWeightMsgAddValidator int = 100
	)

	var weightMsgAddValidator int
	simState.AppParams.GetOrGenerate(opWeightMsgAddValidator, &weightMsgAddValidator, nil,
		func(_ *rand.Rand) {
			weightMsgAddValidator = defaultWeightMsgAddValidator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddValidator,
		validatorregsimulation.SimulateMsgAddValidator(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgRemoveValidator          = "op_weight_msg_validatorreg"
		defaultWeightMsgRemoveValidator int = 100
	)

	var weightMsgRemoveValidator int
	simState.AppParams.GetOrGenerate(opWeightMsgRemoveValidator, &weightMsgRemoveValidator, nil,
		func(_ *rand.Rand) {
			weightMsgRemoveValidator = defaultWeightMsgRemoveValidator
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRemoveValidator,
		validatorregsimulation.SimulateMsgRemoveValidator(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
