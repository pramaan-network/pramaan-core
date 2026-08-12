// Package types (this file) implements the x/docreg module's Params
// helpers. Like x/authority, Params is currently an empty struct — no
// tunable parameters are defined yet — but the module follows the standard
// Params-collection pattern so a real parameter (e.g. a max-metadata-size
// limit) can be added here later without restructuring anything else.
package types

// DefaultParams returns default module parameters
func DefaultParams() Params {
	return Params{}
}

// Validate validates params. Always succeeds today since Params has no
// fields.
func (p Params) Validate() error {
	return nil
}
