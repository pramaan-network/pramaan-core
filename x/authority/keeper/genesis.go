package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/authority/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {

	// set params
	if genState.Params != nil {
		k.Params.Set(ctx, *genState.Params)
	}

	// 🔥 AUTO ROOT IF EMPTY
	if len(genState.Authorities) == 0 {

		root := types.Authority{
			Address: "pramaan1gmwkg6qhd03yjjvps4vj89a8j7peqf474ykd74",
			PubKey:  "genesis-root",
			Role:    "ROOT",
		}

		k.SetAuthority(ctx, root)
		return
	}

	// normal flow
	for _, auth := range genState.Authorities {
		k.SetAuthority(ctx, *auth)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	var authorities []*types.Authority

	// NOTE: you can improve later with iterator
	// for now safe empty export

	return &types.GenesisState{
		Params:      nil,
		Authorities: authorities,
	}
}
