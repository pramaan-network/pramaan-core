// Package types (this file) registers the x/authority module's message
// types with both the legacy Amino codec and the modern proto-based
// interface registry, so transactions containing these messages can be
// signed, decoded, and routed.
package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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
