package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),

		Authorities: []*Authority{
			{
				Address: "pramaan1kzu8xmf6fkml4jalptj6ed597gnf36qt9vnslv", // gov OR your control wallet
				Active:  true,
			},
						{
				Address: "pramaan1zt0hnv8y0vpt3jzy5w84fnr5fjjup5zqcllaqq", // gov OR your control wallet
				Active:  true,
			},
						{
				Address: "pramaan1073yvx37dazc74e7xrxwm6zardzcufdnr7m7qm", // gov OR your control wallet
				Active:  true,
			},
		},

		Threshold: 3,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
