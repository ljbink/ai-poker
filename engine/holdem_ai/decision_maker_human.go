package holdem_ai

import (
	"fmt"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
)

// HumanUserDecisionMaker handles decision making for human players
// It receives input from the frontend through channels
type HumanUserDecisionMaker struct {
	player         holdem.IPlayer
	game           *holdem.Game
	actionChan     chan Action
	validator      *ActionValidator
	timeout        time.Duration
	onActionNeeded ActionNeededCallback
}

// ActionNeededCallback is called when the human player needs to make a decision
// The callback receives the game state and current player
type ActionNeededCallback func(game *holdem.Game, player holdem.IPlayer, validActions []ActionType)

// NewHumanUserDecisionMaker creates a new human decision maker bound to a specific player and game
func NewHumanUserDecisionMaker(player holdem.IPlayer, game *holdem.Game) *HumanUserDecisionMaker {
	return &HumanUserDecisionMaker{
		player:     player,
		game:       game,
		actionChan: make(chan Action, 1), // Buffered channel for single action
		validator:  NewActionValidator(),
		timeout:    30 * time.Second, // 30 second timeout for user decisions
	}
}

// MakeDecision returns a channel that will receive the chosen action
func (h *HumanUserDecisionMaker) MakeDecision() <-chan Action {
	resultChan := make(chan Action, 1)

	go func() {
		// Get valid actions for this player
		validActions := h.validator.GetValidActions(h.game, h.player)

		// Notify frontend that action is needed
		if h.onActionNeeded != nil {
			go h.onActionNeeded(h.game, h.player, validActions)
		}

		// Wait for user input with timeout
		select {
		case action := <-h.actionChan:
			// Validate the action before returning
			if h.validator.IsValidAction(action, h.game, h.player) {
				resultChan <- action
			} else {
				// If action is invalid, auto-fold as fallback
				resultChan <- Action{Type: ActionFold}
			}
		case <-time.After(h.timeout):
			// Timeout - auto-fold
			resultChan <- Action{Type: ActionFold}
		}
		close(resultChan)
	}()

	return resultChan
}

// SendAction allows the frontend to send an action (thread-safe)
func (h *HumanUserDecisionMaker) SendAction(action Action) error {
	// Validate action
	if !h.validator.IsValidAction(action, h.game, h.player) {
		validActions := h.validator.GetValidActions(h.game, h.player)
		return fmt.Errorf("invalid action %s, valid actions: %v", action.Type.String(), validActions)
	}

	// Send action (non-blocking)
	select {
	case h.actionChan <- action:
		return nil
	default:
		return fmt.Errorf("action channel is full, previous action not processed")
	}
}

// GetValidActions returns the currently valid actions for the player
func (h *HumanUserDecisionMaker) GetValidActions() []ActionType {
	return h.validator.GetValidActions(h.game, h.player)
}

// SetTimeout sets the decision timeout duration
func (h *HumanUserDecisionMaker) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
}

// SetActionNeededCallback sets the callback for when action is needed
func (h *HumanUserDecisionMaker) SetActionNeededCallback(callback ActionNeededCallback) {
	h.onActionNeeded = callback
}
