package keeper

import (
	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/docreg/types"
)

type msgServer struct {
	Keeper
	authorityKeeper authoritytypes.AuthorityKeeper
	issuerKeeper    types.IssuerKeeper
}

// ✅ FIXED CONSTRUCTOR
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
