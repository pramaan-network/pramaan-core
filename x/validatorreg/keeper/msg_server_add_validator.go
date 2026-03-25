package keeper

import (
	"context"

	"pramaan/x/validatorreg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddValidator(
	goCtx context.Context,
	msg *types.MsgAddValidator,
) (*types.MsgAddValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔥 AUTHORITY CHECK
	auth, found := k.authorityKeeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator not found in authority")
	}

	if auth.Role != "ROOT" {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "only ROOT can add validator")
	}

	// 🚫 ALREADY EXISTS
	if k.Keeper.IsValidator(ctx, msg.Address) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "validator already exists")
	}

	// ✅ ADD VALIDATOR
	if err := k.Keeper.AddValidator(ctx, msg.Address, msg.Domain); err != nil {
		return nil, err
	}

	// 🔥 EMIT EVENT
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_added",
			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)

	return &types.MsgAddValidatorResponse{}, nil
}
