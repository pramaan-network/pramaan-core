// Package simulation (this file) implements the randomized-testing
// operation for MsgTransferDocument. See register_document.go for the
// package-level overview.
package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"pramaan/x/docreg/keeper"
	"pramaan/x/docreg/types"
)

// SimulateMsgTransferDocument builds a simulation operation for
// MsgTransferDocument. Currently a stub — always a NoOp, no real transfer
// scenario is exercised.
func SimulateMsgTransferDocument(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgTransferDocument{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the TransferDocument simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "TransferDocument simulation not implemented"), nil, nil
	}
}
