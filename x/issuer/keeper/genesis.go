// Package keeper (this file) implements InitGenesis/ExportGenesis for
// x/issuer. See SECURITY_CHANGELOG.md #11: the Issuers field didn't exist
// in this module's genesis proto until that fix, so a genesis export ->
// reinit round-trip used to silently drop every registered issuer.
package keeper

import (
	"context"
	"fmt"

	"pramaan/x/issuer/types"
)

// InitGenesis initializes the module's state from a provided genesis state:
// Params, then every issuer record (restored via a plain Set — see
// types/genesis.go's Validate for the note on why a duplicate address in
// genesis wouldn't be caught here).
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	// Restore issuers so a `pramaand export` -> reinit round-trip doesn't silently drop them.
	for _, issuer := range genState.Issuers {
		if issuer == nil {
			continue
		}
		if err := k.Issuers.Set(ctx, issuer.Address, *issuer); err != nil {
			return fmt.Errorf("failed to restore issuer %s: %w", issuer.Address, err)
		}
	}

	return nil
}

// ExportGenesis returns the module's exported genesis: current Params plus
// every issuer record, walked from the store.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	err = k.Issuers.Walk(ctx, nil, func(_ string, value types.Issuer) (bool, error) {
		v := value
		genesis.Issuers = append(genesis.Issuers, &v)
		return false, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to export issuers: %w", err)
	}

	return genesis, nil
}
