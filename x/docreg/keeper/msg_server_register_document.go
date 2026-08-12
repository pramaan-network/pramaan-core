// Package keeper (this file) implements the x/docreg MsgRegisterDocument
// handler: creating a new tokenized document record after validating the
// issuer's identity, active status, and domain match.
package keeper

import (
	"fmt"
	"context"

	"pramaan/x/docreg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RegisterDocument handles MsgRegisterDocument. Order of checks: required
// fields present, content hash not already registered, then issuer
// validation (issuer exists in x/issuer, is active, matches the message's
// signer, and its registered domain matches the document's doc_type).
// DocType additionally determines tokenization: everything defaults to a
// non-transferable SBT except DocTypeLandProperty, which mints a
// transferable NFT instead (see types.TokenTypeSBT/TokenTypeNFT).
func (k msgServer) RegisterDocument(goCtx context.Context, msg *types.MsgRegisterDocument) (*types.MsgRegisterDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Basic validation: all four identifying fields are required.
	if msg.Id == "" || msg.Hash == "" || msg.Owner == "" || msg.Issuer == "" {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "id, hash, owner, and issuer are required")
	}

	// Bound free-form inputs before they can be persisted (state-bloat / DoS
	// guard — see types.Max*Len). Cheap, stateless checks done up front.
	if len(msg.Id) > types.MaxIDLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "id exceeds max length %d", types.MaxIDLen)
	}
	if len(msg.Hash) > types.MaxHashLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "hash exceeds max length %d", types.MaxHashLen)
	}
	if len(msg.DocType) > types.MaxDocTypeLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "doc_type exceeds max length %d", types.MaxDocTypeLen)
	}
	if len(msg.Owner) > types.MaxAddressLen || len(msg.Issuer) > types.MaxAddressLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "owner/issuer exceeds max length %d", types.MaxAddressLen)
	}
	if len(msg.Metadata) > types.MaxMetadataLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidRequest, "metadata exceeds max length %d", types.MaxMetadataLen)
	}

	// Duplicate check: reject re-registering the same content hash.
	if k.Keeper.IsHashExists(ctx, msg.Hash) {
		return nil, errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "hash already registered")
	}

	// =====================================================
	// Strict issuer validation: the message's declared issuer must be a
	// real, active issuer in x/issuer, and must itself be the signer.
	// =====================================================

	issuer, err := k.issuerKeeper.GetIssuer(ctx, msg.Issuer)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "issuer not registered in issuer module")
	}

	// must be active
	if !issuer.Active {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "issuer is not active")
	}

	// must match signer — an issuer can only register documents as itself,
	// not on behalf of another issuer.
	if msg.Issuer != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "creator must be issuer")
	}

	// Domain enforcement: an issuer registered for one domain (e.g.
	// "education") cannot issue documents of another doc_type (e.g.
	// "land_property") — keeps issuer scope enforced at the protocol level
	// rather than trusting client-side conventions.
	if issuer.Domain != msg.DocType {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "issuer domain does not match document type")
	}

	// =====================================================
	// Tokenization: every doc_type defaults to a non-transferable SBT
	// except land_property, which is minted as a transferable NFT instead.
	// Marked TEMP by the original author — likely needs to become a
	// per-domain/per-issuer configurable rule rather than a single hardcoded
	// doc_type check if more transferable-asset domains are added later.
	// =====================================================

	tokenType := types.TokenTypeSBT
	transferable := false

	if msg.DocType == types.DocTypeLandProperty {
		tokenType = types.TokenTypeNFT
		transferable = true
	}

	// =====================================================
	// Assemble and persist the new document record.
	// =====================================================

	doc := types.Document{
		Id:           msg.Id,
		Hash:         msg.Hash,
		Owner:        msg.Owner,
		Issuer:       msg.Issuer,
		Type:         msg.DocType,
		Status:       types.DocumentStatusIssued,
		Timestamp:    ctx.BlockTime().Unix(),
		Metadata:     msg.Metadata,
		TokenType:    tokenType,
		Transferable: transferable,
	}

	// STORE
	if err := k.Keeper.SetDocument(ctx, doc); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
	sdk.NewEvent(
		"document.registered", // ✅ global standard

		sdk.NewAttribute("module", "docreg"),
		sdk.NewAttribute("action", "register"),

		sdk.NewAttribute("doc_id", doc.Id),
		sdk.NewAttribute("hash", doc.Hash),

		sdk.NewAttribute("issuer", doc.Issuer),
		sdk.NewAttribute("owner", doc.Owner),

		sdk.NewAttribute("doc_type", doc.Type),
		sdk.NewAttribute("domain", issuer.Domain),

		sdk.NewAttribute("token_type", doc.TokenType),
		sdk.NewAttribute("transferable", fmt.Sprintf("%t", doc.Transferable)),

		// 🔥 ENGINE SUPPORT
		sdk.NewAttribute("metadata", msg.Metadata),

		// 🔥 AUDIT
		sdk.NewAttribute("block_time", ctx.BlockTime().String()),
		sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
	),
	)

	return &types.MsgRegisterDocumentResponse{}, nil
}