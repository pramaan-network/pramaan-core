// Package types (this file) defines the x/docreg module's store keys and
// the named constants for document status/token-type/doc-type values.
package types

import "cosmossdk.io/collections"

// Module identity constants used to register this module with the SDK's
// module manager and message/query routers.
const (
	// ModuleName defines the module name
	ModuleName = "docreg"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_docreg")

// Document.Status values. Always reference these constants instead of raw
// string literals so a typo can't silently break a status comparison.
const (
	DocumentStatusIssued = "ISSUED"
	// DocumentStatusRevoked is reserved for when MsgRevokeDocument is added.
	DocumentStatusRevoked = "REVOKED"
)

// Document.TokenType values.
const (
	TokenTypeSBT = "SBT" // non-transferable
	TokenTypeNFT = "NFT" // transferable
)

// DocTypeLandProperty is the doc_type that triggers NFT (transferable) tokenization
// instead of the default SBT (non-transferable) behavior.
const DocTypeLandProperty = "land_property"

// Stateless upper bounds on free-form string inputs accepted by the docreg
// message handlers. Storage is gas-metered, but there is no protocol-level
// cap otherwise, so a single transaction could persist an arbitrarily large
// value into committed state (a state-bloat / DoS vector). These caps keep
// per-record size bounded and predictable. Values are generous relative to
// legitimate use (hashes, IDs, and short metadata blobs) but far below the
// point where they threaten node memory/disk.
const (
	MaxIDLen       = 256
	MaxHashLen     = 256
	MaxDocTypeLen  = 128
	MaxAddressLen  = 128
	MaxMetadataLen = 8192
)
