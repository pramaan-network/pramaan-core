// Package types (this file) declares the "expected keeper" interfaces
// x/validatorreg depends on.
package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// PramaanKeeper is injected into this module's msgServer (see
// keeper/msg_server.go's NewMsgServerImpl) but is currently unused by every
// message handler — all of them check authorityKeeper instead. See
// SECURITY_CHANGELOG.md's "dead pramaan authority subsystem" entry: this is
// wiring for a second, parallel authority system (x/pramaan's own
// Authorities map) that duplicates x/authority's real one and isn't
// actually consulted anywhere in this module.
type PramaanKeeper interface {
	IsAuthority(ctx sdk.Context, addr string) bool
}
