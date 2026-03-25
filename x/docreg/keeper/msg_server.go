package keeper

import (
	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/docreg/types"
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
