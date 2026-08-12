// Package keeper (this file) implements the (disabled) MsgUpdateParams
// handler for x/authority.
package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"pramaan/x/authority/types"
)

// UpdateParams is intentionally disabled for the authority module (not used in PRAMAAN).
// It must still return a non-nil response value alongside the error; returning (nil, nil)
// is unsafe because the SDK's msg service router does not expect a typed-nil response.
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return &types.MsgUpdateParamsResponse{}, errorsmod.Wrap(types.ErrInvalidSigner, "UpdateParams is disabled for the authority module")
}
