// Package keeper (this file) implements the x/docreg MsgUpdateParams
// handler — a standard gov-gated params update, unlike x/authority's
// version of the same message, which is unconditionally disabled.
package keeper

import (
	"bytes"
	"context"

	errorsmod "cosmossdk.io/errors"

	"pramaan/x/docreg/types"
)

// UpdateParams handles MsgUpdateParams: it requires the signer to match the
// module's configured authority (the gov module account by default — see
// module/depinject.go), validates the new Params, and persists them.
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

	if err := k.ParamsStore.Set(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
