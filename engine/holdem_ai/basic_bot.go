package holdem_ai

import (
	"math/rand"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
	"github.com/ljbink/ai-poker/engine/poker"
)

// BasicBot implements a simple poker bot with basic strategy
type BasicBot struct {
	name      string
	validator *ActionValidator
	rng       *rand.Rand
}

// NewBasicBot creates a new basic bot with the given name
func NewBasicBot(name string) DecisionMaker {
	return &BasicBot{
		name:      name,
		validator: NewActionValidator(),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetName returns the bot's name
func (b *BasicBot) GetName() string {
	return b.name
}

// IsBot returns true since this is a bot
func (b *BasicBot) IsBot() bool {
	return true
}

// MakeDecision implements the main decision-making logic
func (b *BasicBot) MakeDecision(game *holdem.Game, player holdem.IPlayer) Action {
	// Get valid actions
	validActions := b.validator.GetValidActions(game, player)

	// Calculate hand strength
	handStrength := b.evaluateHandStrength(game, player)

	// Make decision based on hand strength and game situation
	return b.chooseAction(game, player, validActions, handStrength)
}

// evaluateHandStrength evaluates the strength of the current hand
func (b *BasicBot) evaluateHandStrength(game *holdem.Game, player holdem.IPlayer) float64 {
	playerCards := player.GetHandCards()
	if len(playerCards) != 2 {
		return 0.0 // No cards yet
	}

	// Basic hand strength evaluation
	strength := 0.0

	// Check for pocket pairs
	if playerCards[0].Rank == playerCards[1].Rank {
		pairRank := playerCards[0].Rank
		switch {
		case pairRank == poker.RankAce: // Pocket Aces
			strength = 0.95
		case pairRank >= poker.RankJack: // JJ, QQ, KK
			strength = 0.9
		case pairRank >= poker.RankEight: // 88, 99, TT
			strength = 0.7
		default: // Low pairs
			strength = 0.5
		}
	} else {
		// Non-pair hands - convert ranks for comparison (Ace high)
		rank1 := b.convertRankForComparison(playerCards[0].Rank)
		rank2 := b.convertRankForComparison(playerCards[1].Rank)

		highRank := rank1
		lowRank := rank2
		if rank2 > rank1 {
			highRank = rank2
			lowRank = rank1
		}

		// Premium hands
		if (highRank == 14 && lowRank >= 13) || // AK, AQ
			(highRank == 13 && lowRank == 12) { // KQ
			strength = 0.8
		} else if highRank >= 11 { // Any hand with J, Q, K, A
			strength = 0.6
		} else if highRank >= 8 { // Medium cards
			strength = 0.4
		} else {
			strength = 0.2
		}

		// Suited bonus
		if playerCards[0].Suit == playerCards[1].Suit {
			strength += 0.1
		}

		// Connected bonus (for straights)
		if abs(highRank-lowRank) == 1 {
			strength += 0.05
		}
	}

	// Adjust for position
	strength = b.adjustForPosition(strength, game, player)

	// Adjust for community cards if available
	if len(game.CommunityCards) > 0 {
		strength = b.adjustForCommunityCards(strength, game)
	}

	return min(strength, 1.0)
}

// convertRankForComparison converts poker ranks to numerical values for comparison
// Treats Ace as high (14), King as 13, Queen as 12, Jack as 11, etc.
func (b *BasicBot) convertRankForComparison(rank poker.Rank) int {
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

// adjustForPosition adjusts hand strength based on position
func (b *BasicBot) adjustForPosition(strength float64, game *holdem.Game, player holdem.IPlayer) float64 {
	totalPlayers := game.GetTotalPlayers()
	position := game.GetPlayerPosition(player)

	// Early position (first 1/3 of players) - be more conservative
	if position < totalPlayers/3 {
		return strength * 0.9
	}

	// Late position (last 1/3 of players) - be more aggressive
	if position >= (totalPlayers*2)/3 {
		return strength * 1.1
	}

	// Middle position - no adjustment
	return strength
}

// adjustForCommunityCards adjusts hand strength based on community cards
func (b *BasicBot) adjustForCommunityCards(strength float64, game *holdem.Game) float64 {
	// This is a simplified implementation
	// In a real bot, you'd want to evaluate actual hand rankings

	// For now, just slightly reduce confidence with more community cards
	// as more information makes the hand more defined
	switch len(game.CommunityCards) {
	case 3:
		return strength * 0.95
	case 4:
		return strength * 0.9
	case 5:
		return strength * 0.85
	default:
		return strength
	}
}

// chooseAction selects the best action based on hand strength and game state
func (b *BasicBot) chooseAction(game *holdem.Game, player holdem.IPlayer, validActions []ActionType, handStrength float64) Action {
	// Calculate pot odds
	potOdds := game.CalculatePotOdds(player)

	// Decision thresholds
	foldThreshold := 0.3
	raiseThreshold := 0.7

	// Adjust thresholds based on pot odds
	if potOdds > 0.3 { // Good pot odds
		foldThreshold -= 0.1
	}
	if potOdds > 0.5 { // Excellent pot odds
		foldThreshold -= 0.1
	}

	// Make decision
	if handStrength < foldThreshold {
		return Action{Type: ActionFold}
	}

	if handStrength > raiseThreshold && b.canRaise(validActions) {
		// Decide raise amount
		raiseAmount := b.calculateRaiseAmount(game, player, handStrength)
		return Action{Type: ActionRaise, Amount: raiseAmount}
	}

	// Medium strength hands - call or check
	if b.canCall(validActions) {
		return Action{Type: ActionCall}
	}

	if b.canCheck(validActions) {
		return Action{Type: ActionCheck}
	}

	// Last resort - fold
	return Action{Type: ActionFold}
}

// calculateRaiseAmount determines how much to raise
func (b *BasicBot) calculateRaiseAmount(game *holdem.Game, player holdem.IPlayer, handStrength float64) int {
	// Base raise amount
	baseRaise := game.BigBlind * 3

	// Adjust based on hand strength
	multiplier := 1.0 + (handStrength-0.7)*2 // Range: 1.0 to 1.6

	// Adjust based on position (more aggressive in late position)
	totalPlayers := game.GetTotalPlayers()
	position := game.GetPlayerPosition(player)
	if position >= (totalPlayers*2)/3 {
		multiplier *= 1.2
	}

	// Add some randomness
	multiplier *= 0.8 + b.rng.Float64()*0.4 // ±20% randomness

	raiseAmount := int(float64(baseRaise) * multiplier)

	// Don't raise more than 1/4 of our stack
	maxRaise := player.GetChips() / 4
	if raiseAmount > maxRaise {
		raiseAmount = maxRaise
	}

	// Minimum raise should be at least the big blind
	if raiseAmount < game.BigBlind {
		raiseAmount = game.BigBlind
	}

	return raiseAmount
}

// Helper methods to check if specific actions are available
func (b *BasicBot) canRaise(validActions []ActionType) bool {
	for _, action := range validActions {
		if action == ActionRaise {
			return true
		}
	}
	return false
}

func (b *BasicBot) canCall(validActions []ActionType) bool {
	for _, action := range validActions {
		if action == ActionCall {
			return true
		}
	}
	return false
}

func (b *BasicBot) canCheck(validActions []ActionType) bool {
	for _, action := range validActions {
		if action == ActionCheck {
			return true
		}
	}
	return false
}

// Utility functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
