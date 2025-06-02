package holdem_ai

import (
	"math/rand"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
	"github.com/ljbink/ai-poker/engine/poker"
)

type BasicBotDecisionMaker struct {
	Aggressiveness float64                 // 0.0 = very conservative, 1.0 = very aggressive
	BluffFrequency float64                 // 0.0 = never bluff, 1.0 = always bluff
	evaluator      holdem.IHandEvaluator   // Hand evaluator for strength calculation
	validator      holdem.IActionValidator // Action validator for legal moves
}

// NewBasicBotDecisionMaker creates a new basic bot with specified traits
func NewBasicBotDecisionMaker(aggressiveness, bluffFrequency float64) *BasicBotDecisionMaker {
	return &BasicBotDecisionMaker{
		Aggressiveness: aggressiveness,
		BluffFrequency: bluffFrequency,
		evaluator:      holdem.NewHandEvaluator(),
		validator:      holdem.NewActionValidator(),
	}
}

// MakeDecision implements the IDecisionMaker interface
func (d *BasicBotDecisionMaker) MakeDecision(game *holdem.Game, player holdem.IPlayer) <-chan holdem.Action {
	ch := make(chan holdem.Action, 1)

	go func() {
		defer close(ch)

		// Add realistic thinking time
		thinkingTime := time.Duration(500+rand.Intn(1500)) * time.Millisecond
		time.Sleep(thinkingTime)

		action := d.calculateBestAction(game, player)
		ch <- action
	}()

	return ch
}

// calculateBestAction determines the best action based on hand strength, game state, and bot personality
func (d *BasicBotDecisionMaker) calculateBestAction(game *holdem.Game, player holdem.IPlayer) holdem.Action {
	// Handle nil inputs gracefully
	if game == nil || player == nil {
		return holdem.Action{
			PlayerID: 0,
			Type:     holdem.ActionFold,
			Amount:   0,
		}
	}

	// Get available actions from validator
	availableActions := d.validator.GetAvailableActions(game, player)

	// If no actions available, fold by default
	if len(availableActions) == 0 {
		return holdem.Action{
			PlayerID: player.GetID(),
			Type:     holdem.ActionFold,
			Amount:   0,
		}
	}

	// Check if player is properly seated - if not, fold
	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer == nil || currentPlayer.GetID() != player.GetID() {
		// Player is not the current player or not properly seated, should fold
		if d.isActionAvailable(holdem.ActionFold, availableActions) {
			return holdem.Action{
				PlayerID: player.GetID(),
				Type:     holdem.ActionFold,
				Amount:   0,
			}
		}
	}

	// Evaluate hand strength
	handStrength := d.evaluateHandStrength(game, player)

	// Get betting information
	minRaise := d.validator.GetMinRaiseAmount(game, player)
	maxRaise := d.validator.GetMaxRaiseAmount(game, player)

	// Make decision based on hand strength and available actions
	return d.makeDecisionBasedOnStrength(game, player, handStrength, availableActions, minRaise, maxRaise)
}

// evaluateHandStrength calculates the strength of the current hand (0.0 to 1.0)
func (d *BasicBotDecisionMaker) evaluateHandStrength(game *holdem.Game, player holdem.IPlayer) float64 {
	holeCards := player.GetHandCards()
	communityCards := game.GetCommunityCards()

	if len(holeCards) < 2 {
		return 0.0
	}

	// Evaluate current hand
	handResult := d.evaluator.EvaluateHand(holeCards, communityCards)

	// Convert hand rank to strength percentage
	baseStrength := d.handRankToStrength(handResult.Rank)

	// Adjust for community cards and position
	adjustedStrength := d.adjustStrengthForGameState(baseStrength, game, player, handResult)

	return minFloat64(adjustedStrength, 1.0)
}

// handRankToStrength converts hand rank to base strength value
func (d *BasicBotDecisionMaker) handRankToStrength(rank holdem.HandRank) float64 {
	switch rank {
	case holdem.RoyalFlush:
		return 1.0
	case holdem.StraightFlush:
		return 0.95
	case holdem.FourOfAKind:
		return 0.9
	case holdem.FullHouse:
		return 0.85
	case holdem.Flush:
		return 0.75
	case holdem.Straight:
		return 0.65
	case holdem.ThreeOfAKind:
		return 0.55
	case holdem.TwoPair:
		return 0.45
	case holdem.OnePair:
		return 0.3
	case holdem.HighCard:
		return 0.1
	default:
		return 0.0
	}
}

// adjustStrengthForGameState modifies hand strength based on game context
func (d *BasicBotDecisionMaker) adjustStrengthForGameState(baseStrength float64, game *holdem.Game, player holdem.IPlayer, handResult *holdem.HandResult) float64 {
	adjustment := 0.0

	// Phase adjustments
	switch game.GetCurrentPhase() {
	case holdem.PhasePreflop:
		// Pre-flop: focus on hole card quality
		adjustment += d.evaluatePreflop(player.GetHandCards())
	case holdem.PhaseFlop, holdem.PhaseTurn, holdem.PhaseRiver:
		// Post-flop: consider draws and hand development
		adjustment += d.evaluatePostFlop(handResult, game.GetCommunityCards())
	}

	// Position adjustment (simple implementation)
	activePlayers := d.countActivePlayers(game)
	if activePlayers <= 3 {
		adjustment += 0.1 // Bonus for short-handed play
	}

	return baseStrength + adjustment
}

// evaluatePreflop evaluates hole cards for pre-flop strength
func (d *BasicBotDecisionMaker) evaluatePreflop(holeCards []*poker.Card) float64 {
	if len(holeCards) < 2 {
		return 0.0
	}

	card1, card2 := holeCards[0], holeCards[1]
	rank1 := d.rankToValue(card1.Rank)
	rank2 := d.rankToValue(card2.Rank)

	// Pocket pairs bonus
	if rank1 == rank2 {
		switch {
		case rank1 >= 13: // AA, KK
			return 0.4
		case rank1 >= 10: // QQ, JJ, TT
			return 0.3
		case rank1 >= 7: // 99, 88, 77
			return 0.2
		default:
			return 0.1
		}
	}

	// High cards and suited connectors
	highRank := maxInt(rank1, rank2)
	lowRank := minInt(rank1, rank2)
	suited := card1.Suit == card2.Suit
	connected := abs(highRank-lowRank) == 1

	adjustment := 0.0

	// High card bonus
	if highRank >= 12 { // A, K
		adjustment += 0.15
	} else if highRank >= 10 { // Q, J
		adjustment += 0.1
	}

	// Suited bonus
	if suited {
		adjustment += 0.05
	}

	// Connected bonus
	if connected {
		adjustment += 0.03
	}

	return adjustment
}

// evaluatePostFlop evaluates hand development after the flop
func (d *BasicBotDecisionMaker) evaluatePostFlop(handResult *holdem.HandResult, communityCards poker.Cards) float64 {
	adjustment := 0.0

	// Bonus for made hands vs draws
	if handResult.Rank >= holdem.OnePair {
		adjustment += 0.05
	}

	// Additional bonus for strong made hands
	if handResult.Rank >= holdem.ThreeOfAKind {
		adjustment += 0.1
	}

	return adjustment
}

// makeDecisionBasedOnStrength chooses action based on hand strength and personality
func (d *BasicBotDecisionMaker) makeDecisionBasedOnStrength(game *holdem.Game, player holdem.IPlayer, handStrength float64, availableActions []holdem.ActionType, minRaise, maxRaise int) holdem.Action {
	// Adjust thresholds based on aggressiveness
	foldThreshold := 0.25 - (d.Aggressiveness * 0.1)
	callThreshold := 0.5 - (d.Aggressiveness * 0.15)
	raiseThreshold := 0.7 - (d.Aggressiveness * 0.2)

	// Default action
	action := holdem.Action{
		PlayerID: player.GetID(),
		Type:     holdem.ActionFold,
		Amount:   0,
	}

	// Decision logic
	if handStrength < foldThreshold {
		// Weak hand - fold if possible
		if d.isActionAvailable(holdem.ActionFold, availableActions) {
			action.Type = holdem.ActionFold
		} else if d.isActionAvailable(holdem.ActionCheck, availableActions) {
			action.Type = holdem.ActionCheck
		}
	} else if handStrength < callThreshold {
		// Marginal hand - check/call or bluff
		if d.shouldBluff(handStrength) && d.isActionAvailable(holdem.ActionRaise, availableActions) {
			// Bluff bet
			action.Type = holdem.ActionRaise
			action.Amount = d.calculateBluffAmount(game, player, minRaise)
		} else if d.isActionAvailable(holdem.ActionCall, availableActions) {
			action.Type = holdem.ActionCall
			action.Amount = d.calculateCallAmount(game, player)
		} else if d.isActionAvailable(holdem.ActionCheck, availableActions) {
			action.Type = holdem.ActionCheck
		}
	} else if handStrength < raiseThreshold {
		// Good hand - bet for value or call
		if d.isActionAvailable(holdem.ActionRaise, availableActions) && rand.Float64() < (0.5+d.Aggressiveness*0.3) {
			action.Type = holdem.ActionRaise
			action.Amount = d.calculateValueBetAmount(game, player, handStrength, minRaise)
		} else if d.isActionAvailable(holdem.ActionCall, availableActions) {
			action.Type = holdem.ActionCall
			action.Amount = d.calculateCallAmount(game, player)
		} else if d.isActionAvailable(holdem.ActionCheck, availableActions) {
			action.Type = holdem.ActionCheck
		}
	} else {
		// Strong hand - raise aggressively
		if d.isActionAvailable(holdem.ActionRaise, availableActions) {
			action.Type = holdem.ActionRaise
			action.Amount = d.calculateAggressiveRaiseAmount(game, player, handStrength, minRaise, maxRaise)
		} else if d.isActionAvailable(holdem.ActionCall, availableActions) {
			action.Type = holdem.ActionCall
			action.Amount = d.calculateCallAmount(game, player)
		} else if d.isActionAvailable(holdem.ActionCheck, availableActions) {
			action.Type = holdem.ActionCheck
		}
	}

	// Validate the action before returning
	if err := d.validator.ValidateAction(game, player, action); err != nil {
		// If action is invalid, fallback to check or fold
		if d.isActionAvailable(holdem.ActionCheck, availableActions) {
			action.Type = holdem.ActionCheck
			action.Amount = 0
		} else if d.isActionAvailable(holdem.ActionFold, availableActions) {
			action.Type = holdem.ActionFold
			action.Amount = 0
		}
	}

	return action
}

// Betting amount calculation methods
func (d *BasicBotDecisionMaker) calculateCallAmount(game *holdem.Game, player holdem.IPlayer) int {
	if game == nil || player == nil {
		return 0
	}

	// Get current phase actions
	var actions []holdem.Action
	userActions := game.GetUserActions()
	switch game.GetCurrentPhase() {
	case holdem.PhasePreflop:
		actions = userActions.Preflop
	case holdem.PhaseFlop:
		actions = userActions.Flop
	case holdem.PhaseTurn:
		actions = userActions.Turn
	case holdem.PhaseRiver:
		actions = userActions.River
	default:
		return 0
	}

	// Find highest bet/raise amount in current phase
	currentBet := 0
	for _, action := range actions {
		if action.Type == holdem.ActionRaise || action.Type == holdem.ActionCall {
			if action.Amount > currentBet {
				currentBet = action.Amount
			}
		}
	}

	// Calculate call amount (difference between current bet and player's bet)
	callAmount := currentBet - player.GetBet()
	if callAmount < 0 {
		callAmount = 0
	}

	return callAmount
}

func (d *BasicBotDecisionMaker) calculateBluffAmount(game *holdem.Game, player holdem.IPlayer, minRaise int) int {
	bigBlind := game.GetBigBlind()
	bluffSize := bigBlind + int(float64(bigBlind)*d.Aggressiveness)
	return maxInt(bluffSize, minRaise)
}

func (d *BasicBotDecisionMaker) calculateValueBetAmount(game *holdem.Game, player holdem.IPlayer, handStrength float64, minRaise int) int {
	bigBlind := game.GetBigBlind()
	betSize := int(float64(bigBlind) * (1 + handStrength + d.Aggressiveness) * 2)
	maxBet := player.GetChips() / 3 // Don't bet more than 1/3 of stack

	betAmount := minInt(betSize, maxBet)
	return maxInt(betAmount, minRaise)
}

func (d *BasicBotDecisionMaker) calculateAggressiveRaiseAmount(game *holdem.Game, player holdem.IPlayer, handStrength float64, minRaise, maxRaise int) int {
	bigBlind := game.GetBigBlind()

	// Strong hands warrant bigger bets
	multiplier := 3.0 + (handStrength * 2.0) + (d.Aggressiveness * 2.0)
	raiseAmount := int(float64(bigBlind) * multiplier)

	// Cap at reasonable percentage of stack
	maxBet := player.GetChips() / 2
	raiseAmount = minInt(raiseAmount, maxBet)

	return maxInt(raiseAmount, minRaise)
}

// Helper methods
func (d *BasicBotDecisionMaker) shouldBluff(handStrength float64) bool {
	// Only bluff with marginal hands and based on bluff frequency
	return handStrength > 0.1 && handStrength < 0.4 && rand.Float64() < d.BluffFrequency
}

func (d *BasicBotDecisionMaker) isActionAvailable(actionType holdem.ActionType, availableActions []holdem.ActionType) bool {
	for _, available := range availableActions {
		if available == actionType {
			return true
		}
	}
	return false
}

func (d *BasicBotDecisionMaker) countActivePlayers(game *holdem.Game) int {
	count := 0
	for i := 0; i < 10; i++ {
		if player, err := game.GetPlayerBySit(i); err == nil && player != nil && !player.IsFolded() {
			count++
		}
	}
	return count
}

func (d *BasicBotDecisionMaker) rankToValue(rank poker.Rank) int {
	switch rank {
	case poker.RankAce:
		return 14
	case poker.RankKing:
		return 13
	case poker.RankQueen:
		return 12
	case poker.RankJack:
		return 11
	case poker.RankTen:
		return 10
	case poker.RankNine:
		return 9
	case poker.RankEight:
		return 8
	case poker.RankSeven:
		return 7
	case poker.RankSix:
		return 6
	case poker.RankFive:
		return 5
	case poker.RankFour:
		return 4
	case poker.RankThree:
		return 3
	case poker.RankTwo:
		return 2
	default:
		return 0
	}
}

// Utility functions
func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
