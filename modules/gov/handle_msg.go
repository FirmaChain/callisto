package gov

import (
	"fmt"

	"strconv"

	"github.com/forbole/callisto/v4/types"
	"github.com/forbole/callisto/v4/utils"
	eventutils "github.com/forbole/callisto/v4/utils/events"
	"github.com/rs/zerolog/log"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	juno "github.com/forbole/juno/v6/types"
)

var msgFilter = map[string]bool{
	"/cosmos.gov.v1.MsgSubmitProposal": true,
	"/cosmos.gov.v1.MsgDeposit":        true,
	"/cosmos.gov.v1.MsgVote":           true,

	"/cosmos.gov.v1beta1.MsgSubmitProposal": true,
	"/cosmos.gov.v1beta1.MsgDeposit":        true,
	"/cosmos.gov.v1beta1.MsgVote":           true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	log.Debug().Str("module", "gov").Str("hash", tx.TxHash).Uint64("height", tx.Height).Msg(fmt.Sprintf("handling gov message %s", msg.GetType()))

	switch msg.GetType() {
	case "/cosmos.gov.v1.MsgSubmitProposal":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &govtypesv1.MsgSubmitProposal{})
		return m.handleSubmitProposalEvent(tx, cosmosMsg.Proposer, eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index))
	case "/cosmos.gov.v1beta1.MsgSubmitProposal":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &govtypesv1beta1.MsgSubmitProposal{})
		return m.handleSubmitProposalEvent(tx, cosmosMsg.Proposer, eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index))

	case "/cosmos.gov.v1.MsgDeposit":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &govtypesv1.MsgDeposit{})
		return m.handleDepositEvent(tx, cosmosMsg.Depositor, eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index))
	case "/cosmos.gov.v1beta1.MsgDeposit":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &govtypesv1beta1.MsgDeposit{})
		return m.handleDepositEvent(tx, cosmosMsg.Depositor, eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index))

	case "/cosmos.gov.v1.MsgVote":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &govtypesv1.MsgVote{})
		return m.handleVoteEvent(tx, cosmosMsg.Voter, eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index))
	case "/cosmos.gov.v1beta1.MsgVote":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &govtypesv1beta1.MsgVote{})
		return m.handleVoteEvent(tx, cosmosMsg.Voter, eventutils.FindEventsByMsgIndex(sdk.StringifyEvents(tx.Events), index))
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func (m *Module) handleMsgSubmitProposal(tx *juno.Transaction, index int, msg *govtypes.MsgSubmitProposal) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyProposalID: %s", err)
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing proposal id: %s", err)
	}

	// Get the proposal
	proposal, err := m.source.Proposal(int64(tx.Height), proposalID)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	// Unpack the content
	var content govtypes.Content
	err = m.cdc.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while unpacking proposal content: %s", err)
	}

	// Store the proposal
	proposalObj := types.NewProposal(
		proposal.ProposalId,
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		proposal.GetContent(),
		proposal.Status.String(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		msg.Proposer,
	)
	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, int64(tx.Height))
	return m.db.SaveDeposits([]types.Deposit{deposit})
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func (m *Module) handleMsgDeposit(tx *juno.Tx, msg *govtypes.MsgDeposit) error {
	deposit, err := m.source.ProposalDeposit(int64(tx.Height), msg.ProposalId, msg.Depositor)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}

	return m.db.SaveDeposits([]types.Deposit{
		types.NewDeposit(msg.ProposalId, msg.Depositor, deposit.Amount, int64(tx.Height)),
	})
}

// handleMsgVote allows to properly handle a handleMsgVote
func (m *Module) handleMsgVote(tx *juno.Tx, msg *govtypes.MsgVote) error {
	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, int64(tx.Height))
	return m.db.SaveVote(vote)
}
