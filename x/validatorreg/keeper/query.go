// Package keeper (this file) implements the gRPC QueryServer for
// x/validatorreg: listing active validators and listing all
// validator-admission proposals.
package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/validatorreg/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// backed by the given Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

// queryServer adapts a Keeper to the generated types.QueryServer interface.
type queryServer struct {
	k Keeper
}

// Validators handles the QueryValidatorsRequest gRPC query: returns the
// addresses of every validator currently marked Active (inactive/removed
// validators are filtered out, not just omitted from a "removed" list —
// there is no way to distinguish "never existed" from "removed" via this
// query alone). Not paginated.
func (k queryServer) Validators(goCtx context.Context, req *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var list []string

	err := k.k.Validators.Walk(ctx, nil, func(key string, value types.Validator) (bool, error) {
		if value.Active {
			list = append(list, key)
		}
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorsResponse{
		Validators: list,
	}, nil
}

// Proposals handles the QueryProposalsRequest gRPC query: returns every
// validator-admission proposal regardless of status (PENDING/APPROVED/
// ACTIVATED — there's no REJECTED status yet, see SECURITY_CHANGELOG.md
// bug #10). Not paginated, and has no status filter — a client wanting
// only pending proposals must filter client-side.
func (k queryServer) Proposals(
	goCtx context.Context,
	req *types.QueryProposalsRequest,
) (*types.QueryProposalsResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	var list []*types.ValidatorProposal

	err := k.k.Proposals.Walk(ctx, nil, func(key uint64, value types.ValidatorProposal) (bool, error) {
		v := value
		list = append(list, &v)
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryProposalsResponse{
		Proposals: list,
	}, nil
}
