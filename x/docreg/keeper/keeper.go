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

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte

	Schema      collections.Schema
	ParamsStore collections.Item[types.Params]
	bankKeeper  types.BankKeeper

	Documents  collections.Map[string, types.Document]
	HashIndex  collections.Map[string, string]

	OwnerIndex  collections.Map[string, types.StringList]
	IssuerIndex collections.Map[string, types.StringList]
}

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

// GetAuthority returns module authority
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// ==============================
// 🚀 CORE: STORE DOCUMENT
// ==============================

func (k Keeper) SetDocument(ctx sdk.Context, doc types.Document) error {

	// 🚫 Enforce unique hash
	has, _ := k.HashIndex.Has(ctx, doc.Hash)
	if has {
		return errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "hash already registered")
	}

	// ✅ Store document
	if err := k.Documents.Set(ctx, doc.Id, doc); err != nil {
		return err
	}

	// ✅ Store hash index
	if err := k.HashIndex.Set(ctx, doc.Hash, doc.Id); err != nil {
		return err
	}

	// ==============================
	// 🔥 OWNER INDEX UPDATE
	// ==============================

	ownerList, err := k.OwnerIndex.Get(ctx, doc.Owner)
	if err != nil {
		ownerList = types.StringList{Items: []string{}}
	}

	ownerList.Items = append(ownerList.Items, doc.Id)

	if err := k.OwnerIndex.Set(ctx, doc.Owner, ownerList); err != nil {
		return err
	}

	// ==============================
	// 🔥 ISSUER INDEX UPDATE
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
// 📦 GET DOCUMENT
// ==============================

func (k Keeper) GetDocumentByID(ctx sdk.Context, id string) (types.Document, bool) {
	doc, err := k.Documents.Get(ctx, id)
	if err != nil {
		return types.Document{}, false
	}
	return doc, true
}

// ==============================
// 🔍 HASH CHECK
// ==============================

func (k Keeper) IsHashExists(ctx sdk.Context, hash string) bool {
	has, err := k.HashIndex.Has(ctx, hash)
	if err != nil {
		return false
	}
	return has
}