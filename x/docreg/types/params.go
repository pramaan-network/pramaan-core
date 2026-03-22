package types

// DefaultParams returns default module parameters
func DefaultParams() Params {
	return Params{}
}

// Validate validates params
func (p Params) Validate() error {
	return nil
}