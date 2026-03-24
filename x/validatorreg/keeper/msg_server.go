package keeper

import (
	"pramaan/x/validatorreg/types"
	authoritytypes "pramaan/x/authority/types"
)

type msgServer struct {
	Keeper
	pramaanKeeper types.PramaanKeeper
	authorityKeeper authoritytypes.AuthorityKeeper
}

// ✅ MUST return POINTER
func NewMsgServerImpl(
	keeper Keeper,
	pramaanKeeper types.PramaanKeeper,
	authorityKeeper authoritytypes.AuthorityKeeper,
) types.MsgServer {
	return &msgServer{
		Keeper: keeper,
		pramaanKeeper: pramaanKeeper,
		authorityKeeper: authorityKeeper,
	}
}

// ✅ MUST be POINTER (CRITICAL FIX)
var _ types.MsgServer = &msgServer{}