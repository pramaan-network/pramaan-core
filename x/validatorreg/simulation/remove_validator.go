// Package simulation (this file) implements the randomized-testing
// operation for MsgRemoveValidator. See add_validator.go for the
// package-level overview.
package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"pramaan/x/validatorreg/keeper"
	"pramaan/x/validatorreg/types"
)

// SimulateMsgRemoveValidator builds a simulation operation for
// MsgRemoveValidator. Currently a stub — always a NoOp.
func SimulateMsgRemoveValidator(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRemoveValidator{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the RemoveValidator simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "RemoveValidator simulation not implemented"), nil, nil
	}
}
