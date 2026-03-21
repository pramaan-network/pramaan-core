package keeper

import (
	"context"

	"pramaan/x/docreg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferDocument(goCtx context.Context, msg *types.MsgTransferDocument) (*types.MsgTransferDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Get document
	doc, found := k.Keeper.GetDocumentByID(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "document not found")
	}

	// 2. 🔥 ENFORCE TOKEN RULE
	if !doc.Transferable {
		return nil, errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "this document is non-transferable (SBT)")
	}

	// 3. Update owner
	doc.Owner = msg.NewOwner

	// 4. Save updated document
	if err := k.Keeper.Documents.Set(ctx, doc.Id, doc); err != nil {
		return nil, err
	}

	return &types.MsgTransferDocumentResponse{}, nil
}
