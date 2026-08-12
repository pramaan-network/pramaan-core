// Package keeper (this file) implements a Keeper-level GetDocument lookup.
//
// NOTE: this duplicates queryServer.GetDocument in query.go (same request/
// response types, same underlying GetDocumentByID call). The queryServer
// version is the one actually registered as the gRPC handler
// (types.RegisterQueryServer in module/module.go); this Keeper method
// appears unused by the wired query path. Left as-is rather than removed —
// deleting a method is a bigger call than documenting it, and it's possible
// something outside this package still calls it directly.
package keeper

import (
	"context"

	"pramaan/x/docreg/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDocument looks up a document by ID directly on the Keeper (bypassing
// the queryServer wrapper). See the package-level note above.
func (k Keeper) GetDocument(goCtx context.Context, req *types.QueryGetDocumentRequest) (*types.QueryGetDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	doc, found := k.GetDocumentByID(ctx, req.Id)
	if !found {
		return &types.QueryGetDocumentResponse{}, nil
	}

	return &types.QueryGetDocumentResponse{
		Document: &doc,
	}, nil
}
