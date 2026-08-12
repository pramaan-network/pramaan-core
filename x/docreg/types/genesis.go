// Package types (this file) defines the x/docreg module's genesis default
// and validation logic.
package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon
// any failure. Only validates Params today — it does not check Documents
// for duplicate IDs/hashes or dangling owner/issuer references, since those
// invariants are re-established by SetDocument during InitGenesis (see
// keeper/genesis.go), which will itself error out on a duplicate ID/hash.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
