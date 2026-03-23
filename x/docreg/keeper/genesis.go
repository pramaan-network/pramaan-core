package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/docreg/types"
)

// InitGenesis initializes the module's state from genesis
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	// Set Params
	if err := k.ParamsStore.Set(ctx, genState.Params); err != nil {
		return fmt.Errorf("failed to set params: %w", err)
	}

	return nil
}

// ExportGenesis exports current state to genesis
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	params, err := k.ParamsStore.Get(sdkCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get params: %w", err)
	}

	genesis := &types.GenesisState{
		Params: params,
	}

	return genesis, nil
}
