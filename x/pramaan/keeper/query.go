// Package keeper (this file) declares the x/pramaan queryServer type and
// its constructor. The Params query itself is implemented in
// query_params.go. Note this module's other keeper methods (AddAuthority,
// RemoveAuthority, IsAuthority, SetThreshold, GetThreshold — see keeper.go)
// have no corresponding query RPCs exposed here, unlike x/authority's
// equivalent Authority/Authorities queries.
package keeper

import (
	"pramaan/x/pramaan/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

// queryServer adapts a Keeper to the generated types.QueryServer interface.
type queryServer struct {
	k Keeper
}
