// Package types (this file) registers the x/validatorreg module's message
// types with the proto-based interface registry. See
// SECURITY_CHANGELOG.md #13: MsgApplyValidator/MsgApproveValidator/
// MsgActivateValidator were missing from this registration until that fix,
// making all three permanently undecodable from any transaction.
package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers every Msg type this module defines — the
// direct ROOT-only path (AddValidator/RemoveValidator), the standard
// UpdateParams, and the AUTHORITY-quorum proposal flow
// (ApplyValidator/ApproveValidator/ActivateValidator) — as valid sdk.Msg
// implementations, and registers this module's Msg service descriptor for
// gRPC reflection.
func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {
	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveValidator{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddValidator{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgApplyValidator{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgApproveValidator{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgActivateValidator{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
