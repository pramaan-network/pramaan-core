package keeper

import (
	"pramaan/x/authority/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from genesis data
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	// Set params (pointer → value)
	if genState.Params != nil {
		k.Params.Set(ctx, *genState.Params)
	}

	// If no authorities → auto root
	if len(genState.Authorities) == 0 {
		addrStr, err := k.addressCodec.BytesToString(k.authority)
		if err != nil {
			panic(err)
		}

		root := types.Authority{
			Address: addrStr,
			PubKey:  "genesis-root",
			Role:    "ROOT",
		}

		k.SetAuthority(ctx, root)
		return
	}

	// Load authorities (pointer → value)
	for _, authority := range genState.Authorities {
		if authority == nil {
			continue
		}
		k.SetAuthority(ctx, *authority)
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()

	params, _ := k.Params.Get(ctx)
	genesis.Params = &params

	return genesis
}