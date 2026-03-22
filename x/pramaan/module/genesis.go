package pramaan

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/pramaan/keeper"
	"pramaan/x/pramaan/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// 🔥 STORE AUTHORITIES
	for _, auth := range genState.Authorities {
		err := k.Authorities.Set(ctx, auth.Address, *auth)
		if err != nil {
			panic(err)
		}
		fmt.Println("LOADED AUTHORITY:", auth.Address)
	}

	// 🔥 STORE THRESHOLD
	err := k.Threshold.Set(ctx, genState.Threshold)
	if err != nil {
		panic(err)
	}

	fmt.Println("THRESHOLD:", genState.Threshold)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	genesis := types.DefaultGenesis()
	return *genesis
}