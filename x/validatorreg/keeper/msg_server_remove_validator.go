// Package keeper (this file) implements the x/validatorreg
// MsgRemoveValidator handler — the mirror of AddValidator (see
// msg_server_add_validator.go for the package-level overview).
package keeper

import (
	"context"

	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/validatorreg/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RemoveValidator handles MsgRemoveValidator: requires the signer to hold
// the ROOT role (same authorization model as AddValidator — this is the
// check that was previously a dead empty-byte comparison; see
// SECURITY_CHANGELOG.md bug #1), confirms the target validator exists, then
// removes it outright. Note this only removes the validatorreg-local
// record — it does not also revoke the corresponding entry in
// x/authority's Authorities map, so a removed validator could still show
// up there unless something else cleans that up.
func (k msgServer) RemoveValidator(
	goCtx context.Context,
	msg *types.MsgRemoveValidator,
) (*types.MsgRemoveValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Authority check (same model as AddValidator: ROOT role via the authority module).
	auth, found := k.authorityKeeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "creator not found in authority")
	}

	if auth.Role != authoritytypes.RoleRoot {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "only ROOT can remove validator")
	}

	// Confirm the target validator actually exists.
	if !k.Keeper.IsValidator(ctx, msg.Address) {
		return nil, errorsmod.Wrap(types.ErrInvalidSigner, "validator does not exist")
	}

	// Remove it.
	if err := k.Keeper.RemoveValidator(ctx, msg.Address); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_removed",
			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)

	return &types.MsgRemoveValidatorResponse{}, nil
}
