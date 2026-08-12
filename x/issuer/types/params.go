// Package types (this file) implements the x/issuer module's Params
// helpers. Params is currently an empty struct — no tunable parameters are
// defined yet.
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
