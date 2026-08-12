// Package keeper (this file) declares the x/pramaan msgServer type and its
// constructor. UpdateParams (this module's only message) is implemented in
// msg_update_params.go — no other message handlers exist here, so this file
// has no handler bodies of its own.
package keeper

import (
	"pramaan/x/pramaan/types"
)

// msgServer adapts a Keeper to the generated types.MsgServer interface.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
