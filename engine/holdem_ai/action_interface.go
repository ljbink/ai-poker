package holdem_ai

import (
	"github.com/ljbink/ai-poker/engine/holdem"
)

// ActionType represents the type of action a player can take
type ActionType int

const (
	ActionFold ActionType = iota
	ActionCheck
	ActionCall
	ActionRaise
	ActionAllIn
)

// String returns the string representation of an action type
func (a ActionType) String() string {
	switch a {
	case ActionFold:
		return "fold"
	case ActionCheck:
		return "check"
	case ActionCall:
		return "call"
	case ActionRaise:
		return "raise"
	case ActionAllIn:
		return "all-in"
	default:
		return "unknown"
	}
}

// Action represents a player's action in the game
type Action struct {
	Type   ActionType
	Amount int // For raise/bet actions, this is the raise amount
}

// DecisionMaker interface that both human players and AI bots must implement
type DecisionMaker interface {
	// MakeDecision takes the current game and player and returns the chosen action
	MakeDecision(game *holdem.Game, player holdem.IPlayer) Action

	// GetName returns the name/identifier of this decision maker
	GetName() string

	// IsBot returns true if this is an AI bot, false if human player
	IsBot() bool
}

// ActionValidator provides methods to check if actions are valid
type ActionValidator struct{}

// NewActionValidator creates a new action validator
func NewActionValidator() *ActionValidator {
	return &ActionValidator{}
}

// IsValidAction checks if the proposed action is valid given the current game state
func (av *ActionValidator) IsValidAction(action Action, game *holdem.Game, player holdem.IPlayer) bool {
	callAmount := game.CurrentBet - player.GetBet()

	switch action.Type {
	case ActionFold:
		return true // Can always fold

	case ActionCheck:
		// Can only check if no bet to call
		return game.CurrentBet == player.GetBet()

	case ActionCall:
		// Can only call if there's a bet to call and player has chips
		return callAmount > 0 && callAmount <= player.GetChips()

	case ActionRaise:
		// Can raise if player has enough chips and raise amount is valid
		totalNeeded := callAmount + action.Amount
		return totalNeeded <= player.GetChips() && action.Amount > 0

	case ActionAllIn:
		// Can always go all-in if player has chips
		return player.GetChips() > 0

	default:
		return false
	}
}

// GetValidActions returns all valid actions for the current game state
func (av *ActionValidator) GetValidActions(game *holdem.Game, player holdem.IPlayer) []ActionType {
	var validActions []ActionType

	// Can always fold
	validActions = append(validActions, ActionFold)

	// Check if can check
	if game.CurrentBet == player.GetBet() {
		validActions = append(validActions, ActionCheck)
	}

	// Check if can call
	callAmount := game.CurrentBet - player.GetBet()
	if callAmount > 0 && callAmount <= player.GetChips() {
		validActions = append(validActions, ActionCall)
	}

	// Check if can raise
	if callAmount < player.GetChips() {
		validActions = append(validActions, ActionRaise)
	}

	// Can always go all-in if have chips
	if player.GetChips() > 0 {
		validActions = append(validActions, ActionAllIn)
	}

	return validActions
}
