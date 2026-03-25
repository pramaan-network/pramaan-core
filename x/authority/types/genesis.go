package types

// DefaultGenesis returns default genesis state.
func DefaultGenesis() *GenesisState {
	params := DefaultParams()

	return &GenesisState{
		Params: &params,
	}
}
