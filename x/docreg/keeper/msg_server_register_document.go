package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/docreg/types"
)

func (k msgServer) RegisterDocument(goCtx context.Context, msg *types.MsgRegisterDocument) (*types.MsgRegisterDocumentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🚫 Basic validation
	if msg.Id == "" || msg.Hash == "" || msg.Owner == "" || msg.Issuer == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	// 🚫 Duplicate hash check
	if k.Keeper.IsHashExists(ctx, msg.Hash) {
		return nil, fmt.Errorf("document already exists")
	}

	// ✅ Create document
	doc := types.Document{
		Id:        msg.Id,
		Hash:      msg.Hash,
		Owner:     msg.Owner,
		Issuer:    msg.Issuer,
		Type:      msg.DocType, // 🔥 IMPORTANT FIX
		Status:    "ISSUED",
		Timestamp: ctx.BlockTime().Unix(),
		Metadata:  msg.Metadata,
	}

	// ✅ Store
	if err := k.Keeper.SetDocument(ctx, doc); err != nil {
		return nil, err
	}

	return &types.MsgRegisterDocumentResponse{}, nil
}
