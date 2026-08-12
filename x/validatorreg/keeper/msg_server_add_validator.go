// Package keeper (this file) implements the x/validatorreg MsgAddValidator
// handler: a direct, ROOT-only path to register a validator, bypassing the
// AUTHORITY-quorum proposal flow (ApplyValidator/ApproveValidator/
// ActivateValidator in msg_server.go) entirely. Useful for genesis-adjacent
// bootstrap or emergency admission; RemoveValidator is its mirror image.
package keeper

import (
	"context"

	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/validatorreg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddValidator handles MsgAddValidator: requires the signer to hold the
// ROOT role in x/authority, rejects a duplicate address, then registers
// the validator directly (Active by construction — see keeper.AddValidator).
func (k msgServer) AddValidator(
	goCtx context.Context,
	msg *types.MsgAddValidator,
) (*types.MsgAddValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Bound free-form inputs (state-bloat / DoS guard — see types.Max*Len).
	if len(msg.Address) > types.MaxAddressLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "address exceeds max length %d", types.MaxAddressLen)
	}
	if len(msg.Domain) > types.MaxDomainLen {
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "domain exceeds max length %d", types.MaxDomainLen)
	}

	// Authority check: creator must be a registered ROOT.
	auth, found := k.authorityKeeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator not found in authority")
	}

	if auth.Role != authoritytypes.RoleRoot {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "only ROOT can add validator")
	}

	// Reject a duplicate address.
	if k.Keeper.IsValidator(ctx, msg.Address) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "validator already exists")
	}

	// Register the validator.
	if err := k.Keeper.AddValidator(ctx, msg.Address, msg.Domain); err != nil {
		return nil, err
	}

	// Emit a standard event for indexers/explorers.
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_added",
			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)

	return &types.MsgAddValidatorResponse{}, nil
}
