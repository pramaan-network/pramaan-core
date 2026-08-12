// Package types (this file) registers the x/issuer module's message types
// with the proto-based interface registry. See SECURITY_CHANGELOG.md #12:
// MsgCreateIssuer/MsgRevokeIssuer were missing from this registration until
// that fix, which made both messages permanently undecodable from any
// transaction.
package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers MsgUpdateParams, MsgCreateIssuer, and
// MsgRevokeIssuer as valid sdk.Msg implementations, and registers this
// module's Msg service descriptor for gRPC reflection.
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateIssuer{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRevokeIssuer{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
