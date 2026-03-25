package keeper

import (
	"context"

	"pramaan/x/docreg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RegisterDocument(goCtx context.Context, msg *types.MsgRegisterDocument) (*types.MsgRegisterDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// ✅ BASIC VALIDATION
	if msg.Id == "" || msg.Hash == "" || msg.Owner == "" || msg.Issuer == "" {
		return nil, errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "missing required fields")
	}

	// 🔥 DUPLICATE CHECK
	if k.Keeper.IsHashExists(ctx, msg.Hash) {
		return nil, errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "hash already registered")
	}

	// =====================================================
	// 🔐 NEW: STRICT ISSUER VALIDATION (CORE FIX)
	// =====================================================

	issuer, err := k.issuerKeeper.GetIssuer(ctx, msg.Issuer)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "issuer not registered in issuer module")
	}

	// must be active
	if !issuer.Active {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "issuer is not active")
	}

	// must match signer
	if msg.Issuer != msg.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "creator must be issuer")
	}

	// 🔥 DOMAIN ENFORCEMENT (YOUR RULE)
	if issuer.Domain != msg.DocType {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "issuer domain does not match document type")
	}

	// =====================================================
	// TOKENIZATION LOGIC (TEMP)
	// =====================================================

	tokenType := "SBT"
	transferable := false

	if msg.DocType == "land_property" {
		tokenType = "NFT"
		transferable = true
	}

	// =====================================================
	// CREATE DOCUMENT
	// =====================================================

	doc := types.Document{
		Id:           msg.Id,
		Hash:         msg.Hash,
		Owner:        msg.Owner,
		Issuer:       msg.Issuer,
		Type:         msg.DocType,
		Status:       "ISSUED",
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
			"document_registered",
			sdk.NewAttribute("id", doc.Id),
			sdk.NewAttribute("hash", doc.Hash),
			sdk.NewAttribute("issuer", doc.Issuer),
			sdk.NewAttribute("owner", doc.Owner),
			sdk.NewAttribute("type", doc.Type),
			sdk.NewAttribute("creator", msg.Issuer),
		),
	)

	return &types.MsgRegisterDocumentResponse{}, nil
}
