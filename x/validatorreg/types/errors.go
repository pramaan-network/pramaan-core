// Package types (this file) defines the x/validatorreg module's sentinel
// errors.
package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/validatorreg module sentinel errors, registered against ModuleName.
// Reused (somewhat loosely — see msg_server_add_validator.go/
// msg_server_remove_validator.go) as the generic "unauthorized/invalid
// request" error for this module's non-UpdateParams messages too, rather
// than each having its own dedicated error code.
var (
	// ErrInvalidSigner is returned when a message's signer does not match
	// the expected authority/role.
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
)
