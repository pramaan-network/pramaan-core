// Package types (this file) defines the x/issuer module's sentinel errors.
package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/issuer module sentinel errors, registered against ModuleName so error
// codes are namespaced per-module.
var (
	// ErrInvalidSigner is returned when a message's signer does not match
	// the expected authority (used by UpdateParams).
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
)
