package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/docreg/types"
)

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