// Package types (this file) defines the x/pramaan module's sentinel
// errors.
package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/pramaan module sentinel errors, registered against ModuleName.
var (
	// ErrInvalidSigner is returned when a message's signer does not match
	// the expected authority (used by UpdateParams).
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
)
