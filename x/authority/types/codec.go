package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	msgservice "github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateParams{}, "authority/UpdateParams", nil)
	cdc.RegisterConcrete(&MsgAddAuthority{}, "authority/AddAuthority", nil)
}

// RegisterInterfaces registers the module interfaces
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgAddAuthority{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}