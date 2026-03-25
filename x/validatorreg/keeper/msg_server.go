package keeper

import (
	"context"
	"fmt"

	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/validatorreg/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
	pramaanKeeper   types.PramaanKeeper
	authorityKeeper authoritytypes.AuthorityKeeper
}

// ✅ MUST return POINTER
func NewMsgServerImpl(
	keeper Keeper,
	pramaanKeeper types.PramaanKeeper,
	authorityKeeper authoritytypes.AuthorityKeeper,
) types.MsgServer {
	return &msgServer{
		Keeper:          keeper,
		pramaanKeeper:   pramaanKeeper,
		authorityKeeper: authorityKeeper,
	}
}

func (k msgServer) ApplyValidator(
	goCtx context.Context,
	msg *types.MsgApplyValidator,
) (*types.MsgApplyValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔹 1. get current proposal count
	count, err := k.Keeper.ProposalCount.Get(ctx)
	if err != nil {
		count = 0
	}

	newID := count + 1

	// 🔹 2. create proposal
	proposal := types.ValidatorProposal{
		Id:        newID,
		Applicant: msg.Creator,
		Domain:    msg.Domain,
		Data:      msg.Data,
		Approvals: []string{},
		Status:    "PENDING",
	}

	// 🔹 3. store proposal
	if err := k.Keeper.Proposals.Set(ctx, newID, proposal); err != nil {
		return nil, err
	}

	// 🔹 4. update counter
	if err := k.Keeper.ProposalCount.Set(ctx, newID); err != nil {
		return nil, err
	}

	// 🔹 5. emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_proposal_created",
			sdk.NewAttribute("id", fmt.Sprintf("%d", newID)),
			sdk.NewAttribute("applicant", msg.Creator),
			sdk.NewAttribute("domain", msg.Domain),
		),
	)

	return &types.MsgApplyValidatorResponse{}, nil
}

func (k msgServer) ApproveValidator(
	goCtx context.Context,
	msg *types.MsgApproveValidator,
) (*types.MsgApproveValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔹 1. check authority
	auth, found := k.authorityKeeper.GetAuthority(ctx, msg.Creator)
	if !found {
		return nil, fmt.Errorf("not an authority")
	}

	if auth.Role != "AUTHORITY" {
		return nil, fmt.Errorf("only AUTHORITY can approve")
	}

	// 🔹 2. get proposal
	proposal, err := k.Keeper.Proposals.Get(ctx, msg.ProposalId)
	if err != nil {
		return nil, fmt.Errorf("proposal not found")
	}

	// 🔹 3. check already approved
	for _, a := range proposal.Approvals {
		if a == msg.Creator {
			return nil, fmt.Errorf("already approved")
		}
	}

	// 🔹 4. append approval
	proposal.Approvals = append(proposal.Approvals, msg.Creator)

	// 🔹 5. threshold check (3)
	if len(proposal.Approvals) >= 3 && proposal.Status != "APPROVED" {
		proposal.Status = "APPROVED"
	}

	// 🔹 6. save
	if err := k.Keeper.Proposals.Set(ctx, msg.ProposalId, proposal); err != nil {
		return nil, err
	}

	// 🔹 7. emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_approved",
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute("approver", msg.Creator),
		),
	)

	return &types.MsgApproveValidatorResponse{}, nil
}

func (k msgServer) ActivateValidator(
	goCtx context.Context,
	msg *types.MsgActivateValidator,
) (*types.MsgActivateValidatorResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// 🔹 1. get proposal
	proposal, err := k.Keeper.Proposals.Get(ctx, msg.ProposalId)
	if err != nil {
		return nil, fmt.Errorf("proposal not found")
	}

	// 🔹 2. must be approved
	if proposal.Status != "APPROVED" {
		return nil, fmt.Errorf("proposal not approved")
	}

	// 🔹 3. only applicant can activate
	if proposal.Applicant != msg.Creator {
		return nil, fmt.Errorf("only applicant can activate")
	}

	// 🔹 4. create validator
	err = k.Keeper.AddValidator(ctx, proposal.Applicant, proposal.Domain)
	if err != nil {
		return nil, err
	}

	newAuthority := authoritytypes.Authority{
		Address: proposal.Applicant,
		PubKey:  "validator", // temp (can improve later)
		Role:    "VALIDATOR",
	}

	k.authorityKeeper.SetAuthority(ctx, newAuthority)

	// 🔹 5. mark proposal as used
	proposal.Status = "ACTIVATED"
	if err := k.Keeper.Proposals.Set(ctx, msg.ProposalId, proposal); err != nil {
		return nil, err
	}

	// 🔹 6. emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"validator_activated",
			sdk.NewAttribute("address", proposal.Applicant),
			sdk.NewAttribute("domain", proposal.Domain),
			sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),
		),
	)

	return &types.MsgActivateValidatorResponse{}, nil
}

var _ types.MsgServer = &msgServer{}
