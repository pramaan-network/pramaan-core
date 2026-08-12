// Package keeper (this file) implements InitGenesis/ExportGenesis for
// x/validatorreg. See SECURITY_CHANGELOG.md #11: the Validators/Proposals
// fields didn't exist in this module's genesis proto until that fix, so a
// genesis export -> reinit round-trip used to silently drop every
// registered validator and pending/approved proposal.
package keeper

import (
	"context"
	"fmt"

	"pramaan/x/validatorreg/types"
)

// InitGenesis initializes the module's state from a provided genesis state:
// Params, then every validator record, then every proposal — and, since
// ProposalCount is what ApplyValidator uses to allocate the next proposal
// ID, it also fast-forwards that counter to the highest restored proposal
// ID so a freshly re-initialized chain can't hand out an ID that collides
// with one already present in the restored genesis.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	// Restore validators so a `pramaand export` -> reinit round-trip doesn't silently drop them.
	for _, val := range genState.Validators {
		if val == nil {
			continue
		}
		if err := k.Validators.Set(ctx, val.Address, *val); err != nil {
			return fmt.Errorf("failed to restore validator %s: %w", val.Address, err)
		}
	}

	// Restore validator-admission proposals, and keep the proposal counter consistent
	// with the highest restored proposal ID so newly submitted proposals don't collide.
	var maxID uint64
	for _, p := range genState.Proposals {
		if p == nil {
			continue
		}
		if err := k.Proposals.Set(ctx, p.Id, *p); err != nil {
			return fmt.Errorf("failed to restore proposal %d: %w", p.Id, err)
		}
		if p.Id > maxID {
			maxID = p.Id
		}
	}
	if maxID > 0 {
		if err := k.ProposalCount.Set(ctx, maxID); err != nil {
			return fmt.Errorf("failed to restore proposal count: %w", err)
		}
	}

	return nil
}

// ExportGenesis returns the module's exported genesis: current Params plus
// every validator and proposal record, walked from the store.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	err = k.Validators.Walk(ctx, nil, func(_ string, value types.Validator) (bool, error) {
		v := value
		genesis.Validators = append(genesis.Validators, &v)
		return false, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to export validators: %w", err)
	}

	err = k.Proposals.Walk(ctx, nil, func(_ uint64, value types.ValidatorProposal) (bool, error) {
		v := value
		genesis.Proposals = append(genesis.Proposals, &v)
		return false, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to export proposals: %w", err)
	}

	return genesis, nil
}
