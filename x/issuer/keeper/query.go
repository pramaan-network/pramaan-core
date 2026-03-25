package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/issuer/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

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