package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),

		Authorities: []*Authority{
			{
				Address: "cosmos1sh954jkt35mpzjnsv6ywq2dgewhjp5de8uq5ft", // gov OR your control wallet
				Active:  true,
			},
						{
				Address: "cosmos17lnnehn6u663avxm2lcd6f7h3xhtuxrcrat33f", // gov OR your control wallet
				Active:  true,
			},
						{
				Address: "cosmos1ycq48wjzxxuecalj7y7ucv2sggt2yfswn3vfzq", // gov OR your control wallet
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
