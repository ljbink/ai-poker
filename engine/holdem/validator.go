package holdem

import (
	"fmt"
)

// ValidationError represents an action validation error
type ValidationError struct {
	Message string
	Code    ValidationErrorCode
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ValidationErrorCode represents different types of validation errors
type ValidationErrorCode int

const (
	ErrorInvalidPlayer ValidationErrorCode = iota
	ErrorInvalidAction
	ErrorInsufficientChips
	ErrorInvalidAmount
	ErrorOutOfTurn
	ErrorGameState
	ErrorActionNotAllowed
)

type IActionValidator interface {
	ValidateAction(game *Game, player IPlayer, action Action) *ValidationError
	GetAvailableActions(game *Game, player IPlayer) []ActionType
	GetMinRaiseAmount(game *Game, player IPlayer) int
	GetMaxRaiseAmount(game *Game, player IPlayer) int
}

// ActionValidator provides methods for validating poker actions
type ActionValidator struct{}

// NewActionValidator creates a new action validator
func NewActionValidator() *ActionValidator {
	return &ActionValidator{}
}

// ValidateAction validates if an action is legal in the current game state
func (v *ActionValidator) ValidateAction(game *Game, player IPlayer, action Action) *ValidationError {
	// Basic validations
	if err := v.validateBasicAction(action); err != nil {
		return err
	}

	if err := v.validatePlayer(game, player, action.PlayerID); err != nil {
		return err
	}

	if err := v.validateGameState(game); err != nil {
		return err
	}

	if err := v.validatePlayerTurn(game, player); err != nil {
		return err
	}

	// Action-specific validations
	switch action.Type {
	case ActionFold:
		return v.validateFold(game, player, action)
	case ActionCheck:
		return v.validateCheck(game, player, action)
	case ActionCall:
		return v.validateCall(game, player, action)
	case ActionRaise:
		return v.validateRaise(game, player, action)
	case ActionAllIn:
		return v.validateAllIn(game, player, action)
	default:
		return &ValidationError{
			Message: fmt.Sprintf("Unknown action type: %d", action.Type),
			Code:    ErrorInvalidAction,
		}
	}
}

// GetAvailableActions returns all valid actions for a player in current game state
func (v *ActionValidator) GetAvailableActions(game *Game, player IPlayer) []ActionType {
	var actions []ActionType

	// Basic validations
	if game == nil || player == nil {
		return actions
	}

	if player.IsFolded() {
		return actions // No actions available for folded players
	}

	currentBet := v.getCurrentBet(game)
	playerBet := player.GetBet()
	callAmount := currentBet - playerBet

	// Ensure call amount is not negative
	if callAmount < 0 {
		callAmount = 0
	}

	// Always can fold (unless already folded)
	actions = append(actions, ActionFold)

	// Check if player can check
	if callAmount == 0 {
		actions = append(actions, ActionCheck)
	}

	// Check if player can call
	if callAmount > 0 && player.GetChips() >= callAmount {
		actions = append(actions, ActionCall)
	}

	// Check if player can raise
	if v.canPlayerRaise(game, player) {
		actions = append(actions, ActionRaise)
	}

	// Check if player can go all-in
	if player.GetChips() > 0 {
		actions = append(actions, ActionAllIn)
	}

	return actions
}

// GetMinRaiseAmount returns the minimum raise amount for a player
func (v *ActionValidator) GetMinRaiseAmount(game *Game, player IPlayer) int {
	if game == nil || player == nil {
		return 0
	}

	currentBet := v.getCurrentBet(game)
	playerBet := player.GetBet()
	callAmount := currentBet - playerBet

	if callAmount < 0 {
		callAmount = 0
	}

	// Minimum raise is typically the big blind
	minRaise := game.GetBigBlind()

	// If there's already a bet, minimum raise is the difference between current bet and previous bet
	if currentBet > 0 {
		// Find the previous bet amount to calculate minimum raise
		actions := v.getCurrentPhaseActions(game)
		prevBet := 0
		for _, action := range actions {
			if action.Type == ActionRaise && action.Amount < currentBet {
				prevBet = action.Amount
			}
		}
		if prevBet > 0 {
			minRaise = currentBet - prevBet
		}
	}

	return callAmount + minRaise
}

// GetMaxRaiseAmount returns the maximum raise amount for a player (all-in)
func (v *ActionValidator) GetMaxRaiseAmount(game *Game, player IPlayer) int {
	if game == nil || player == nil {
		return 0
	}

	return player.GetChips()
}

// Basic validation functions
func (v *ActionValidator) validateBasicAction(action Action) *ValidationError {
	if action.PlayerID <= 0 {
		return &ValidationError{
			Message: "Invalid player ID",
			Code:    ErrorInvalidPlayer,
		}
	}

	if action.Amount < 0 {
		return &ValidationError{
			Message: "Action amount cannot be negative",
			Code:    ErrorInvalidAmount,
		}
	}

	return nil
}

func (v *ActionValidator) validatePlayer(game *Game, player IPlayer, actionPlayerID int) *ValidationError {
	if player == nil {
		return &ValidationError{
			Message: "Player is nil",
			Code:    ErrorInvalidPlayer,
		}
	}

	if player.GetID() != actionPlayerID {
		return &ValidationError{
			Message: "Player ID mismatch",
			Code:    ErrorInvalidPlayer,
		}
	}

	if player.IsFolded() {
		return &ValidationError{
			Message: "Player has already folded",
			Code:    ErrorActionNotAllowed,
		}
	}

	return nil
}

func (v *ActionValidator) validateGameState(game *Game) *ValidationError {
	if game == nil {
		return &ValidationError{
			Message: "Game is nil",
			Code:    ErrorGameState,
		}
	}

	// Check if game is in a valid phase for actions
	phase := game.GetCurrentPhase()
	if phase == PhaseShowdown {
		return &ValidationError{
			Message: "No actions allowed during showdown",
			Code:    ErrorGameState,
		}
	}

	return nil
}

func (v *ActionValidator) validatePlayerTurn(game *Game, player IPlayer) *ValidationError {
	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer == nil {
		return &ValidationError{
			Message: "No current player",
			Code:    ErrorOutOfTurn,
		}
	}

	if currentPlayer.GetID() != player.GetID() {
		return &ValidationError{
			Message: "Not player's turn",
			Code:    ErrorOutOfTurn,
		}
	}

	return nil
}

// Action-specific validation functions
func (v *ActionValidator) validateFold(game *Game, player IPlayer, action Action) *ValidationError {
	if action.Amount != 0 {
		return &ValidationError{
			Message: "Fold action should have amount 0",
			Code:    ErrorInvalidAmount,
		}
	}

	return nil
}

func (v *ActionValidator) validateCheck(game *Game, player IPlayer, action Action) *ValidationError {
	if action.Amount != 0 {
		return &ValidationError{
			Message: "Check action should have amount 0",
			Code:    ErrorInvalidAmount,
		}
	}

	currentBet := v.getCurrentBet(game)
	playerBet := player.GetBet()

	if currentBet > playerBet {
		return &ValidationError{
			Message: "Cannot check when there is a bet to call",
			Code:    ErrorActionNotAllowed,
		}
	}

	return nil
}

func (v *ActionValidator) validateCall(game *Game, player IPlayer, action Action) *ValidationError {
	currentBet := v.getCurrentBet(game)
	playerBet := player.GetBet()
	callAmount := currentBet - playerBet

	if callAmount <= 0 {
		return &ValidationError{
			Message: "No bet to call",
			Code:    ErrorActionNotAllowed,
		}
	}

	if action.Amount != callAmount {
		return &ValidationError{
			Message: fmt.Sprintf("Call amount should be %d, got %d", callAmount, action.Amount),
			Code:    ErrorInvalidAmount,
		}
	}

	if player.GetChips() < callAmount {
		return &ValidationError{
			Message: "Insufficient chips to call",
			Code:    ErrorInsufficientChips,
		}
	}

	return nil
}

func (v *ActionValidator) validateRaise(game *Game, player IPlayer, action Action) *ValidationError {
	if action.Amount <= 0 {
		return &ValidationError{
			Message: "Raise amount must be positive",
			Code:    ErrorInvalidAmount,
		}
	}

	currentBet := v.getCurrentBet(game)
	playerBet := player.GetBet()
	callAmount := currentBet - playerBet

	if callAmount < 0 {
		callAmount = 0
	}

	totalRequired := callAmount + action.Amount

	if player.GetChips() < totalRequired {
		return &ValidationError{
			Message: "Insufficient chips to raise",
			Code:    ErrorInsufficientChips,
		}
	}

	minRaise := v.GetMinRaiseAmount(game, player)
	if totalRequired < minRaise {
		return &ValidationError{
			Message: fmt.Sprintf("Raise amount too small. Minimum: %d, got: %d", minRaise, totalRequired),
			Code:    ErrorInvalidAmount,
		}
	}

	return nil
}

func (v *ActionValidator) validateAllIn(game *Game, player IPlayer, action Action) *ValidationError {
	if action.Amount != player.GetChips() {
		return &ValidationError{
			Message: fmt.Sprintf("All-in amount should be %d (all chips), got %d", player.GetChips(), action.Amount),
			Code:    ErrorInvalidAmount,
		}
	}

	if player.GetChips() <= 0 {
		return &ValidationError{
			Message: "Player has no chips to go all-in",
			Code:    ErrorInsufficientChips,
		}
	}

	return nil
}

// Helper functions
func (v *ActionValidator) getCurrentBet(game *Game) int {
	actions := v.getCurrentPhaseActions(game)
	maxBet := 0

	for _, action := range actions {
		if action.Type == ActionRaise || action.Type == ActionCall {
			if action.Amount > maxBet {
				maxBet = action.Amount
			}
		}
	}

	return maxBet
}

func (v *ActionValidator) getCurrentPhaseActions(game *Game) []Action {
	userActions := game.GetUserActions()

	switch game.GetCurrentPhase() {
	case PhasePreflop:
		return userActions.Preflop
	case PhaseFlop:
		return userActions.Flop
	case PhaseTurn:
		return userActions.Turn
	case PhaseRiver:
		return userActions.River
	default:
		return []Action{}
	}
}

func (v *ActionValidator) canPlayerRaise(game *Game, player IPlayer) bool {
	currentBet := v.getCurrentBet(game)
	playerBet := player.GetBet()
	callAmount := currentBet - playerBet

	if callAmount < 0 {
		callAmount = 0
	}

	minRaise := v.GetMinRaiseAmount(game, player)
	return player.GetChips() >= minRaise
}

// Utility functions for external use

// IsValidActionType checks if an action type is valid
func IsValidActionType(actionType ActionType) bool {
	switch actionType {
	case ActionFold, ActionCheck, ActionCall, ActionRaise, ActionAllIn:
		return true
	case ActionSystemShuffle, ActionSystemDealHole, ActionSystemDealFlop, ActionSystemDealTurn, ActionSystemDealRiver, ActionSystemPhaseChange:
		return true
	default:
		return false
	}
}

// ActionTypeToString converts action type to string
func ActionTypeToString(actionType ActionType) string {
	switch actionType {
	case ActionFold:
		return "Fold"
	case ActionCheck:
		return "Check"
	case ActionCall:
		return "Call"
	case ActionRaise:
		return "Raise"
	case ActionAllIn:
		return "All-In"
	case ActionSystemShuffle:
		return "System: Shuffle"
	case ActionSystemDealHole:
		return "System: Deal Hole Cards"
	case ActionSystemDealFlop:
		return "System: Deal Flop"
	case ActionSystemDealTurn:
		return "System: Deal Turn"
	case ActionSystemDealRiver:
		return "System: Deal River"
	case ActionSystemPhaseChange:
		return "System: Phase Change"
	default:
		return "Unknown"
	}
}

// ValidationErrorCodeToString converts validation error code to string
func ValidationErrorCodeToString(code ValidationErrorCode) string {
	switch code {
	case ErrorInvalidPlayer:
		return "Invalid Player"
	case ErrorInvalidAction:
		return "Invalid Action"
	case ErrorInsufficientChips:
		return "Insufficient Chips"
	case ErrorInvalidAmount:
		return "Invalid Amount"
	case ErrorOutOfTurn:
		return "Out of Turn"
	case ErrorGameState:
		return "Game State"
	case ErrorActionNotAllowed:
		return "Action Not Allowed"
	default:
		return "Unknown"
	}
}
