// Package keeper (this file) implements InitGenesis/ExportGenesis for
// x/authority — the module whose genesis is most load-bearing in this
// chain, since it's the only way a ROOT authority can ever come into
// existence (there is no message to create one; see msg_server.go's
// AddAuthority, which explicitly refuses to create RoleRoot).
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/authority/types"
)

// InitGenesis loads the module's genesis state into the store at chain
// start. It re-validates the ROOT/duplicate invariants that
// types.GenesisState.Validate() also checks (defense in depth: this runs
// even if ValidateGenesis was somehow bypassed) and panics rather than
// returning an error, per Cosmos SDK convention for InitGenesis — an
// invalid genesis should halt startup, not silently produce a broken chain.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {

	// Set params
	if genState.Params != nil {
		k.Params.Set(ctx, *genState.Params)
	}

	if len(genState.Authorities) == 0 {
		panic("❌ genesis must contain at least ROOT")
	}

	rootFound := false
	seen := make(map[string]bool)

	for _, auth := range genState.Authorities {

		addr := auth.Address

		// prevent duplicate
		if seen[addr] {
			panic("❌ duplicate authority in genesis: " + addr)
		}
		seen[addr] = true

		// check role
		if auth.Role == types.RoleRoot {
			rootFound = true
		}

		k.SetAuthority(ctx, *auth)
	}

	if !rootFound {
		panic("❌ ROOT authority missing in genesis")
	}
}

// ExportGenesis reads the full authority set and params back out of the
// store, for `pramaand export` to write into a new genesis file.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	var authorities []*types.Authority

	k.IterateAuthorities(ctx, func(auth types.Authority) bool {
		a := auth
		authorities = append(authorities, &a)
		return false
	})

	params, _ := k.Params.Get(ctx)

	return &types.GenesisState{
		Params:      &params,
		Authorities: authorities,
	}
}