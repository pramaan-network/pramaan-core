// Package keeper (this file) declares the x/docreg msgServer type and its
// constructor. The actual message handlers (RegisterDocument,
// TransferDocument, UpdateParams) are implemented in their own files
// (msg_server_register_document.go, msg_server_transfer_document.go,
// msg_update_params.go) to keep each handler's logic isolated.
package keeper

import (
	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/docreg/types"
)

// msgServer adapts a Keeper (plus its authority/issuer dependencies) to the
// generated types.MsgServer interface.
type msgServer struct {
	Keeper
	authorityKeeper authoritytypes.AuthorityKeeper
	issuerKeeper    types.IssuerKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface,
// wiring in the cross-module keepers RegisterDocument needs to validate
// issuer identity/domain (issuerKeeper) and any future authority checks
// (authorityKeeper).
func NewMsgServerImpl(
	keeper Keeper,
	authorityKeeper authoritytypes.AuthorityKeeper,
	issuerKeeper types.IssuerKeeper,
) types.MsgServer {
	return &msgServer{
		Keeper:          keeper,
		authorityKeeper: authorityKeeper,
		issuerKeeper:    issuerKeeper,
	}
}

var _ types.MsgServer = msgServer{}
