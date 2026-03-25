package keeper

import (
	"context"

	"pramaan/x/docreg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferDocument(
	goCtx context.Context,
	msg *types.MsgTransferDocument,
) (*types.MsgTransferDocumentResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔹 Get existing document
	doc, found := k.Keeper.GetDocumentByID(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "document not found")
	}

	// 🔹 Check ownership
	if doc.Owner != msg.Creator {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "not document owner")
	}

	// 🔹 Check transferability
	if !doc.Transferable {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "document is non-transferable (SBT)")
	}

	oldOwner := doc.Owner
	newOwner := msg.NewOwner

	// ==============================
	// 🔥 REMOVE FROM OLD OWNER INDEX
	// ==============================

	oldList, err := k.Keeper.OwnerIndex.Get(ctx, oldOwner)
	if err == nil {
		var updated []string
		for _, id := range oldList.Items {
			if id != doc.Id {
				updated = append(updated, id)
			}
		}
		oldList.Items = updated
		_ = k.Keeper.OwnerIndex.Set(ctx, oldOwner, oldList)
	}

	// ==============================
	// 🔥 ADD TO NEW OWNER INDEX
	// ==============================

	newList, err := k.Keeper.OwnerIndex.Get(ctx, newOwner)
	if err != nil {
		newList = types.StringList{Items: []string{}}
	}

	newList.Items = append(newList.Items, doc.Id)

	if err := k.Keeper.OwnerIndex.Set(ctx, newOwner, newList); err != nil {
		return nil, err
	}

	// ==============================
	// 🔥 UPDATE DOCUMENT
	// ==============================

	doc.Owner = newOwner

	if err := k.Keeper.Documents.Set(ctx, doc.Id, doc); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"document_transferred",
			sdk.NewAttribute("id", doc.Id),
			sdk.NewAttribute("from", oldOwner),
			sdk.NewAttribute("to", newOwner),
			sdk.NewAttribute("type", doc.Type),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)

	return &types.MsgTransferDocumentResponse{}, nil
}
