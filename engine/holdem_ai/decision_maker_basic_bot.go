package holdem_ai

import (
	"math/rand"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
	"github.com/ljbink/ai-poker/engine/poker"
)

// BasicBotDecisionMaker implements a simple poker bot with basic strategy
type BasicBotDecisionMaker struct {
	player    holdem.IPlayer
	game      *holdem.Game
	validator *ActionValidator
	rng       *rand.Rand
}

// NewBasicBotDecisionMaker creates a new basic bot bound to a specific player and game
func NewBasicBotDecisionMaker(player holdem.IPlayer, game *holdem.Game) DecisionMaker {
	return &BasicBotDecisionMaker{
		player:    player,
		game:      game,
		validator: NewActionValidator(),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// MakeDecision returns a channel that will receive the chosen action
func (b *BasicBotDecisionMaker) MakeDecision() <-chan Action {
	resultChan := make(chan Action, 1)

	go func() {
		// Simulate thinking time for more realistic behavior
		thinkingTime := time.Duration(100+b.rng.Intn(500)) * time.Millisecond
		time.Sleep(thinkingTime)

		// Get valid actions
		validActions := b.validator.GetValidActions(b.game, b.player)

		// Calculate hand strength
		handStrength := b.evaluateHandStrength()

		// Make decision based on hand strength and game situation
		action := b.chooseAction(validActions, handStrength)

		resultChan <- action
		close(resultChan)
	}()

	return resultChan
}

// evaluateHandStrength evaluates the strength of the current hand
func (b *BasicBotDecisionMaker) evaluateHandStrength() float64 {
	playerCards := b.player.GetHandCards()
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
	strength = b.adjustForPosition(strength)

	// Adjust for community cards if available
	if len(b.game.CommunityCards) > 0 {
		strength = b.adjustForCommunityCards(strength)
	}

	return min(strength, 1.0)
}

// convertRankForComparison converts poker ranks to numerical values for comparison
// Treats Ace as high (14), King as 13, Queen as 12, Jack as 11, etc.
func (b *BasicBotDecisionMaker) convertRankForComparison(rank poker.Rank) int {
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
func (b *BasicBotDecisionMaker) adjustForPosition(strength float64) float64 {
	totalPlayers := b.game.GetTotalPlayers()
	position := b.game.GetPlayerPosition(b.player)

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
func (b *BasicBotDecisionMaker) adjustForCommunityCards(strength float64) float64 {
	// This is a simplified implementation
	// In a real bot, you'd want to evaluate actual hand rankings

	// For now, just slightly reduce confidence with more community cards
	// as more information makes the hand more defined
	switch len(b.game.CommunityCards) {
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
func (b *BasicBotDecisionMaker) chooseAction(validActions []ActionType, handStrength float64) Action {
	// Calculate pot odds
	potOdds := b.game.CalculatePotOdds(b.player)

	// Decision thresholds
	foldThreshold := 0.3
	raiseThreshold := 0.7

	// Adjust thresholds based on pot odds
	if potOdds > 0 {
		// Better pot odds mean we can play looser
		foldThreshold -= potOdds * 0.5
		raiseThreshold -= potOdds * 0.3
	}

	// Decision logic
	if handStrength < foldThreshold {
		// Weak hand - fold unless we can check
		if b.canCheck(validActions) {
			return Action{Type: ActionCheck}
		}
		return Action{Type: ActionFold}
	} else if handStrength > raiseThreshold {
		// Strong hand - raise or bet
		if b.canRaise(validActions) {
			raiseAmount := b.calculateRaiseAmount(handStrength)
			return Action{Type: ActionRaise, Amount: raiseAmount}
		} else if b.canCall(validActions) {
			return Action{Type: ActionCall}
		} else if b.canCheck(validActions) {
			return Action{Type: ActionCheck}
		}
		return Action{Type: ActionFold}
	} else {
		// Medium hand - call or check
		if b.canCheck(validActions) {
			return Action{Type: ActionCheck}
		} else if b.canCall(validActions) {
			// Sometimes fold medium hands if pot odds are bad
			if potOdds < 0.3 && b.rng.Float64() < 0.3 {
				return Action{Type: ActionFold}
			}
			return Action{Type: ActionCall}
		}
		return Action{Type: ActionFold}
	}
}

// calculateRaiseAmount determines how much to raise based on hand strength
func (b *BasicBotDecisionMaker) calculateRaiseAmount(handStrength float64) int {
	// Base raise amount
	baseRaise := b.game.BigBlind * 2

	// Scale based on hand strength
	multiplier := 1.0 + (handStrength-0.5)*2 // Range: 0 to 3

	raiseAmount := int(float64(baseRaise) * multiplier)

	// Ensure we don't raise more than we can afford
	callAmount := b.game.CurrentBet - b.player.GetBet()
	maxRaise := b.player.GetChips() - callAmount

	if raiseAmount > maxRaise {
		raiseAmount = maxRaise
	}

	// Minimum raise
	if raiseAmount < b.game.BigBlind {
		raiseAmount = b.game.BigBlind
	}

	return raiseAmount
}

// Helper methods to check valid actions
func (b *BasicBotDecisionMaker) canRaise(validActions []ActionType) bool {
	for _, action := range validActions {
		if action == ActionRaise {
			return true
		}
	}
	return false
}

func (b *BasicBotDecisionMaker) canCall(validActions []ActionType) bool {
	for _, action := range validActions {
		if action == ActionCall {
			return true
		}
	}
	return false
}

func (b *BasicBotDecisionMaker) canCheck(validActions []ActionType) bool {
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
