package holdem

import (
	"sort"

	"github.com/ljbink/ai-poker/engine/poker"
	"github.com/samber/lo"
)

// HandRank represents the strength of a poker hand
type HandRank int

const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// HandResult represents the evaluation of a hand
type HandResult struct {
	Rank        HandRank
	Description string
	Value       int64       // Numeric value for comparison
	Cards       poker.Cards // The 5 cards that make the hand
}

// EvaluatePlayerHand evaluates the best 5-card hand from player's cards + community cards
func EvaluatePlayerHand(player IPlayer, communityCards poker.Cards) *HandResult {
	playerCards := player.GetHandCards()
	if len(playerCards) < 2 {
		return &HandResult{
			Rank:        HighCard,
			Description: "Insufficient cards",
			Value:       0,
			Cards:       poker.Cards{},
		}
	}

	allCards := poker.Cards{}
	allCards.Append(playerCards...)
	allCards.Append(communityCards...)

	if len(allCards) < 5 {
		return &HandResult{
			Rank:        HighCard,
			Description: "Insufficient cards",
			Value:       0,
			Cards:       allCards,
		}
	}

	// Generate all possible 5-card combinations
	combinations := generateCombinations(allCards, 5)

	handResults := lo.Map(combinations, func(combo poker.Cards, _ int) *HandResult {
		return evaluateFiveCards(combo)
	})

	return lo.MaxBy(handResults, func(a, b *HandResult) bool {
		return a.Value > b.Value
	})
}

// evaluateFiveCards evaluates exactly 5 cards
func evaluateFiveCards(cards poker.Cards) *HandResult {
	if len(cards) != 5 {
		return &HandResult{Value: 0}
	}

	// Sort cards by rank for easier analysis
	sortedCards := make(poker.Cards, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return getRankValue(sortedCards[i].Rank) > getRankValue(sortedCards[j].Rank)
	})

	isFlush := checkFlush(sortedCards)
	isStraight := checkStraight(sortedCards)
	rankCounts := getRankCounts(sortedCards)

	// Check for each hand type from highest to lowest
	if isStraight && isFlush {
		if isRoyalFlush(sortedCards) {
			return &HandResult{
				Rank:        RoyalFlush,
				Description: "Royal Flush",
				Value:       calculateValue(RoyalFlush, sortedCards),
				Cards:       sortedCards,
			}
		}
		return &HandResult{
			Rank:        StraightFlush,
			Description: "Straight Flush",
			Value:       calculateValue(StraightFlush, sortedCards),
			Cards:       sortedCards,
		}
	}

	if hasFourOfAKind(rankCounts) {
		return &HandResult{
			Rank:        FourOfAKind,
			Description: "Four of a Kind",
			Value:       calculateValue(FourOfAKind, sortedCards),
			Cards:       sortedCards,
		}
	}

	if hasFullHouse(rankCounts) {
		return &HandResult{
			Rank:        FullHouse,
			Description: "Full House",
			Value:       calculateValue(FullHouse, sortedCards),
			Cards:       sortedCards,
		}
	}

	if isFlush {
		return &HandResult{
			Rank:        Flush,
			Description: "Flush",
			Value:       calculateValue(Flush, sortedCards),
			Cards:       sortedCards,
		}
	}

	if isStraight {
		return &HandResult{
			Rank:        Straight,
			Description: "Straight",
			Value:       calculateValue(Straight, sortedCards),
			Cards:       sortedCards,
		}
	}

	if hasThreeOfAKind(rankCounts) {
		return &HandResult{
			Rank:        ThreeOfAKind,
			Description: "Three of a Kind",
			Value:       calculateValue(ThreeOfAKind, sortedCards),
			Cards:       sortedCards,
		}
	}

	if hasTwoPair(rankCounts) {
		return &HandResult{
			Rank:        TwoPair,
			Description: "Two Pair",
			Value:       calculateValue(TwoPair, sortedCards),
			Cards:       sortedCards,
		}
	}

	if hasOnePair(rankCounts) {
		return &HandResult{
			Rank:        OnePair,
			Description: "One Pair",
			Value:       calculateValue(OnePair, sortedCards),
			Cards:       sortedCards,
		}
	}

	return &HandResult{
		Rank:        HighCard,
		Description: "High Card",
		Value:       calculateValue(HighCard, sortedCards),
		Cards:       sortedCards,
	}
}

// Helper functions
func getRankValue(rank poker.Rank) int {
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

func checkFlush(cards poker.Cards) bool {
	if len(cards) != 5 {
		return false
	}
	suit := cards[0].Suit
	return lo.EveryBy(cards[1:], func(card *poker.Card) bool {
		return card.Suit == suit
	})
}

func checkStraight(cards poker.Cards) bool {
	if len(cards) != 5 {
		return false
	}

	values := lo.Map(cards, func(card *poker.Card, _ int) int {
		return getRankValue(card.Rank)
	})
	sort.Ints(values)

	// Check for regular straight
	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1]+1 {
			// Check for low ace straight (A-2-3-4-5)
			if values[0] == 2 && values[1] == 3 && values[2] == 4 && values[3] == 5 && values[4] == 14 {
				return true
			}
			return false
		}
	}
	return true
}

func isRoyalFlush(cards poker.Cards) bool {
	if !checkFlush(cards) || !checkStraight(cards) {
		return false
	}
	values := lo.Map(cards, func(card *poker.Card, _ int) int {
		return getRankValue(card.Rank)
	})
	sort.Ints(values)
	return values[0] == 10 && values[4] == 14 // 10-J-Q-K-A
}

func getRankCounts(cards poker.Cards) map[poker.Rank]int {
	counts := make(map[poker.Rank]int)
	lo.ForEach(cards, func(card *poker.Card, _ int) {
		counts[card.Rank]++
	})
	return counts
}

func hasFourOfAKind(rankCounts map[poker.Rank]int) bool {
	return lo.SomeBy(lo.Values(rankCounts), func(count int) bool {
		return count == 4
	})
}

func hasFullHouse(rankCounts map[poker.Rank]int) bool {
	counts := lo.Values(rankCounts)
	hasThree := lo.Contains(counts, 3)
	hasTwo := lo.Contains(counts, 2)
	return hasThree && hasTwo
}

func hasThreeOfAKind(rankCounts map[poker.Rank]int) bool {
	return lo.Contains(lo.Values(rankCounts), 3)
}

func hasTwoPair(rankCounts map[poker.Rank]int) bool {
	pairCount := lo.CountBy(lo.Values(rankCounts), func(count int) bool {
		return count == 2
	})
	return pairCount == 2
}

func hasOnePair(rankCounts map[poker.Rank]int) bool {
	return lo.Contains(lo.Values(rankCounts), 2)
}

func calculateValue(handRank HandRank, cards poker.Cards) int64 {
	value := int64(handRank) * 1000000000000 // Base value for hand rank

	// Add kicker values for tie-breaking
	lo.ForEach(cards, func(card *poker.Card, i int) {
		value += int64(getRankValue(card.Rank)) * int64(1000000000/(i*10+1))
	})

	return value
}

func generateCombinations(cards poker.Cards, r int) []poker.Cards {
	if r > len(cards) {
		return nil
	}

	var result []poker.Cards
	var combination poker.Cards

	var generate func(start, depth int)
	generate = func(start, depth int) {
		if depth == r {
			combo := make(poker.Cards, len(combination))
			copy(combo, combination)
			result = append(result, combo)
			return
		}

		for i := start; i < len(cards); i++ {
			combination = append(combination, cards[i])
			generate(i+1, depth+1)
			combination = combination[:len(combination)-1]
		}
	}

	generate(0, 0)
	return result
}
