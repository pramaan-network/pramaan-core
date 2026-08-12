// Package types (this file) defines the x/issuer module's store keys.
package types

import "cosmossdk.io/collections"

// Module identity constants used to register this module with the SDK's
// module manager and message/query routers.
const (
	// ModuleName defines the module name
	ModuleName = "issuer"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_issuer")

// Stateless upper bounds on free-form string inputs accepted by the issuer
// message handlers (state-bloat / DoS guard). Generous for legitimate use,
// well below anything that threatens node memory/disk.
const (
	MaxDomainLen  = 128
	MaxAddressLen = 128
)
