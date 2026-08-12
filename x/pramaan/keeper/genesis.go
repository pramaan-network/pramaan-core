package keeper

import (
	"context"

	"pramaan/x/pramaan/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
//
// NOTE: genState.Authorities/Threshold were previously silently dropped here —
// the proto/genesis.pb.go has always supported these fields (and a dead,
// never-called free function in module/genesis.go already implemented this
// logic correctly), but this method — the one actually wired to
// AppModule.InitGenesis — only ever restored Params. Any genesis.json that
// declared this module's authorities/threshold had that data discarded on
// every chain init.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	for _, auth := range genState.Authorities {
		if auth == nil {
			continue
		}
		if err := k.Authorities.Set(ctx, auth.Address, *auth); err != nil {
			return err
		}
	}

	if err := k.Threshold.Set(ctx, genState.Threshold); err != nil {
		return err
	}

	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	var authorities []*types.Authority
	err = k.Authorities.Walk(ctx, nil, func(_ string, auth types.Authority) (bool, error) {
		a := auth
		authorities = append(authorities, &a)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	genesis.Authorities = authorities

	threshold, err := k.Threshold.Get(ctx)
	if err != nil {
		threshold = 0
	}
	genesis.Threshold = threshold

	return genesis, nil
}
