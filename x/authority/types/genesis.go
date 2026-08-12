// Package types (this file) defines the x/authority module's genesis
// default and validation logic. Genesis for this module is unusual among
// the chain's modules in that it's load-bearing at boot: without at least
// one valid ROOT authority in genesis, the chain has no way to bootstrap
// any other role, so Validate() enforces that invariant strictly.
package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns default genesis state.
func DefaultGenesis() *GenesisState {
	params := DefaultParams()

	return &GenesisState{
		Params: &params,
	}
}

// Validate checks that at least one ROOT authority is present, every address is
// a valid Bech32 account address, no address appears twice, and every role is one
// of the recognised constants.
func (gs *GenesisState) Validate() error {
	if len(gs.Authorities) == 0 {
		return fmt.Errorf("genesis must contain at least one ROOT authority")
	}

	validRoles := map[string]bool{
		RoleRoot:      true,
		RoleAuthority: true,
		RoleValidator: true,
		RoleIssuer:    true,
	}

	seen := make(map[string]struct{})
	rootFound := false

	for _, auth := range gs.Authorities {
		if auth == nil {
			return fmt.Errorf("nil authority entry in genesis")
		}
		if auth.Address == "" {
			return fmt.Errorf("authority address cannot be empty")
		}
		if _, err := sdk.AccAddressFromBech32(auth.Address); err != nil {
			return fmt.Errorf("invalid authority address %q: %w", auth.Address, err)
		}
		if _, dup := seen[auth.Address]; dup {
			return fmt.Errorf("duplicate authority address: %s", auth.Address)
		}
		seen[auth.Address] = struct{}{}

		if !validRoles[auth.Role] {
			return fmt.Errorf("invalid authority role %q for address %s", auth.Role, auth.Address)
		}
		if auth.Role == RoleRoot {
			rootFound = true
		}
	}

	if !rootFound {
		return fmt.Errorf("genesis must contain at least one ROOT authority")
	}

	return nil
}
