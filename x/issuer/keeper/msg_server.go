// Package keeper (this file) implements the x/issuer MsgServer:
// CreateIssuer and RevokeIssuer.
package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/issuer/types"
)

// msgServer adapts a Keeper to the generated types.MsgServer interface.
// Unlike docreg/validatorreg, the authority/validator keepers this server
// needs aren't passed to the constructor separately — they're already
// embedded fields on Keeper itself (see keeper.go), so NewMsgServerImpl
// only needs the one Keeper argument.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// backed by the given Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = &msgServer{}

// ==============================
// CREATE ISSUER
// ==============================

// CreateIssuer handles MsgCreateIssuer. Authorization chain, checked in
// order: the creator must hold the VALIDATOR role in x/authority; the
// creator must also be a registered, active validator in x/validatorreg;
// the validator's registered domain must match the requested issuer domain
// (a validator for "education" can't create an issuer for "land_property");
// and the target issuer address must not already exist. Only after all
// four checks pass is the new Issuer record created.
func (k msgServer) CreateIssuer(
	goCtx context.Context,
	msg *types.MsgCreateIssuer,
) (*types.MsgCreateIssuerResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔴 0. BOUND FREE-FORM INPUTS (state-bloat / DoS guard — see types.Max*Len)
	if len(msg.Domain) > types.MaxDomainLen {
		return nil, fmt.Errorf("domain exceeds max length %d", types.MaxDomainLen)
	}
	if len(msg.Address) > types.MaxAddressLen {
		return nil, fmt.Errorf("address exceeds max length %d", types.MaxAddressLen)
	}

	// 🔴 1. CHECK AUTHORITY ROLE = VALIDATOR
	auth, found := k.authorityKeeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, fmt.Errorf("creator not found in authority")
	}

	if auth.Role != authoritytypes.RoleValidator {
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

// RevokeIssuer handles MsgRevokeIssuer: deactivates (sets Active = false)
// an issuer, but only its own creating validator may do so — there's no
// path here for a higher role (e.g. ROOT or AUTHORITY) to revoke a
// misbehaving issuer if its creating validator is unavailable or
// compromised, which is worth keeping in mind for incident response.
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

