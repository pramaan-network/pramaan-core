// Package types (this file) registers the x/docreg module's message types
// with the proto-based interface registry, so transactions containing them
// can be decoded and routed to the msg server.
package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers MsgTransferDocument, MsgRegisterDocument, and
// MsgUpdateParams as valid sdk.Msg implementations, and registers this
// module's Msg service descriptor for gRPC reflection.
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTransferDocument{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterDocument{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
