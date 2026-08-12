// Package types (this file) defines the x/pramaan module's genesis default
// and validation logic.
package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation. Only validates Params —
// Authorities/Threshold are not checked here (no invariant is enforced on
// them, unlike x/authority's genesis which requires a ROOT entry).
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
