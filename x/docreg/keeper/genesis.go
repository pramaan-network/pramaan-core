// Package keeper (this file) implements InitGenesis/ExportGenesis for
// x/docreg. See SECURITY_CHANGELOG.md #11 for why this file matters: the
// Documents field didn't exist in this module's genesis proto until that
// fix, so a genesis export -> reinit round-trip used to silently drop every
// registered document.
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

	// Restore documents (and rebuild the hash/owner/issuer indexes) via SetDocument
	// so a `pramaand export` -> reinit round-trip doesn't silently drop them.
	for _, doc := range genState.Documents {
		if doc == nil {
			continue
		}
		if err := k.SetDocument(ctx, *doc); err != nil {
			return fmt.Errorf("failed to restore document %s: %w", doc.Id, err)
		}
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

	var documents []*types.Document
	iter, err := k.Documents.Iterate(sdkCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to iterate documents: %w", err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		doc, err := iter.Value()
		if err != nil {
			return nil, fmt.Errorf("failed to read document during export: %w", err)
		}
		d := doc
		documents = append(documents, &d)
	}

	genesis := &types.GenesisState{
		Params:    params,
		Documents: documents,
	}

	return genesis, nil
}
