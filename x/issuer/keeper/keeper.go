// Package keeper implements the x/issuer module's state storage: the
// registry of issuers, each tied to the VALIDATOR account that created it
// and scoped to a single domain. This file defines the Keeper type, its
// constructor, and the GetIssuer accessor that docreg's RegisterDocument
// depends on via types.IssuerKeeper.
package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"pramaan/x/issuer/types"
)

// Keeper holds this module's dependencies and store accessors, plus the
// cross-module authority/validator keepers CreateIssuer needs for its
// role/domain checks.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	Issuers collections.Map[string, types.Issuer]

	authorityKeeper types.AuthorityKeeper
	validatorKeeper types.ValidatorKeeper
}

// NewKeeper constructs a Keeper, builds its collections schema, and panics
// if the given authority address doesn't decode — a fail-fast check at
// app-wiring time.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	authorityKeeper types.AuthorityKeeper,
	validatorKeeper types.ValidatorKeeper,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),

		Issuers: collections.NewMap(
		sb,
		collections.NewPrefix("issuer"),
		"issuers",
		collections.StringKey,
		codec.CollValue[types.Issuer](cdc),
	),

		authorityKeeper: authorityKeeper,
		validatorKeeper: validatorKeeper,

	}



	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the address authorized to submit MsgUpdateParams for
// this module (the gov module account by default — see module/depinject.go).
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// GetIssuer looks up an issuer by address. This is the method docreg's
// RegisterDocument calls (via types.IssuerKeeper) to confirm a document's
// declared issuer is real before allowing registration. Note: any error
// from the underlying Get (not just "not found") is collapsed into the
// same generic "issuer not found" message — a transient store error would
// look identical to a genuinely unregistered issuer to callers.
func (k Keeper) GetIssuer(ctx sdk.Context, address string) (types.Issuer, error) {
	issuer, err := k.Issuers.Get(ctx, address)
	if err != nil {
		return types.Issuer{}, fmt.Errorf("issuer not found")
	}
	return issuer, nil
}