package holdem_ai

import (
	"github.com/ljbink/ai-poker/engine/holdem"
)

// GameAdapter handles all game actions and integrates AI with the game engine
type GameAdapter struct{}

// NewGameAdapter creates a new game adapter
func NewGameAdapter() *GameAdapter {
	return &GameAdapter{}
}

// ExecuteAction translates an AI action into game state changes
func (ga *GameAdapter) ExecuteAction(game *holdem.Game, action Action) bool {
	switch action.Type {
	case ActionFold:
		return ga.executeFold(game)
	case ActionCheck:
		return ga.executeCheck(game)
	case ActionCall:
		return ga.executeCall(game)
	case ActionRaise:
		return ga.executeRaise(game, action.Amount)
	case ActionAllIn:
		// All-in is essentially raising all remaining chips
		player := game.GetCurrentPlayer()
		allInAmount := player.GetChips()
		return ga.executeRaise(game, allInAmount)
	default:
		// Default to fold if unknown action
		return ga.executeFold(game)
	}
}

// executeFold handles fold action
func (ga *GameAdapter) executeFold(game *holdem.Game) bool {
	player := game.GetCurrentPlayer()
	player.Fold()

	game.ActionsThisRound++
	ga.nextPlayer(game)

	// Check if only one player remains
	activePlayers := game.GetActivePlayers()
	if len(activePlayers) <= 1 {
		ga.endHandEarly(game)
		return true
	}

	return ga.checkAndAdvanceIfBettingComplete(game)
}

// executeCall handles call action
func (ga *GameAdapter) executeCall(game *holdem.Game) bool {
	player := game.GetCurrentPlayer()
	callAmount := game.CurrentBet - player.GetBet()

	if callAmount > 0 {
		actualAmount := callAmount
		if callAmount > player.GetChips() {
			actualAmount = player.GetChips()
		}
		player.Bet(actualAmount)
		game.Pot += actualAmount
	}

	game.ActionsThisRound++
	ga.nextPlayer(game)
	return ga.checkAndAdvanceIfBettingComplete(game)
}

// executeRaise handles raise action
func (ga *GameAdapter) executeRaise(game *holdem.Game, raiseAmount int) bool {
	player := game.GetCurrentPlayer()
	totalBet := game.CurrentBet + raiseAmount
	betNeeded := totalBet - player.GetBet()

	actualAmount := betNeeded
	if betNeeded > player.GetChips() {
		actualAmount = player.GetChips()
	}
	player.Bet(actualAmount)
	game.Pot += actualAmount

	if player.GetBet() > game.CurrentBet {
		game.CurrentBet = player.GetBet()
	}

	game.ActionsThisRound++
	ga.nextPlayer(game)
	return ga.checkAndAdvanceIfBettingComplete(game)
}

// executeCheck handles check action
func (ga *GameAdapter) executeCheck(game *holdem.Game) bool {
	game.ActionsThisRound++
	ga.nextPlayer(game)
	return ga.checkAndAdvanceIfBettingComplete(game)
}

// Helper methods moved from Game struct

func (ga *GameAdapter) nextPlayer(game *holdem.Game) {
	original := game.CurrentPlayerPosition
	for {
		game.CurrentPlayerPosition = (game.CurrentPlayerPosition + 1) % len(game.Players)
		player := game.Players[game.CurrentPlayerPosition]

		// Found active player or completed full circle
		if !player.IsFolded() || game.CurrentPlayerPosition == original {
			break
		}
	}
}

func (ga *GameAdapter) checkAndAdvanceIfBettingComplete(game *holdem.Game) bool {
	if ga.isBettingRoundComplete(game) {
		ga.advanceToNextPhase(game)
		return true // Betting round was completed
	}
	return false
}

func (ga *GameAdapter) isBettingRoundComplete(game *holdem.Game) bool {
	activePlayers := game.GetActivePlayers()
	if len(activePlayers) <= 1 {
		return true
	}

	// In preflop, we need special handling because of blinds
	if game.CurrentPhase == holdem.PhasePreflop {
		// Everyone needs to match the current bet or be all-in
		// AND we need at least as many actions as players
		if game.ActionsThisRound < len(activePlayers) {
			return false
		}

		for _, player := range activePlayers {
			if player.GetBet() != game.CurrentBet && player.GetChips() > 0 {
				return false
			}
		}
		return true
	}

	// Post-flop: everyone needs to act at least once
	// and all have same bet (usually 0) or are all-in
	if game.ActionsThisRound < len(activePlayers) {
		return false
	}

	for _, player := range activePlayers {
		if player.GetBet() != game.CurrentBet && player.GetChips() > 0 {
			return false
		}
	}

	return true
}

func (ga *GameAdapter) advanceToNextPhase(game *holdem.Game) bool {
	// Reset bets and action count for next round
	for _, player := range game.Players {
		player.ResetBet()
	}
	game.CurrentBet = 0
	game.ActionsThisRound = 0

	switch game.CurrentPhase {
	case holdem.PhasePreflop:
		game.CurrentPhase = holdem.PhaseFlop
		ga.dealFlop(game)
		ga.setFirstPlayerPostFlop(game)
		return false
	case holdem.PhaseFlop:
		game.CurrentPhase = holdem.PhaseTurn
		ga.dealTurn(game)
		ga.setFirstPlayerPostFlop(game)
		return false
	case holdem.PhaseTurn:
		game.CurrentPhase = holdem.PhaseRiver
		ga.dealRiver(game)
		ga.setFirstPlayerPostFlop(game)
		return false
	case holdem.PhaseRiver:
		game.CurrentPhase = holdem.PhaseShowdown
		ga.showdown(game)
		return true // Hand is over
	case holdem.PhaseShowdown:
		return true // Already at showdown
	}
	return false
}

func (ga *GameAdapter) endHandEarly(game *holdem.Game) {
	activePlayers := game.GetActivePlayers()
	if len(activePlayers) == 1 {
		winner := activePlayers[0]
		winner.GrandChips(game.Pot)
		game.CurrentPhase = holdem.PhaseShowdown
	}
}

func (ga *GameAdapter) showdown(game *holdem.Game) {
	activePlayers := game.GetActivePlayers()
	if len(activePlayers) == 0 {
		return
	}

	// Evaluate all active hands
	type playerResult struct {
		player holdem.IPlayer
		hand   *holdem.HandResult
	}

	var results []playerResult
	for _, player := range activePlayers {
		hand := holdem.EvaluatePlayerHand(player, game.CommunityCards)
		results = append(results, playerResult{player, hand})
	}

	// Find the best hand value
	bestValue := int64(0)
	for _, result := range results {
		if result.hand.Value > bestValue {
			bestValue = result.hand.Value
		}
	}

	// Find all winners (in case of tie)
	var winners []holdem.IPlayer
	for _, result := range results {
		if result.hand.Value == bestValue {
			winners = append(winners, result.player)
		}
	}

	// Distribute pot among winners
	potShare := game.Pot / len(winners)
	remainder := game.Pot % len(winners)

	for i, winner := range winners {
		share := potShare
		if i == 0 {
			share += remainder // Give remainder to first winner
		}
		winner.GrandChips(share)
	}

	game.CurrentPhase = holdem.PhaseShowdown
}

func (ga *GameAdapter) setFirstPlayerPostFlop(game *holdem.Game) {
	game.CurrentPlayerPosition = (game.DealerPosition + 1) % len(game.Players)

	// Find first active player
	original := game.CurrentPlayerPosition
	for game.Players[game.CurrentPlayerPosition].IsFolded() {
		game.CurrentPlayerPosition = (game.CurrentPlayerPosition + 1) % len(game.Players)
		if game.CurrentPlayerPosition == original {
			break
		}
	}
}

func (ga *GameAdapter) dealFlop(game *holdem.Game) {
	// Burn one card
	if len(game.Deck) > 0 {
		game.Deck = game.Deck[1:]
	}

	// Deal 3 cards
	for i := 0; i < 3 && len(game.Deck) > 0; i++ {
		game.CommunityCards.Append(game.Deck[0])
		game.Deck = game.Deck[1:]
	}
}

func (ga *GameAdapter) dealTurn(game *holdem.Game) {
	// Burn one card
	if len(game.Deck) > 0 {
		game.Deck = game.Deck[1:]
	}

	// Deal 1 card
	if len(game.Deck) > 0 {
		game.CommunityCards.Append(game.Deck[0])
		game.Deck = game.Deck[1:]
	}
}

func (ga *GameAdapter) dealRiver(game *holdem.Game) {
	// Burn one card
	if len(game.Deck) > 0 {
		game.Deck = game.Deck[1:]
	}

	// Deal 1 card
	if len(game.Deck) > 0 {
		game.CommunityCards.Append(game.Deck[0])
		game.Deck = game.Deck[1:]
	}
}
