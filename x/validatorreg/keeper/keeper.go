package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/validatorreg/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte // ⚠️ kept for SDK compatibility (NOT USED)

	Schema collections.Schema
	Params collections.Item[types.Params]

	Validators collections.Map[string, bool] // ✅ MUST INIT
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
) Keeper {

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			codec.CollValue[types.Params](cdc),
		),

		// 🔥 CRITICAL FIX: initialize Validators store
		Validators: collections.NewMap(
			sb,
			types.ValidatorKeyPrefix, // ⚠️ must exist in types/keys.go
			"validators",
			collections.StringKey,
			collections.BoolValue,
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns module authority (NOT USED in your design)
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// ✅ Add Validator
func (k Keeper) AddValidator(ctx sdk.Context, addr string) error {
	return k.Validators.Set(ctx, addr, true)
}

// ✅ Remove Validator
func (k Keeper) RemoveValidator(ctx sdk.Context, addr string) error {
	return k.Validators.Remove(ctx, addr)
}

// ✅ Check Validator
func (k Keeper) IsValidator(ctx sdk.Context, addr string) bool {
	has, err := k.Validators.Has(ctx, addr)
	if err != nil {
		return false
	}
	return has
}