// Package types (this file) defines the x/issuer module's genesis default
// and validation logic.
package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon
// any failure. Only validates Params — it does not check Issuers for
// duplicate addresses or dangling validator references (see
// keeper/genesis.go's InitGenesis, which restores issuers via a plain
// Set — a duplicate address in genesis would just overwrite silently rather
// than error, unlike docreg's SetDocument which rejects duplicates).
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
