// Package keeper (this file) implements the x/validatorreg
// AUTHORITY-quorum validator-admission workflow: ApplyValidator (any
// account proposes itself) -> ApproveValidator (distinct AUTHORITY-role
// accounts vote, see ValidatorApprovalThreshold below) -> ActivateValidator
// (the original applicant, once approved, finalizes registration). This is
// separate from the direct ROOT-only AddValidator/RemoveValidator path in
// msg_server_add_validator.go/msg_server_remove_validator.go.
package keeper

import (
	"context"
	"fmt"

	authoritytypes "pramaan/x/authority/types"
	"pramaan/x/validatorreg/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorApprovalThreshold is the number of distinct AUTHORITY approvals a
// validator proposal needs before it flips to APPROVED.
//
// NOTE: this should really be a governed on-chain param (types.Params) so it can be
// changed via MsgUpdateParams instead of a binary upgrade. That requires adding a
// field to the validatorreg Params proto message and regenerating the *.pb.go code
// (buf/ignite codegen), which isn't available in this environment right now. Tracked
// as a follow-up; for now this constant is at least a single named, documented place
// instead of a bare magic number.
const ValidatorApprovalThreshold = 3

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

	// Bound free-form inputs before persisting a proposal (state-bloat / DoS
	// guard — see types.Max*Len).
	if len(msg.Domain) > types.MaxDomainLen {
		return nil, fmt.Errorf("domain exceeds max length %d", types.MaxDomainLen)
	}
	if len(msg.Data) > types.MaxDataLen {
		return nil, fmt.Errorf("data exceeds max length %d", types.MaxDataLen)
	}

	// ==============================
	// 🔥 0. PREVENT DUPLICATE APPLY
	// ==============================
	err := k.Keeper.Proposals.Walk(ctx, nil, func(id uint64, p types.ValidatorProposal) (bool, error) {

		// same applicant already has pending/approved proposal
		if p.Applicant == msg.Creator && (p.Status == types.ProposalStatusPending || p.Status == types.ProposalStatusApproved) {
			return true, fmt.Errorf("validator already has active proposal")
		}

		return false, nil
	})

	if err != nil {
		return nil, err
	}

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
		Status:    types.ProposalStatusPending,
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
		"validator.proposal.created",

		sdk.NewAttribute("module", "validatorreg"),
		sdk.NewAttribute("action", "create_proposal"),

		sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", newID)),
		sdk.NewAttribute("applicant", msg.Creator),
		sdk.NewAttribute("domain", msg.Domain),

		sdk.NewAttribute("status", types.ProposalStatusPending),

		// 🔥 ENGINE SUPPORT
		sdk.NewAttribute("metadata", msg.Data),

		// 🔥 AUDIT
		sdk.NewAttribute("block_time", ctx.BlockTime().String()),
		sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
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

	if auth.Role != authoritytypes.RoleAuthority {
		return nil, fmt.Errorf("only AUTHORITY can approve")
	}

	// 🔹 2. get proposal
	proposal, err := k.Keeper.Proposals.Get(ctx, msg.ProposalId)
	if err != nil {
		return nil, fmt.Errorf("proposal not found")
	}

	// 🔹 2b. only a still-pending proposal can accept approvals. Without this
	// guard an AUTHORITY could keep appending approvals to a proposal that is
	// already APPROVED or ACTIVATED — noise at best, and it leaves the
	// terminal states mutable, which any later logic keying off Approvals
	// would have to defend against. Reject anything not PENDING outright.
	if proposal.Status != types.ProposalStatusPending {
		return nil, fmt.Errorf("proposal is not pending (status %s)", proposal.Status)
	}

	// 🔹 3. check already approved
	for _, a := range proposal.Approvals {
		if a == msg.Creator {
			return nil, fmt.Errorf("already approved")
		}
	}

	// 🔹 4. append approval
	proposal.Approvals = append(proposal.Approvals, msg.Creator)

	// 🔹 5. threshold check
	if len(proposal.Approvals) >= ValidatorApprovalThreshold && proposal.Status != types.ProposalStatusApproved {
		proposal.Status = types.ProposalStatusApproved
	}

	// 🔹 6. save
	if err := k.Keeper.Proposals.Set(ctx, msg.ProposalId, proposal); err != nil {
		return nil, err
	}

	// 🔹 7. emit event
	ctx.EventManager().EmitEvent(
	sdk.NewEvent(
		"validator.proposal.approved",

		sdk.NewAttribute("module", "validatorreg"),
		sdk.NewAttribute("action", "approve"),

		sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),
		sdk.NewAttribute("approver", msg.Creator),

		sdk.NewAttribute("total_approvals", fmt.Sprintf("%d", len(proposal.Approvals))),
		sdk.NewAttribute("status", proposal.Status),

		// 🔥 AUDIT
		sdk.NewAttribute("block_time", ctx.BlockTime().String()),
		sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
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
	if proposal.Status != types.ProposalStatusApproved {
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
		Role:    authoritytypes.RoleValidator,
	}

	k.authorityKeeper.SetAuthority(ctx, newAuthority)

	// 🔹 5. mark proposal as used
	proposal.Status = types.ProposalStatusActivated
	if err := k.Keeper.Proposals.Set(ctx, msg.ProposalId, proposal); err != nil {
		return nil, err
	}

	// 🔹 6. emit event
	ctx.EventManager().EmitEvent(
	sdk.NewEvent(
		"validator.activated",

		sdk.NewAttribute("module", "validatorreg"),
		sdk.NewAttribute("action", "activate"),

		sdk.NewAttribute("validator", proposal.Applicant),
		sdk.NewAttribute("domain", proposal.Domain),
		sdk.NewAttribute("proposal_id", fmt.Sprintf("%d", msg.ProposalId)),

		sdk.NewAttribute("status", proposal.Status),

		// 🔥 AUDIT
		sdk.NewAttribute("block_time", ctx.BlockTime().String()),
		sdk.NewAttribute("height", fmt.Sprintf("%d", ctx.BlockHeight())),
	),
	)		

	return &types.MsgActivateValidatorResponse{}, nil
}

var _ types.MsgServer = &msgServer{}
