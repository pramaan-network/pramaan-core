// Package simulation implements randomized-testing operations for x/docreg,
// used by the Cosmos SDK simulation framework (`go test ./app -run
// TestFullAppSimulation ...`). This file covers MsgRegisterDocument.
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

// SimulateMsgRegisterDocument builds a simulation operation for
// MsgRegisterDocument. Currently a stub: it always returns a NoOp — actual
// field population (hash/owner/issuer/doc_type) has not been implemented,
// so this message is never really exercised by the simulator.
func SimulateMsgRegisterDocument(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRegisterDocument{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the RegisterDocument simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "RegisterDocument simulation not implemented"), nil, nil
	}
}
