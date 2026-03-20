package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/docreg/types"
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

	// ✅ CREATE DOCUMENT
	doc := types.Document{
		Id:        msg.Id,
		Hash:      msg.Hash,
		Owner:     msg.Owner,
		Issuer:    msg.Issuer,
		Type:      msg.DocType,
		Status:    "ISSUED",
		Timestamp: ctx.BlockTime().Unix(),
		Metadata:  msg.Metadata,
	}

	// ✅ STORE
	if err := k.Keeper.SetDocument(ctx, doc); err != nil {
		return nil, err
	}

	return &types.MsgRegisterDocumentResponse{}, nil
}