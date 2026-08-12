// Package types (this file) declares the "expected keeper" interfaces this
// module depends on. Following the standard Cosmos SDK pattern, x/authority
// depends on narrow interfaces (not concrete keeper types) from the modules
// it needs, so it can be compiled/tested without importing x/auth or x/bank
// directly, and so mocking those dependencies in tests is straightforward.
package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface for the Bank module.
//
// NOTE: this is injected into the keeper (see keeper/keeper.go) but no
// x/authority keeper method currently calls it — see SECURITY_CHANGELOG.md
// for the related finding that this module's bank permissions
// (Minter/Burner/Staking in app/app_config.go) are likewise unused.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
}

// ParamSubspace defines the expected legacy x/params Subspace interface.
// Unused by this module (it stores Params via the collections API, not a
// legacy subspace) — kept for parity with the standard ignite scaffold.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

// AuthorityKeeper is the interface other modules (docreg, issuer,
// validatorreg) depend on to read/write this module's authority registry
// without importing x/authority/keeper directly. This is the module's real
// public surface for cross-module role checks.
type AuthorityKeeper interface {
	GetAuthority(ctx sdk.Context, address string) (Authority, bool)
	SetAuthority(ctx sdk.Context, authority Authority)
}
