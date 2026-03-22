package keeper

import (
	"pramaan/x/validatorreg/types"
)

type msgServer struct {
	Keeper
	pramaanKeeper types.PramaanKeeper
}

// ✅ MUST return POINTER
func NewMsgServerImpl(keeper Keeper, pk types.PramaanKeeper) types.MsgServer {
	return &msgServer{
		Keeper:         keeper,
		pramaanKeeper:  pk,
	}
}

// ✅ MUST be POINTER (CRITICAL FIX)
var _ types.MsgServer = &msgServer{}