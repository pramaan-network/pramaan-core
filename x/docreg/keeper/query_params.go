package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"pramaan/x/docreg/types"
)

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := k.ParamsStore.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{
		Params: &params,
	}, nil
}