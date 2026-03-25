package types

import "cosmossdk.io/collections"

const (
	ModuleName = "authority"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

var (
	AuthorityKeyPrefix = []byte{0x01}

	// ✅ ONLY THIS ONE (correct type)
	ParamsKey = collections.NewPrefix("params")
)

// GetAuthorityKey returns the store key for an authority address
func GetAuthorityKey(address string) []byte {
	return append(AuthorityKeyPrefix, []byte(address)...)
}
