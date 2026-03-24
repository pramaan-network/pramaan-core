package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/validatorreg/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

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