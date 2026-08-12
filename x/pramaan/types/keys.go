// Package types (this file) defines the x/pramaan module's store keys.
//
// x/pramaan is this repo's original/base module (predating the later split
// into authority/docreg/issuer/validatorreg — see keeper/keeper.go and
// SECURITY_CHANGELOG.md's "dead pramaan authority subsystem" note). It
// still holds a legacy Authorities/Threshold subsystem that duplicates
// x/authority's real one but is not wired into any live message handler.
package types

import "cosmossdk.io/collections"

// Module identity constants used to register this module with the SDK's
// module manager and message/query routers.
const (
	// ModuleName defines the module name
	ModuleName = "pramaan"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_pramaan")
