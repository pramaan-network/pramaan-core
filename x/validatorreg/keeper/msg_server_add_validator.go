package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/validatorreg/types"
)

func (k msgServer) AddValidator(
	goCtx context.Context,
	msg *types.MsgAddValidator,
) (*types.MsgAddValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔥 NEW AUTHORITY CHECK (MULTI-AUTHORITY)
	if !k.pramaanKeeper.IsAuthority(ctx, msg.Creator) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "only authority can add validator")
	}

	// 🚫 ALREADY EXISTS
	if k.Keeper.IsValidator(ctx, msg.Address) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "validator already exists")
	}

	// ✅ ADD VALIDATOR
	if err := k.Keeper.AddValidator(ctx, msg.Address); err != nil {
		return nil, err
	}

	return &types.MsgAddValidatorResponse{}, nil
}