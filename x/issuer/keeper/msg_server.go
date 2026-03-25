package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"pramaan/x/issuer/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = &msgServer{}

// ==============================
// CREATE ISSUER
// ==============================

func (k msgServer) CreateIssuer(
	goCtx context.Context,
	msg *types.MsgCreateIssuer,
) (*types.MsgCreateIssuerResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔴 1. CHECK AUTHORITY ROLE = VALIDATOR
	auth, found := k.authorityKeeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, fmt.Errorf("creator not found in authority")
	}

	if auth.Role != "VALIDATOR" {
		return nil, fmt.Errorf("only VALIDATOR can create issuer")
	}

	// 🔴 2. CHECK VALIDATOR EXISTS + ACTIVE
	val, err := k.validatorKeeper.GetValidator(ctx, msg.Creator)
	if err != nil {
		return nil, fmt.Errorf("validator not found")
	}

	if !val.Active {
		return nil, fmt.Errorf("validator not active")
	}

	// 🔴 3. DOMAIN CHECK
	if val.Domain != msg.Domain {
		return nil, fmt.Errorf("validator not allowed for this domain")
	}

	// 🔴 4. CHECK ISSUER EXISTS
	_, err = k.Keeper.Issuers.Get(ctx, msg.Address)
	if err == nil {
		return nil, fmt.Errorf("issuer already exists")
	}

	// 🔴 5. CREATE ISSUER
	issuer := types.Issuer{
		Id:        msg.Address,
		Address:   msg.Address,
		Validator: msg.Creator,
		Domain:    msg.Domain,
		Active:    true,
	}

	if err := k.Keeper.Issuers.Set(ctx, msg.Address, issuer); err != nil {
		return nil, err
	}

	// 🔥 EVENT
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"issuer_added",
			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("validator", msg.Creator),
			sdk.NewAttribute("domain", msg.Domain),
		),
	)

	return &types.MsgCreateIssuerResponse{}, nil
}

// ==============================
// REVOKE ISSUER
// ==============================

func (k msgServer) RevokeIssuer(
	goCtx context.Context,
	msg *types.MsgRevokeIssuer,
) (*types.MsgRevokeIssuerResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	issuer, err := k.Keeper.Issuers.Get(ctx, msg.Address)
	if err != nil {
		return nil, fmt.Errorf("issuer not found")
	}

	// ❌ only creator validator
	if issuer.Validator != msg.Creator {
		return nil, fmt.Errorf("not allowed")
	}

	issuer.Active = false

	if err := k.Keeper.Issuers.Set(ctx, msg.Address, issuer); err != nil {
		return nil, err
	}

	// 🔥 event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"issuer_revoked",
			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("creator", msg.Creator),
		),
	)

	return &types.MsgRevokeIssuerResponse{}, nil
}

