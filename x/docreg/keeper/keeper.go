package keeper

import (
	"fmt"
	"pramaan/x/docreg/types"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte
	Schema       collections.Schema
	ParamsStore collections.Item[types.Params]
	bankKeeper   types.BankKeeper
	Documents    collections.Map[string, types.Document]
	HashIndex    collections.Map[string, string]
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
		ParamsStore:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
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
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// Store Document
func (k Keeper) SetDocument(ctx sdk.Context, doc types.Document) error {

	// 🚫 HARD ENFORCEMENT HERE
	has, _ := k.HashIndex.Has(ctx, doc.Hash)
	if has {
		return errorsmod.Wrapf(types.ErrDocumentAlreadyExists, "hash already registered")
	}

	// store document
	if err := k.Documents.Set(ctx, doc.Id, doc); err != nil {
		return err
	}

	// store hash index
	return k.HashIndex.Set(ctx, doc.Hash, doc.Id)
}

// Get Document
func (k Keeper) GetDocumentByID(ctx sdk.Context, id string) (types.Document, bool) {
	doc, err := k.Documents.Get(ctx, id)
	if err != nil {
		return types.Document{}, false
	}
	return doc, true
}

// Hash Check
func (k Keeper) IsHashExists(ctx sdk.Context, hash string) bool {
	has, err := k.HashIndex.Has(ctx, hash)
	if err != nil {
		return false
	}
	return has
}