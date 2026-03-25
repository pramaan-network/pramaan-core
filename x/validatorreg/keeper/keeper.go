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
	authority    []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	Validators collections.Map[string, types.Validator]

	Proposals     collections.Map[uint64, types.ValidatorProposal]
	ProposalCount collections.Item[uint64]
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
			codec.CollValue[types.Validator](cdc),
		),

		Proposals: collections.NewMap(
			sb,
			collections.NewPrefix(1),
			"proposals",
			collections.Uint64Key,
			codec.CollValue[types.ValidatorProposal](cdc),
		),

		ProposalCount: collections.NewItem(
			sb,
			collections.NewPrefix(2),
			"proposal_count",
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

// GetAuthority returns module authority (NOT USED in your design)
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// ✅ Add Validator
func (k Keeper) AddValidator(ctx sdk.Context, addr string, domain string) error {
	val := types.Validator{
		Address: addr,
		Active:  true,
		Domain:  domain,
	}
	return k.Validators.Set(ctx, addr, val)
}

// ✅ Remove Validator
func (k Keeper) RemoveValidator(ctx sdk.Context, addr string) error {
	return k.Validators.Remove(ctx, addr)
}

// ✅ Check Validator
func (k Keeper) IsValidator(ctx sdk.Context, addr string) bool {
	val, err := k.Validators.Get(ctx, addr)
	if err != nil {
		return false
	}
	return val.Active
}

func (k Keeper) GetValidator(ctx sdk.Context, addr string) (types.Validator, error) {
	return k.Validators.Get(ctx, addr)
}