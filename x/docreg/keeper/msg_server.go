package keeper

import (
	"pramaan/x/docreg/types"
	authoritytypes "pramaan/x/authority/types"
)

type msgServer struct {
	Keeper
	authorityKeeper authoritytypes.AuthorityKeeper
}

// ✅ FIXED CONSTRUCTOR
func NewMsgServerImpl(
	keeper Keeper,
	authorityKeeper authoritytypes.AuthorityKeeper,
) types.MsgServer {
	return &msgServer{
		Keeper:          keeper,
		authorityKeeper: authorityKeeper,
	}
}

var _ types.MsgServer = msgServer{}