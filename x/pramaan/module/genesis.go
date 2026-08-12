// Package pramaan (this file) — DEAD CODE, not called by anything.
//
// AppModule.InitGenesis/ExportGenesis (module.go) call the Keeper *methods*
// of the same name in x/pramaan/keeper/genesis.go, not these free functions.
// This file predates that keeper-method version and was left behind when
// the logic moved; it happened to be the only place Authorities/Threshold
// genesis restoration was implemented correctly (see
// SECURITY_CHANGELOG.md #22), which is exactly why that bug existed — the
// working code was sitting right here, just never wired in.
//
// Kept as-is (not deleted) since removing code is a bigger, more
// judgment-laden call than documenting it — see SECURITY_CHANGELOG.md's
// "dead pramaan authority subsystem" entry for the fuller writeup and the
// case for deleting this along with the rest of that subsystem.
package pramaan

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/pramaan/keeper"
	"pramaan/x/pramaan/types"
)

// InitGenesis is UNUSED — see package comment above. Restores Authorities
// and Threshold from genesis state; kept only as a reference for what
// keeper.Keeper.InitGenesis now does properly.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// 🔥 STORE AUTHORITIES
	for _, auth := range genState.Authorities {
		err := k.Authorities.Set(ctx, auth.Address, *auth)
		if err != nil {
			panic(err)
		}
		fmt.Println("LOADED AUTHORITY:", auth.Address)
	}

	// 🔥 STORE THRESHOLD
	err := k.Threshold.Set(ctx, genState.Threshold)
	if err != nil {
		panic(err)
	}

	fmt.Println("THRESHOLD:", genState.Threshold)
}

// ExportGenesis is UNUSED — see package comment above. Also incomplete even
// as a reference: it returns types.DefaultGenesis() unconditionally rather
// than reading back live state, unlike keeper.Keeper.ExportGenesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	genesis := types.DefaultGenesis()
	return *genesis
}
