package keeper

import (
	"context"

	"pramaan/x/docreg/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	return k.ParamsStore.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.ParamsStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}
