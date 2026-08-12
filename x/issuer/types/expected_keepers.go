// Package types (this file) declares the "expected keeper" interfaces
// x/issuer depends on. Unlike most sibling modules, issuer depends on two
// other custom modules (authority, validatorreg) in addition to the
// standard auth/bank — CreateIssuer needs to confirm its caller is both a
// registered AUTHORITY-role VALIDATOR (via AuthorityKeeper) and an active
// validator for the requested domain (via ValidatorKeeper).
package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "pramaan/x/authority/types"
	validatortypes "pramaan/x/validatorreg/types"
)

// ==============================
// AUTHORITY KEEPER
// ==============================

// AuthorityKeeper is the interface used to confirm a CreateIssuer caller
// holds the VALIDATOR role in x/authority.
type AuthorityKeeper interface {
	GetAuthority(ctx sdk.Context, address string) (authoritytypes.Authority, bool)
}

// ==============================
// VALIDATOR KEEPER
// ==============================

// ValidatorKeeper is the interface used to confirm a CreateIssuer caller is
// an active validator registered for the domain it's trying to issue an
// issuer for.
type ValidatorKeeper interface {
	IsValidator(ctx sdk.Context, address string) bool
	GetValidator(ctx sdk.Context, address string) (validatortypes.Validator, error)
}

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
}