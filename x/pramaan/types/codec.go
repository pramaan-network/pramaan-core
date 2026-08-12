// Package types (this file) registers the x/pramaan module's message types
// (just MsgUpdateParams — this module defines no other messages) with the
// proto-based interface registry.
package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers MsgUpdateParams as a valid sdk.Msg
// implementation, and registers this module's Msg service descriptor for
// gRPC reflection.
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
