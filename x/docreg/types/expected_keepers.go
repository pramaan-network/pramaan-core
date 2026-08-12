// Package types (this file) declares the "expected keeper" interfaces
// x/docreg depends on: narrow interfaces onto auth/bank/issuer rather than
// their concrete keeper types, so this module doesn't need to import those
// packages directly and is easier to unit-test with mocks.
package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	issuertypes "pramaan/x/issuer/types"
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

// IssuerKeeper is the interface RegisterDocument uses (see
// keeper/msg_server_register_document.go) to confirm the message's issuer
// is a real, active issuer whose registered domain matches the document's
// declared doc_type before allowing a document to be created.
type IssuerKeeper interface {
	GetIssuer(ctx sdk.Context, address string) (issuertypes.Issuer, error)
}