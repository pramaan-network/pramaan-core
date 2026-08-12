// Package keeper (this file) implements the x/authority MsgServer: the
// AddAuthority handler that lets an existing ROOT/VALIDATOR account grant a
// role to a new address. (UpdateParams — the other message this module
// defines — is implemented separately in msg_update_params.go and is
// unconditionally disabled.)
package keeper

import (
	"fmt"
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"pramaan/x/authority/types"
)

// msgServer adapts a Keeper to the generated types.MsgServer interface.
type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// backed by the given Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// AddAuthority handles MsgAddAuthority: it registers a new address under a
// role (AUTHORITY/VALIDATOR/ISSUER), gated by the creator's own role.
//
// Role-grant rules enforced here:
//   - ROOT can never be created through this message (RoleRoot is rejected
//     outright) — ROOT accounts only ever come from genesis.
//   - AUTHORITY can only be granted by an existing ROOT.
//   - VALIDATOR cannot be granted here at all — validators must go through
//     x/validatorreg's proposal-and-approval flow (ApplyValidator /
//     ApproveValidator / ActivateValidator), not a direct grant.
//   - ISSUER can only be granted by an existing VALIDATOR.
//   - Any other role string is rejected as an invalid request.
func (k msgServer) AddAuthority(
	goCtx context.Context,
	msg *types.MsgAddAuthority,
) (*types.MsgAddAuthorityResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 0. Bound free-form inputs before persisting (state-bloat / DoS guard —
	// see types.Max*Len).
	if len(msg.Address) > types.MaxAddressLen {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "address exceeds max length %d", types.MaxAddressLen)
	}
	if len(msg.PubKey) > types.MaxPubKeyLen {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "pubkey exceeds max length %d", types.MaxPubKeyLen)
	}
	if len(msg.Metadata) > types.MaxMetadataLen {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "metadata exceeds max length %d", types.MaxMetadataLen)
	}

	// 1. Get creator authority
	creatorAuth, found := k.Keeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "creator is not an authority")
	}

	// 2. Prevent ROOT creation
	if msg.Role == types.RoleRoot {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot create ROOT authority")
	}

	// 3. Role-based validation
	switch msg.Role {

	case types.RoleAuthority:
		if creatorAuth.Role != types.RoleRoot {
			return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only ROOT can add AUTHORITY")
		}

	case types.RoleValidator:
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "validators must use proposal system")

	case types.RoleIssuer:
		if creatorAuth.Role != types.RoleValidator {
			return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only VALIDATOR can add ISSUER")
		}

	default:
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid role")
	}

	// 4. Check if already exists
	_, exists := k.Keeper.GetAuthority(ctx, msg.Address)
	if exists {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "authority already exists")
	}

	// 5. Create authority
	newAuthority := types.Authority{
		Address: msg.Address,
		PubKey:  msg.PubKey,
		Role:    msg.Role,
	}

	// 6. Store
	k.Keeper.SetAuthority(ctx, newAuthority)

	// =====================================================
	// 🔥 STANDARDIZED EVENT (ENGINE READY)
	// =====================================================

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"authority.created", // ✅ standardized naming

			sdk.NewAttribute("module", "authority"),
			sdk.NewAttribute("action", "create"),

			sdk.NewAttribute("address", msg.Address),
			sdk.NewAttribute("role", msg.Role),
			sdk.NewAttribute("creator", msg.Creator),

			// 🔥 NEW (GENERIC ENGINE SUPPORT)
			sdk.NewAttribute("metadata", msg.Metadata),

			// 🔥 AUDIT SUPPORT
			sdk.NewAttribute("block_time", ctx.BlockTime().String()),
			sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
		),
	)

	return &types.MsgAddAuthorityResponse{}, nil
}