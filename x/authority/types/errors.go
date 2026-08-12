// Package types holds the x/authority module's protobuf-independent Go
// types: sentinel errors, params helpers, genesis validation, codec
// registration, and expected-keeper interfaces. This file defines the
// module's sentinel errors.
package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/authority module sentinel errors. Registered against ModuleName so error
// codes are namespaced per-module (Cosmos SDK convention) and don't collide
// with error codes registered by other modules in this chain.
var (
	// ErrInvalidSigner is returned when a message's signer does not match
	// the expected authority. Used both for the disabled UpdateParams path
	// (see keeper/msg_update_params.go) and for AddAuthority role checks.
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
)
