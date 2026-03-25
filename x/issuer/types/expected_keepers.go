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

type AuthorityKeeper interface {
	GetAuthority(ctx sdk.Context, address string) (authoritytypes.Authority, bool)
}

// ==============================
// VALIDATOR KEEPER
// ==============================

type ValidatorKeeper interface {
	IsValidator(ctx sdk.Context, address string) bool
	GetValidator(ctx sdk.Context, address string) (validatortypes.Validator, error)
}

type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
}