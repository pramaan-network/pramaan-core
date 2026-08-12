// Package keeper (this file) implements the gRPC QueryServer for looking up
// individual authorities and listing all of them.
package keeper

import (
	"context"

	"pramaan/x/authority/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

// --------------------
// Get single authority
// --------------------

// Authority handles the QueryGetAuthorityRequest gRPC query, looking up a
// single authority by address. A not-found address returns an empty
// response rather than an error, since "this address is not an authority"
// is a normal, expected answer for a query, not a fault condition.
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

// Authorities handles the QueryAllAuthoritiesRequest gRPC query, returning
// every registered authority in the store. Not paginated — fine while the
// authority set is small (expected to be at most a few dozen ROOT/AUTHORITY/
// VALIDATOR/ISSUER accounts), but would need pagination support if this set
// ever grows large, since it currently loads the entire result into memory.
func (q queryServer) Authorities(ctx context.Context, req *types.QueryAllAuthoritiesRequest) (*types.QueryAllAuthoritiesResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	store := q.k.storeService.OpenKVStore(sdkCtx)

	// Iterate every key under the authority prefix.
	iterator, err := store.Iterator(types.AuthorityKeyPrefix, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var authorities []*types.Authority

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()

		// Defensive check: confirm the key actually starts with the
		// authority prefix before decoding its value as an Authority. Should
		// always be true given the iterator bounds above, but guards against
		// a future key-space change silently corrupting query results.
		if len(key) == 0 || key[0] != types.AuthorityKeyPrefix[0] {
			continue
		}

		var authority types.Authority
		q.k.cdc.MustUnmarshal(iterator.Value(), &authority)

		// Skip zero-value entries (empty address) rather than returning
		// them to the client — a decoded-but-empty Authority indicates a
		// stale/corrupt record, not a real registered authority.
		if authority.Address == "" {
			continue
		}

		authorities = append(authorities, &authority)
	}

	return &types.QueryAllAuthoritiesResponse{
		Authorities: authorities,
	}, nil
}
