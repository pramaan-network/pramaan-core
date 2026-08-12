// Package types (this file) defines the x/authority module's store keys,
// key-prefixes, and the shared role-name constants (RoleRoot/RoleAuthority/
// RoleValidator/RoleIssuer) that this module and its siblings (docreg,
// issuer, validatorreg) all reference for role checks.
package types

import "cosmossdk.io/collections"
import "bytes"

// Module identity constants used to register this module with the SDK's
// module manager, message router, and query router.
const (
	ModuleName = "authority"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	GovModuleName = "gov"
)

var (
	// AuthorityKeyPrefix prefixes every raw-KV authority record key (see
	// GetAuthorityKey/IsAuthorityKey below).
	AuthorityKeyPrefix = []byte{0x01}

	// ParamsKey is the collections-API prefix for this module's Params item.
	ParamsKey = collections.NewPrefix("params")
)

// Authority roles. These are shared across modules (authority, validatorreg,
// issuer) — always reference these constants instead of raw string literals so
// a typo can't silently break a role check.
const (
	RoleRoot      = "ROOT"
	RoleAuthority = "AUTHORITY"
	RoleValidator = "VALIDATOR"
	RoleIssuer    = "ISSUER"
)

// Stateless upper bounds on free-form string inputs accepted by the
// authority message handlers (state-bloat / DoS guard). Generous for
// legitimate use, well below anything that threatens node memory/disk.
const (
	MaxAddressLen  = 128
	MaxPubKeyLen   = 1024
	MaxMetadataLen = 8192
)

// GetAuthorityKey returns the raw store key for an authority address:
// AuthorityKeyPrefix followed by the address bytes.
func GetAuthorityKey(address string) []byte {
	return append(AuthorityKeyPrefix, []byte(address)...)
}

// IsAuthorityKey reports whether a raw store key belongs to the authority
// key-space (i.e. starts with AuthorityKeyPrefix), used when iterating the
// whole module store to skip over non-authority entries (e.g. Params).
func IsAuthorityKey(key []byte) bool {
	return len(key) >= len(AuthorityKeyPrefix) &&
		bytes.Equal(key[:len(AuthorityKeyPrefix)], AuthorityKeyPrefix)
}
