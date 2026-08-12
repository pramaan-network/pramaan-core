// Package keeper implements the x/authority module's state storage: the
// authority (role) registry that ROOT/AUTHORITY/VALIDATOR/ISSUER accounts
// live in, plus the params store every module keeps. This file defines the
// Keeper type and its core CRUD/iteration methods; msg_server.go and
// query.go build the gRPC surface on top of these.
package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"pramaan/x/authority/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper holds this module's dependencies and store accessors. Authority
// records themselves are stored directly via storeService (raw KV, see
// GetAuthority/SetAuthority below) rather than through the collections API;
// only Params uses collections.Item here.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// authority is the address allowed to submit MsgUpdateParams for this
	// module. UpdateParams is unconditionally disabled for x/authority (see
	// msg_update_params.go), so this field is currently inert, but it's kept
	// so the module follows the same NewKeeper(..., authority, ...)
	// signature as every sibling module.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	// bankKeeper is injected but not currently called by any keeper method
	// — see types/expected_keepers.go and SECURITY_CHANGELOG.md.
	bankKeeper types.BankKeeper
}

// NewKeeper constructs a Keeper, builds its collections schema, and panics
// if the given authority address doesn't decode — a fail-fast check at
// app-wiring time rather than discovering a bad address later during a
// transaction.
func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,

	bankKeeper types.BankKeeper,
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

		bankKeeper: bankKeeper,
		Params:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// SetAuthority stores (creates or overwrites) an authority record, keyed by
// its address. Callers are responsible for role/permission checks before
// calling this — it performs no authorization itself.
func (k Keeper) SetAuthority(ctx sdk.Context, authority types.Authority) {
	store := k.storeService.OpenKVStore(ctx)

	bz := k.cdc.MustMarshal(&authority)

	store.Set(types.GetAuthorityKey(authority.Address), bz)
}

// GetAuthority retrieves an authority record by address. The bool return
// reports whether a record was found — false means "not registered as an
// authority", not an error.
func (k Keeper) GetAuthority(ctx sdk.Context, address string) (types.Authority, bool) {
	store := k.storeService.OpenKVStore(ctx)

	bz, _ := store.Get(types.GetAuthorityKey(address))
	if bz == nil {
		return types.Authority{}, false
	}

	var authority types.Authority
	k.cdc.MustUnmarshal(bz, &authority)

	return authority, true
}

// IsAuthority reports whether the given address has any authority record
// registered, regardless of role. Callers that need a specific role (e.g.
// "must be ROOT") should call GetAuthority and check the Role field
// themselves — this only answers "is this address known at all".
func (k Keeper) IsAuthority(ctx sdk.Context, address string) bool {
	_, found := k.GetAuthority(ctx, address)
	return found
}

// IterateAuthorities walks every authority record in the store in key
// order, invoking cb for each. Returning true from cb stops iteration early
// (a "break", not a "continue"). Used by ExportGenesis to dump the full
// authority set.
func (k Keeper) IterateAuthorities(ctx sdk.Context, cb func(types.Authority) bool) {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {

		key := iterator.Key()

		// filter only authority keys
		if !types.IsAuthorityKey(key) {
			continue
		}

		var auth types.Authority
		k.cdc.MustUnmarshal(iterator.Value(), &auth)

		if cb(auth) {
			break
		}
	}
}