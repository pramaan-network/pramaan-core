// Package keeper implements the x/docreg module's state storage: the
// Documents map plus three secondary indexes (by content hash, by owner, by
// issuer) that keep lookups off a full table scan. This file defines the
// Keeper type, its constructor, and the core SetDocument/GetDocumentByID/
// IsHashExists methods that every message and query handler builds on.
package keeper

import (
	"fmt"

	"pramaan/x/docreg/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"errors"
)

// Keeper holds this module's dependencies and store accessors. Documents
// are the primary record; HashIndex/OwnerIndex/IssuerIndex are secondary
// indexes maintained alongside every write to Documents (see SetDocument)
// so queries like "all documents owned by X" don't require a full scan.
type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte

	Schema      collections.Schema
	ParamsStore collections.Item[types.Params]
	bankKeeper  types.BankKeeper

	Documents collections.Map[string, types.Document]
	HashIndex collections.Map[string, string]

	OwnerIndex  collections.Map[string, types.StringList]
	IssuerIndex collections.Map[string, types.StringList]
}

// NewKeeper constructs a Keeper, builds its collections schema (Params +
// Documents + the three indexes), and panics if the given authority address
// doesn't decode — a fail-fast check at app-wiring time.
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

		ParamsStore: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			codec.CollValue[types.Params](cdc),
		),

		Documents: collections.NewMap(
			sb,
			collections.NewPrefix("documents"),
			"documents",
			collections.StringKey,
			codec.CollValue[types.Document](cdc),
		),

		HashIndex: collections.NewMap(
			sb,
			collections.NewPrefix("hash_index"),
			"hash_index",
			collections.StringKey,
			collections.StringValue,
		),

		OwnerIndex: collections.NewMap(
			sb,
			collections.NewPrefix("owner_index"),
			"owner_index",
			collections.StringKey,
			codec.CollValue[types.StringList](cdc),
		),

		IssuerIndex: collections.NewMap(
			sb,
			collections.NewPrefix("issuer_index"),
			"issuer_index",
			collections.StringKey,
			codec.CollValue[types.StringList](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the address authorized to submit MsgUpdateParams for
// this module (the gov module account by default — see module/depinject.go
// and SECURITY_CHANGELOG.md #21 for the bug this used to have).
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// ==============================
// CORE: STORE DOCUMENT
// ==============================

// SetDocument stores a new document and updates all three secondary
// indexes (hash, owner, issuer) to match. Rejects a duplicate ID or a
// duplicate content hash outright rather than overwriting — used by both
// RegisterDocument (new documents) and InitGenesis (restoring documents
// from a genesis export), which is why it does full index maintenance
// rather than assuming Documents.Set alone is enough.
func (k Keeper) SetDocument(ctx sdk.Context, doc types.Document) error {

	// Enforce a unique ID — prevents silently overwriting an existing
	// document (this was bug #3; see SECURITY_CHANGELOG.md).
	idExists, err := k.Documents.Has(ctx, doc.Id)
	if err != nil {
		return err
	}
	if idExists {
		return errorsmod.Wrapf(types.ErrDocumentIDExists, "document id already registered")
	}

	// Enforce a unique content hash — a different document ID should never
	// be allowed to register the same underlying content twice.
	has, _ := k.HashIndex.Has(ctx, doc.Hash)
	if has {
		return errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "hash already registered")
	}

	// Store the document itself.
	if err := k.Documents.Set(ctx, doc.Id, doc); err != nil {
		return err
	}

	// Store the hash -> ID index entry.
	if err := k.HashIndex.Set(ctx, doc.Hash, doc.Id); err != nil {
		return err
	}

	// ==============================
	// Owner index update: append this document's ID to its owner's list.
	// ==============================

	ownerList, err := k.OwnerIndex.Get(ctx, doc.Owner)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return err
		}
		ownerList = types.StringList{Items: []string{}}
	}

	ownerList.Items = append(ownerList.Items, doc.Id)

	if err := k.OwnerIndex.Set(ctx, doc.Owner, ownerList); err != nil {
		return err
	}

	// ==============================
	// Issuer index update: append this document's ID to its issuer's list.
	// ==============================

	issuerList, err := k.IssuerIndex.Get(ctx, doc.Issuer)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return err
		}
		issuerList = types.StringList{Items: []string{}}
	}

	issuerList.Items = append(issuerList.Items, doc.Id)

	if err := k.IssuerIndex.Set(ctx, doc.Issuer, issuerList); err != nil {
		return err
	}

	return nil
}

// ==============================
// GET DOCUMENT
// ==============================

// GetDocumentByID looks up a document by its ID. The bool return reports
// whether it was found — any Get error (including not-found) is treated as
// "not found" here, so a real backend error would be indistinguishable from
// a missing document to callers of this method.
func (k Keeper) GetDocumentByID(ctx sdk.Context, id string) (types.Document, bool) {
	doc, err := k.Documents.Get(ctx, id)
	if err != nil {
		return types.Document{}, false
	}
	return doc, true
}

// ==============================
// HASH CHECK
// ==============================

// IsHashExists reports whether a document with the given content hash is
// already registered, via HashIndex.
func (k Keeper) IsHashExists(ctx sdk.Context, hash string) bool {
	has, err := k.HashIndex.Has(ctx, hash)
	if err != nil {
		return false
	}
	return has
}
