// Package keeper (this file) implements the x/docreg MsgTransferDocument
// handler: moving ownership of a transferable (NFT-type) document from one
// address to another, keeping the OwnerIndex in sync.
package keeper

import (
	"context"
	"errors"

	"pramaan/x/docreg/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TransferDocument handles MsgTransferDocument. Validation order: reject an
// empty new_owner first (cheapest check), then confirm the document exists,
// then that the signer is its current owner, then that it's actually
// transferable (not an SBT), then that the transfer isn't a no-op
// self-transfer. Only after all of that does it touch the OwnerIndex —
// removing the document ID from the old owner's list and adding it to the
// new owner's list — before finally updating the document's Owner field.
func (k msgServer) TransferDocument(
	goCtx context.Context,
	msg *types.MsgTransferDocument,
) (*types.MsgTransferDocumentResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate new owner up front — cheapest check, do it before touching
	// the store.
	if msg.NewOwner == "" {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "new_owner is required")
	}

	// Bound free-form inputs (state-bloat / DoS guard — see types.Max*Len).
	if len(msg.Id) > types.MaxIDLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "id exceeds max length %d", types.MaxIDLen)
	}
	if len(msg.NewOwner) > types.MaxAddressLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "new_owner exceeds max length %d", types.MaxAddressLen)
	}

	// Get existing document.
	doc, found := k.Keeper.GetDocumentByID(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrDocumentNotFound, "document %s not found", msg.Id)
	}

	// Check ownership: only the current owner may initiate a transfer.
	if doc.Owner != msg.Creator {
		return nil, errorsmod.Wrapf(types.ErrNotDocumentOwner, "not document owner")
	}

	// Check transferability: SBT-type documents (see types.TokenTypeSBT)
	// can never be transferred, by design.
	if !doc.Transferable {
		return nil, errorsmod.Wrapf(types.ErrDocumentNotTransferable, "document is non-transferable (SBT)")
	}

	// Reject no-op self-transfers.
	if msg.NewOwner == doc.Owner {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "new_owner must differ from current owner")
	}

	oldOwner := doc.Owner
	newOwner := msg.NewOwner

	// ==============================
	// Remove the document ID from the old owner's index entry.
	// ==============================

	oldList, err := k.Keeper.OwnerIndex.Get(ctx, oldOwner)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		var updated []string
		for _, id := range oldList.Items {
			if id != doc.Id {
				updated = append(updated, id)
			}
		}
		oldList.Items = updated
		if err := k.Keeper.OwnerIndex.Set(ctx, oldOwner, oldList); err != nil {
			return nil, err
		}
	}

	// ==============================
	// Add the document ID to the new owner's index entry (creating it if
	// this is their first document).
	// ==============================

	newList, err := k.Keeper.OwnerIndex.Get(ctx, newOwner)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		newList = types.StringList{Items: []string{}}
	}

	newList.Items = append(newList.Items, doc.Id)

	if err := k.Keeper.OwnerIndex.Set(ctx, newOwner, newList); err != nil {
		return nil, err
	}

	// ==============================
	// Persist the document with its Owner field updated.
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
