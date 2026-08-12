// Package keeper (this file) implements the gRPC QueryServer for x/issuer:
// single-issuer lookup and a full listing.
package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/issuer/types"
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

// ==============================
// GET ISSUER
// ==============================

// GetIssuer handles the QueryGetIssuerRequest gRPC query: look up a single
// issuer by address.
func (k queryServer) GetIssuer(
	goCtx context.Context,
	req *types.QueryGetIssuerRequest,
) (*types.QueryGetIssuerResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	issuer, err := k.k.Issuers.Get(ctx, req.Address)
	if err != nil {
		return nil, fmt.Errorf("issuer not found")
	}

	return &types.QueryGetIssuerResponse{
		Issuer: &issuer,
	}, nil
}

// ==============================
// LIST ISSUERS
// ==============================

// Issuers handles the QueryIssuersRequest gRPC query: returns every
// registered issuer. Not paginated — fine while the issuer set is small,
// but would need pagination if it grows large (loads the full result set
// into memory, same tradeoff as x/authority's Authorities query).
func (k queryServer) Issuers(
	goCtx context.Context,
	req *types.QueryIssuersRequest,
) (*types.QueryIssuersResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	var list []*types.Issuer

	err := k.k.Issuers.Walk(ctx, nil, func(key string, value types.Issuer) (bool, error) {
		v := value
		list = append(list, &v)
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryIssuersResponse{
		Issuers: list,
	}, nil
}