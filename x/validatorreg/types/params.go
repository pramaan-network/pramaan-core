// Package types (this file) implements the x/validatorreg module's Params
// helpers. Params is currently an empty struct — notably, the
// ValidatorApprovalThreshold constant (keeper/msg_server.go) is NOT a
// module param, so it can't be changed via governance; see that file's
// comment for why.
package types

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params.
func (p Params) Validate() error {

	return nil
}
