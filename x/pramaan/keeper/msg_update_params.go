// Package keeper (this file) implements the x/pramaan MsgUpdateParams
// handler: a standard gov-gated params update.
package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"

	"pramaan/x/pramaan/types"
)

// UpdateParams handles MsgUpdateParams: requires the signer to match the
// module's configured authority (the gov module account by default), then
// validates and persists the new Params.
func (k msgServer) UpdateParams(ctx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	authority, err := k.addressCodec.StringToBytes(req.Authority)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if !bytes.Equal(k.GetAuthority(), authority) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", expectedAuthorityStr, req.Authority)
	}

	if err := req.Params.Validate(); err != nil {
		return nil, err
	}

	if err := k.Params.Set(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
