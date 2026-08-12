// Package types (this file) defines the x/validatorreg module's genesis
// default and validation logic.
package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon
// any failure. Only validates Params — Validators/Proposals aren't checked
// for duplicate addresses/IDs here (InitGenesis's plain Set calls would
// silently overwrite a duplicate rather than error — see keeper/genesis.go).
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
