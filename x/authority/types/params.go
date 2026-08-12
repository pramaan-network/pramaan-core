// Package types (this file) implements the x/authority module's Params
// type helpers. The module currently defines no tunable parameters — Params
// is an empty struct — but the constructors/Validate are kept so the module
// follows the same Params-collection pattern as every other module in this
// chain, and so adding a real parameter later is a small, local change.
package types

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params. Always succeeds today since Params
// has no fields; kept as a named hook so future fields get validated here
// instead of scattered across call sites.
func (p Params) Validate() error {

	return nil
}
