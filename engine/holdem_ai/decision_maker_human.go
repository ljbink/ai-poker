package holdem_ai

import (
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
)

type HumanDecisionMaker struct {
	validator     holdem.IActionValidator // Action validator for legal moves
	actionChannel chan holdem.Action      // Channel to receive actions from external frontend
}

func NewHumanDecisionMaker() *HumanDecisionMaker {
	return &HumanDecisionMaker{
		validator:     holdem.NewActionValidator(),
		actionChannel: make(chan holdem.Action, 1),
	}
}

// MakeDecision implements the IDecisionMaker interface
// This will wait for an action to be provided via SetAction method
func (d *HumanDecisionMaker) MakeDecision(game *holdem.Game, player holdem.IPlayer) <-chan holdem.Action {
	ch := make(chan holdem.Action, 1)

	go func() {
		defer close(ch)

		// Wait for external frontend to provide an action
		select {
		case action := <-d.actionChannel:
			// Validate the action before returning
			if err := d.validator.ValidateAction(game, player, action); err != nil {
				// If action is invalid, return a fold action as fallback
				fallbackAction := holdem.Action{
					PlayerID: player.GetID(),
					Type:     holdem.ActionFold,
					Amount:   0,
				}
				ch <- fallbackAction
			} else {
				ch <- action
			}
		case <-time.After(60 * time.Second): // 60 second timeout
			// Timeout - return fold action
			timeoutAction := holdem.Action{
				PlayerID: player.GetID(),
				Type:     holdem.ActionFold,
				Amount:   0,
			}
			ch <- timeoutAction
		}
	}()

	return ch
}

// SetAction allows external frontend to provide the human player's action
func (d *HumanDecisionMaker) SetAction(action holdem.Action) {
	select {
	case d.actionChannel <- action:
		// Action sent successfully
	default:
		// Channel is full or not ready, ignore
	}
}

// GetAvailableActions returns the valid actions for the current game state
// This can be used by external frontend to show available options
func (d *HumanDecisionMaker) GetAvailableActions(game *holdem.Game, player holdem.IPlayer) []holdem.ActionType {
	return d.validator.GetAvailableActions(game, player)
}

// GetMinRaiseAmount returns the minimum raise amount
func (d *HumanDecisionMaker) GetMinRaiseAmount(game *holdem.Game, player holdem.IPlayer) int {
	return d.validator.GetMinRaiseAmount(game, player)
}

// GetMaxRaiseAmount returns the maximum raise amount (all-in)
func (d *HumanDecisionMaker) GetMaxRaiseAmount(game *holdem.Game, player holdem.IPlayer) int {
	return d.validator.GetMaxRaiseAmount(game, player)
}

// ValidateAction validates if an action is legal - useful for frontend validation
func (d *HumanDecisionMaker) ValidateAction(game *holdem.Game, player holdem.IPlayer, action holdem.Action) error {
	if err := d.validator.ValidateAction(game, player, action); err != nil {
		return err
	}
	return nil
}

// GetCallAmount calculates the amount needed to call
func (d *HumanDecisionMaker) GetCallAmount(game *holdem.Game, player holdem.IPlayer) int {
	actions := d.getCurrentPhaseActions(game)
	currentBet := 0

	for _, action := range actions {
		if action.Type == holdem.ActionRaise || action.Type == holdem.ActionCall {
			if action.Amount > currentBet {
				currentBet = action.Amount
			}
		}
	}

	callAmount := currentBet - player.GetBet()
	if callAmount < 0 {
		callAmount = 0
	}

	return callAmount
}

// Helper function to get current phase actions
func (d *HumanDecisionMaker) getCurrentPhaseActions(game *holdem.Game) []holdem.Action {
	userActions := game.GetUserActions()

	switch game.GetCurrentPhase() {
	case holdem.PhasePreflop:
		return userActions.Preflop
	case holdem.PhaseFlop:
		return userActions.Flop
	case holdem.PhaseTurn:
		return userActions.Turn
	case holdem.PhaseRiver:
		return userActions.River
	default:
		return []holdem.Action{}
	}
}
