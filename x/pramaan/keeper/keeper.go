package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"pramaan/x/pramaan/types"

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

	Authorities collections.Map[string, types.Authority]
	Threshold   collections.Item[uint64]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
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

		Params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			codec.CollValue[types.Params](cdc),
		),

		Authorities: collections.NewMap(
			sb,
			collections.NewPrefix("authorities"),
			"authorities",
			collections.StringKey,
			codec.CollValue[types.Authority](cdc),
		),

		Threshold: collections.NewItem(
			sb,
			collections.NewPrefix("threshold"),
			"threshold",
			collections.Uint64Value,
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// ==============================
// 🔐 AUTHORITY FUNCTIONS
// ==============================

// Add Authority
func (k Keeper) AddAuthority(ctx sdk.Context, addr string) error {
	auth := types.Authority{
		Address: addr,
		Active:  true,
	}
	return k.Authorities.Set(ctx, addr, auth)
}

// Remove Authority
func (k Keeper) RemoveAuthority(ctx sdk.Context, addr string) error {
	return k.Authorities.Remove(ctx, addr)
}

// Check Authority
func (k Keeper) IsAuthority(ctx sdk.Context, addr string) bool {
	has, err := k.Authorities.Has(ctx, addr)
	if err != nil {
		return false
	}
	return has
}

// Set Threshold
func (k Keeper) SetThreshold(ctx sdk.Context, t uint64) error {
	return k.Threshold.Set(ctx, t)
}

// Get Threshold
func (k Keeper) GetThreshold(ctx sdk.Context) uint64 {
	t, err := k.Threshold.Get(ctx)
	if err != nil {
		return 1 // default fallback
	}
	return t
}