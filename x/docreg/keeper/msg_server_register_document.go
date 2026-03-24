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

	// 🔥 DUPLICATE CHECK (EARLY REJECTION)
	if k.Keeper.IsHashExists(ctx, msg.Hash) {
		return nil, errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "hash already registered")
	}

	auth, found := k.authorityKeeper.GetAuthority(ctx, msg.Issuer)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "issuer not registered")
	}

	if auth.Role != "ISSUER" {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "not a valid issuer")
	}

	// CREATE DOCUMENT
	// TOKENIZATION LOGIC
	tokenType := "SBT"
	transferable := false

	if msg.DocType == "land_property" {
		tokenType = "NFT"
		transferable = true
	}

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
