package keeper

import (
	"fmt"
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"pramaan/x/authority/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// AddAuthority handles adding new authority
func (k msgServer) AddAuthority(
	goCtx context.Context,
	msg *types.MsgAddAuthority,
) (*types.MsgAddAuthorityResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. Get creator authority
	creatorAuth, found := k.Keeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "creator is not an authority")
	}

	// 2. Prevent ROOT creation
	if msg.Role == "ROOT" {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot create ROOT authority")
	}

	// 3. Role-based validation
	switch msg.Role {

	case "AUTHORITY":
		if creatorAuth.Role != "ROOT" {
			return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only ROOT can add AUTHORITY")
		}

	case "VALIDATOR":
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "validators must use proposal system")

	case "ISSUER":
		if creatorAuth.Role != "VALIDATOR" {
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