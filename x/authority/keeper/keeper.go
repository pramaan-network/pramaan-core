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

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	bankKeeper types.BankKeeper
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
		Params:     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// SetAuthority stores an authority
func (k Keeper) SetAuthority(ctx sdk.Context, authority types.Authority) {
	store := k.storeService.OpenKVStore(ctx)

	bz := k.cdc.MustMarshal(&authority)

	store.Set(types.GetAuthorityKey(authority.Address), bz)
}

// GetAuthority retrieves an authority by address
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

// IsAuthority checks if address is registered authority
func (k Keeper) IsAuthority(ctx sdk.Context, address string) bool {
	_, found := k.GetAuthority(ctx, address)
	return found
}