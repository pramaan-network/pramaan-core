// Package simulation implements randomized-testing operations for
// x/validatorreg, used by the Cosmos SDK simulation framework. This file
// covers MsgAddValidator.
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

// SimulateMsgAddValidator builds a simulation operation for
// MsgAddValidator. Currently a stub — always a NoOp.
func SimulateMsgAddValidator(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddValidator{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the AddValidator simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "AddValidator simulation not implemented"), nil, nil
	}
}
