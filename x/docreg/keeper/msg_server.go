package keeper

import (
	"pramaan/x/docreg/types"
)

type msgServer struct {
	Keeper
}

// constructor
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// interface check
var _ types.MsgServer = msgServer{}