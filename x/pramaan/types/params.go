// Package types (this file) implements the x/pramaan module's Params
// helpers. Params is currently an empty struct.
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
