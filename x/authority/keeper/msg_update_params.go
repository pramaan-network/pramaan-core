package keeper

import (
	"context"

	"pramaan/x/authority/types"
)

// Disable MsgUpdateParams completely (not used in PRAMAAN)
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	return nil, nil
}
