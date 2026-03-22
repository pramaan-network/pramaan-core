package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/validatorreg/types"
)

func (k msgServer) RemoveValidator(
	goCtx context.Context,
	msg *types.MsgRemoveValidator,
) (*types.MsgRemoveValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔐 AUTHORITY CHECK
	authorityAddr, err := k.Keeper.addressCodec.BytesToString(k.Keeper.authority)
	if err != nil {
		return nil, err
	}

	if msg.Creator != authorityAddr {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "only authority can remove validator")
	}

	// 🚫 NOT EXISTS
	if !k.Keeper.IsValidator(ctx, msg.Address) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "validator does not exist")
	}

	// ✅ REMOVE
	if err := k.Keeper.RemoveValidator(ctx, msg.Address); err != nil {
		return nil, err
	}

	return &types.MsgRemoveValidatorResponse{}, nil
}