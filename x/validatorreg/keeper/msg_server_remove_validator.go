package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/validatorreg/types"
)

func (k msgServer) RemoveValidator(goCtx context.Context, msg *types.MsgRemoveValidator) (*types.MsgRemoveValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔒 ADMIN CHECK
	if msg.Creator != k.Keeper.Admin {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "only admin can remove validator")
	}

	// 🚫 Not exists
	if !k.Keeper.IsValidator(ctx, msg.Address) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "validator not found")
	}

	if err := k.Keeper.RemoveValidator(ctx, msg.Address); err != nil {
		return nil, err
	}

	return &types.MsgRemoveValidatorResponse{}, nil
}