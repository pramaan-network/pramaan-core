package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/validatorreg/types"
)

func (k msgServer) AddValidator(goCtx context.Context, msg *types.MsgAddValidator) (*types.MsgAddValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔒 ADMIN CHECK (CRITICAL)
	if msg.Creator != k.Keeper.Admin {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "only admin can add validator")
	}

	// 🚫 Already exists
	if k.Keeper.IsValidator(ctx, msg.Address) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "validator already exists")
	}

	// ✅ Add
	if err := k.Keeper.AddValidator(ctx, msg.Address); err != nil {
		return nil, err
	}

	return &types.MsgAddValidatorResponse{}, nil
}