// package keeper

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"pramaan/x/authority/types"
// )

// func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {

// 	// set params
// 	if genState.Params != nil {
// 		k.Params.Set(ctx, *genState.Params)
// 	}

// 	// 🔐 NO AUTO ROOT — MUST COME FROM GENESIS
// 	for _, auth := range genState.Authorities {
// 		k.SetAuthority(ctx, *auth)
// 	}
// }

// func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

// 	var authorities []*types.Authority

// 	// TODO: (optional later) export authorities properly

// 	return &types.GenesisState{
// 		Params:      nil,
// 		Authorities: authorities,
// 	}
// }


package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/authority/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {

	if genState.Params != nil {
		k.Params.Set(ctx, *genState.Params)
	}

	if len(genState.Authorities) == 0 {

	root := types.Authority{
		Address: "pramaan1y7y6ym9j6y4kygrnwng0cfseasy5r72tf4qyuz",
		PubKey:  "root-multisig",
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
