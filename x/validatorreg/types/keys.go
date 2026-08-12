// Package types (this file) defines the x/validatorreg module's store keys
// and the named constants for validator-proposal status values.
package types

import "cosmossdk.io/collections"

// Module identity constants used to register this module with the SDK's
// module manager and message/query routers.
const (
	// ModuleName defines the module name
	ModuleName = "validatorreg"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// GovModuleName duplicates the gov module's name to avoid a dependency with x/gov.
	// It should be synced with the gov module's name if it is ever changed.
	// See: https://github.com/cosmos/cosmos-sdk/blob/v0.52.0-beta.2/x/gov/types/keys.go#L9
	GovModuleName = "gov"
)

// ParamsKey is the prefix to retrieve all Params
var ParamsKey = collections.NewPrefix("p_validatorreg")

// ValidatorKeyPrefix is the collections-API prefix for the Validators map.
var ValidatorKeyPrefix = collections.NewPrefix("validator")

// ValidatorProposal.Status values. Always reference these constants instead
// of raw string literals so a typo can't silently break a status
// comparison. There is deliberately no ProposalStatusRejected yet — see
// SECURITY_CHANGELOG.md bug #10: a proposal that fails to reach quorum has
// no terminal rejection state, it just stays PENDING forever.
const (
	ProposalStatusPending   = "PENDING"
	ProposalStatusApproved  = "APPROVED"
	ProposalStatusActivated = "ACTIVATED"
)

// Stateless upper bounds on free-form string inputs accepted by the
// validatorreg message handlers, to keep a single transaction from
// persisting an arbitrarily large value into committed state (state-bloat /
// DoS guard). Generous relative to legitimate use, well below anything that
// threatens node memory/disk.
const (
	MaxDomainLen  = 128
	MaxDataLen    = 8192
	MaxAddressLen = 128
)
