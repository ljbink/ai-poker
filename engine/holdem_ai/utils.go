package holdem_ai

import "github.com/ljbink/ai-poker/engine/holdem"

// Utility functions for formatting game information for display
// These are helper functions that frontends can use to format data

// FormatGamePhase converts GamePhase to string
func FormatGamePhase(phase holdem.GamePhase) string {
	phaseNames := map[holdem.GamePhase]string{
		holdem.PhasePreflop:  "Preflop",
		holdem.PhaseFlop:     "Flop",
		holdem.PhaseTurn:     "Turn",
		holdem.PhaseRiver:    "River",
		holdem.PhaseShowdown: "Showdown",
	}
	return phaseNames[phase]
}

// FormatPlayerCards formats player's hole cards as string
func FormatPlayerCards(player holdem.IPlayer) string {
	playerCards := ""
	for i, card := range player.GetHandCards() {
		if i > 0 {
			playerCards += " "
		}
		playerCards += card.String()
	}
	return playerCards
}

// FormatCommunityCards formats community cards as string
func FormatCommunityCards(game *holdem.Game) string {
	communityCards := ""
	for i, card := range game.CommunityCards {
		if i > 0 {
			communityCards += " "
		}
		communityCards += card.String()
	}
	return communityCards
}

// CalculateCallAmount calculates how much the player needs to call
func CalculateCallAmount(game *holdem.Game, player holdem.IPlayer) int {
	callAmount := game.CurrentBet - player.GetBet()
	if callAmount < 0 {
		return 0
	}
	return callAmount
}

// CalculateMinRaise calculates minimum raise amount
func CalculateMinRaise(game *holdem.Game) int {
	return game.BigBlind
}

// CalculateMaxRaise calculates maximum raise amount (all-in)
func CalculateMaxRaise(game *holdem.Game, player holdem.IPlayer) int {
	callAmount := CalculateCallAmount(game, player)
	maxRaise := player.GetChips() - callAmount
	if maxRaise < 0 {
		return 0
	}
	return maxRaise
}

// Convenience functions for creating common actions

// CreateFoldAction creates a fold action
func CreateFoldAction() Action {
	return Action{Type: ActionFold}
}

// CreateCheckAction creates a check action
func CreateCheckAction() Action {
	return Action{Type: ActionCheck}
}

// CreateCallAction creates a call action
func CreateCallAction() Action {
	return Action{Type: ActionCall}
}

// CreateRaiseAction creates a raise action with specified amount
func CreateRaiseAction(amount int) Action {
	return Action{Type: ActionRaise, Amount: amount}
}

// CreateAllInAction creates an all-in action
func CreateAllInAction() Action {
	return Action{Type: ActionAllIn}
}
