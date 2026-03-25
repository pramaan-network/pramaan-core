package keeper

import (
	"context"

	"pramaan/x/authority/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// --------------------
// Get single authority
// --------------------
func (q queryServer) Authority(ctx context.Context, req *types.QueryGetAuthorityRequest) (*types.QueryGetAuthorityResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	authority, found := q.k.GetAuthority(sdkCtx, req.Address)
	if !found {
		return &types.QueryGetAuthorityResponse{}, nil
	}

	return &types.QueryGetAuthorityResponse{
		Authority: &authority,
	}, nil
}

// --------------------
// Get all authorities
// --------------------
func (q queryServer) Authorities(ctx context.Context, req *types.QueryAllAuthoritiesRequest) (*types.QueryAllAuthoritiesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := q.k.storeService.OpenKVStore(sdkCtx)

	// ✅ Correct prefix iterator
	iterator, err := store.Iterator(types.AuthorityKeyPrefix, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var authorities []*types.Authority

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()

		// ✅ SAFETY: ensure key actually starts with prefix
		if len(key) == 0 || key[0] != types.AuthorityKeyPrefix[0] {
			continue
		}

		var authority types.Authority
		q.k.cdc.MustUnmarshal(iterator.Value(), &authority)

		// ✅ SKIP EMPTY ENTRIES (CRITICAL FIX)
		if authority.Address == "" {
			continue
		}

		authorities = append(authorities, &authority)
	}

	return &types.QueryAllAuthoritiesResponse{
		Authorities: authorities,
	}, nil
}
