package holdem

import (
	"sort"

	"github.com/ljbink/ai-poker/engine/poker"
)

// HandRank represents the ranking of a poker hand
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

// HandResult contains the evaluation result of a poker hand
type HandResult struct {
	Rank        HandRank     // The hand rank (pair, flush, etc.)
	Description string       // Human-readable description
	Value       int          // Numeric value for comparison (higher is better)
	Cards       poker.Cards  // The cards that make up the hand
	Kickers     []poker.Rank // Kicker cards for tie-breaking
}

type IHandEvaluator interface {
	EvaluateHand(holeCards []*poker.Card, communityCards poker.Cards) *HandResult
	CompareHands(hand1, hand2 *HandResult) int
}

// HandEvaluator provides methods for evaluating poker hands
type HandEvaluator struct{}

// NewHandEvaluator creates a new hand evaluator
func NewHandEvaluator() *HandEvaluator {
	return &HandEvaluator{}
}

// EvaluateHand evaluates a player's best 5-card hand from hole cards and community cards
func (e *HandEvaluator) EvaluateHand(holeCards []*poker.Card, communityCards poker.Cards) *HandResult {
	if len(holeCards) < 2 {
		return &HandResult{
			Rank:        HighCard,
			Description: "No cards",
			Value:       0,
			Cards:       poker.Cards{},
			Kickers:     []poker.Rank{},
		}
	}

	// Combine hole cards and community cards
	allCards := poker.Cards{}
	allCards.Append(holeCards...)
	allCards.Append(communityCards...)

	// Filter out nil cards
	validCards := poker.Cards{}
	for _, card := range allCards {
		if card != nil {
			validCards.Append(card)
		}
	}

	if len(validCards) < 2 {
		return &HandResult{
			Rank:        HighCard,
			Description: "Insufficient cards",
			Value:       0,
			Cards:       validCards,
			Kickers:     []poker.Rank{},
		}
	}

	// Evaluate the best 5-card hand
	return e.findBestHand(validCards)
}

// CompareHands compares two hand results and returns:
// 1 if hand1 > hand2, -1 if hand1 < hand2, 0 if equal
func (e *HandEvaluator) CompareHands(hand1, hand2 *HandResult) int {
	// First compare by hand rank
	if hand1.Rank > hand2.Rank {
		return 1
	}
	if hand1.Rank < hand2.Rank {
		return -1
	}

	// Same rank, compare by value
	if hand1.Value > hand2.Value {
		return 1
	}
	if hand1.Value < hand2.Value {
		return -1
	}

	// Same value, compare kickers
	return e.compareKickers(hand1.Kickers, hand2.Kickers)
}

// findBestHand finds the best 5-card hand from available cards
func (e *HandEvaluator) findBestHand(cards poker.Cards) *HandResult {
	if len(cards) < 5 {
		return e.evaluatePartialHand(cards)
	}

	bestHand := &HandResult{
		Rank:  HighCard,
		Value: 0,
	}

	// Try all combinations of 5 cards
	e.generateCombinations(cards, 5, func(combination poker.Cards) {
		hand := e.evaluateFiveCardHand(combination)
		if e.CompareHands(hand, bestHand) > 0 {
			bestHand = hand
		}
	})

	return bestHand
}

// evaluateFiveCardHand evaluates exactly 5 cards
func (e *HandEvaluator) evaluateFiveCardHand(cards poker.Cards) *HandResult {
	if len(cards) != 5 {
		return e.evaluatePartialHand(cards)
	}

	// Sort cards by rank (descending)
	sortedCards := make(poker.Cards, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return e.rankValue(sortedCards[i].Rank) > e.rankValue(sortedCards[j].Rank)
	})

	// Check for each hand type (highest to lowest)
	if result := e.checkRoyalFlush(sortedCards); result != nil {
		return result
	}
	if result := e.checkStraightFlush(sortedCards); result != nil {
		return result
	}
	if result := e.checkFourOfAKind(sortedCards); result != nil {
		return result
	}
	if result := e.checkFullHouse(sortedCards); result != nil {
		return result
	}
	if result := e.checkFlush(sortedCards); result != nil {
		return result
	}
	if result := e.checkStraight(sortedCards); result != nil {
		return result
	}
	if result := e.checkThreeOfAKind(sortedCards); result != nil {
		return result
	}
	if result := e.checkTwoPair(sortedCards); result != nil {
		return result
	}
	if result := e.checkOnePair(sortedCards); result != nil {
		return result
	}

	return e.checkHighCard(sortedCards)
}

// evaluatePartialHand evaluates hands with less than 5 cards
func (e *HandEvaluator) evaluatePartialHand(cards poker.Cards) *HandResult {
	if len(cards) == 0 {
		return &HandResult{
			Rank:        HighCard,
			Description: "No cards",
			Value:       0,
			Cards:       cards,
			Kickers:     []poker.Rank{},
		}
	}

	// Sort cards by rank (descending)
	sortedCards := make(poker.Cards, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return e.rankValue(sortedCards[i].Rank) > e.rankValue(sortedCards[j].Rank)
	})

	// Check for pairs/trips with available cards
	if len(cards) >= 4 {
		if result := e.checkFourOfAKind(sortedCards); result != nil {
			return result
		}
	}
	if len(cards) >= 3 {
		if result := e.checkThreeOfAKind(sortedCards); result != nil {
			return result
		}
	}
	if len(cards) >= 2 {
		if result := e.checkOnePair(sortedCards); result != nil {
			return result
		}
	}

	return e.checkHighCard(sortedCards)
}

// Hand checking functions
func (e *HandEvaluator) checkRoyalFlush(cards poker.Cards) *HandResult {
	if !e.isFlush(cards) {
		return nil
	}
	if !e.isRoyalStraight(cards) {
		return nil
	}

	return &HandResult{
		Rank:        RoyalFlush,
		Description: "Royal Flush",
		Value:       9000000,
		Cards:       cards,
		Kickers:     []poker.Rank{},
	}
}

func (e *HandEvaluator) checkStraightFlush(cards poker.Cards) *HandResult {
	if !e.isFlush(cards) {
		return nil
	}
	highCard := e.getStraightHighCard(cards)
	if highCard == poker.RankNone {
		return nil
	}

	return &HandResult{
		Rank:        StraightFlush,
		Description: "Straight Flush",
		Value:       8000000 + e.rankValue(highCard),
		Cards:       cards,
		Kickers:     []poker.Rank{highCard},
	}
}

func (e *HandEvaluator) checkFourOfAKind(cards poker.Cards) *HandResult {
	rankCounts := e.getRankCounts(cards)

	var quadRank poker.Rank
	var kicker poker.Rank

	for rank, count := range rankCounts {
		if count == 4 {
			quadRank = rank
		} else if count >= 1 {
			kicker = rank
		}
	}

	if quadRank == poker.RankNone {
		return nil
	}

	return &HandResult{
		Rank:        FourOfAKind,
		Description: "Four of a Kind",
		Value:       7000000 + e.rankValue(quadRank)*1000 + e.rankValue(kicker),
		Cards:       cards,
		Kickers:     []poker.Rank{quadRank, kicker},
	}
}

func (e *HandEvaluator) checkFullHouse(cards poker.Cards) *HandResult {
	rankCounts := e.getRankCounts(cards)

	var tripRank, pairRank poker.Rank

	for rank, count := range rankCounts {
		if count == 3 {
			tripRank = rank
		} else if count == 2 {
			pairRank = rank
		}
	}

	if tripRank == poker.RankNone || pairRank == poker.RankNone {
		return nil
	}

	return &HandResult{
		Rank:        FullHouse,
		Description: "Full House",
		Value:       6000000 + e.rankValue(tripRank)*1000 + e.rankValue(pairRank),
		Cards:       cards,
		Kickers:     []poker.Rank{tripRank, pairRank},
	}
}

func (e *HandEvaluator) checkFlush(cards poker.Cards) *HandResult {
	if !e.isFlush(cards) {
		return nil
	}

	// Get all ranks for kickers
	var kickers []poker.Rank
	for _, card := range cards {
		kickers = append(kickers, card.Rank)
	}

	// Sort kickers descending
	sort.Slice(kickers, func(i, j int) bool {
		return e.rankValue(kickers[i]) > e.rankValue(kickers[j])
	})

	value := 5000000
	for i, rank := range kickers {
		value += e.rankValue(rank) * (1000 / (i + 1))
	}

	return &HandResult{
		Rank:        Flush,
		Description: "Flush",
		Value:       value,
		Cards:       cards,
		Kickers:     kickers,
	}
}

func (e *HandEvaluator) checkStraight(cards poker.Cards) *HandResult {
	highCard := e.getStraightHighCard(cards)
	if highCard == poker.RankNone {
		return nil
	}

	return &HandResult{
		Rank:        Straight,
		Description: "Straight",
		Value:       4000000 + e.rankValue(highCard),
		Cards:       cards,
		Kickers:     []poker.Rank{highCard},
	}
}

func (e *HandEvaluator) checkThreeOfAKind(cards poker.Cards) *HandResult {
	rankCounts := e.getRankCounts(cards)

	var tripRank poker.Rank
	var kickers []poker.Rank

	for rank, count := range rankCounts {
		if count == 3 {
			tripRank = rank
		} else if count >= 1 {
			for i := 0; i < count; i++ {
				kickers = append(kickers, rank)
			}
		}
	}

	if tripRank == poker.RankNone {
		return nil
	}

	// Sort kickers descending
	sort.Slice(kickers, func(i, j int) bool {
		return e.rankValue(kickers[i]) > e.rankValue(kickers[j])
	})

	value := 3000000 + e.rankValue(tripRank)*1000
	for i, rank := range kickers {
		if i < 2 { // Only consider top 2 kickers
			value += e.rankValue(rank) * (100 / (i + 1))
		}
	}

	allKickers := []poker.Rank{tripRank}
	allKickers = append(allKickers, kickers...)

	return &HandResult{
		Rank:        ThreeOfAKind,
		Description: "Three of a Kind",
		Value:       value,
		Cards:       cards,
		Kickers:     allKickers,
	}
}

func (e *HandEvaluator) checkTwoPair(cards poker.Cards) *HandResult {
	rankCounts := e.getRankCounts(cards)

	var pairs []poker.Rank
	var kickers []poker.Rank

	for rank, count := range rankCounts {
		if count == 2 {
			pairs = append(pairs, rank)
		} else if count >= 1 {
			for i := 0; i < count; i++ {
				kickers = append(kickers, rank)
			}
		}
	}

	if len(pairs) < 2 {
		return nil
	}

	// Sort pairs descending
	sort.Slice(pairs, func(i, j int) bool {
		return e.rankValue(pairs[i]) > e.rankValue(pairs[j])
	})

	// Sort kickers descending
	sort.Slice(kickers, func(i, j int) bool {
		return e.rankValue(kickers[i]) > e.rankValue(kickers[j])
	})

	value := 2000000 + e.rankValue(pairs[0])*1000 + e.rankValue(pairs[1])*100
	if len(kickers) > 0 {
		value += e.rankValue(kickers[0])
	}

	allKickers := pairs
	allKickers = append(allKickers, kickers...)

	return &HandResult{
		Rank:        TwoPair,
		Description: "Two Pair",
		Value:       value,
		Cards:       cards,
		Kickers:     allKickers,
	}
}

func (e *HandEvaluator) checkOnePair(cards poker.Cards) *HandResult {
	rankCounts := e.getRankCounts(cards)

	var pairRank poker.Rank
	var kickers []poker.Rank

	for rank, count := range rankCounts {
		if count == 2 {
			pairRank = rank
		} else if count >= 1 {
			for i := 0; i < count; i++ {
				kickers = append(kickers, rank)
			}
		}
	}

	if pairRank == poker.RankNone {
		return nil
	}

	// Sort kickers descending
	sort.Slice(kickers, func(i, j int) bool {
		return e.rankValue(kickers[i]) > e.rankValue(kickers[j])
	})

	value := 1000000 + e.rankValue(pairRank)*1000
	for i, rank := range kickers {
		if i < 3 { // Only consider top 3 kickers
			value += e.rankValue(rank) * (100 / (i + 1))
		}
	}

	allKickers := []poker.Rank{pairRank}
	allKickers = append(allKickers, kickers...)

	return &HandResult{
		Rank:        OnePair,
		Description: "One Pair",
		Value:       value,
		Cards:       cards,
		Kickers:     allKickers,
	}
}

func (e *HandEvaluator) checkHighCard(cards poker.Cards) *HandResult {
	// Sort cards by rank descending
	sortedCards := make(poker.Cards, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return e.rankValue(sortedCards[i].Rank) > e.rankValue(sortedCards[j].Rank)
	})

	var kickers []poker.Rank
	for _, card := range sortedCards {
		kickers = append(kickers, card.Rank)
	}

	value := 0
	for i, rank := range kickers {
		if i < 5 { // Only consider top 5 cards
			value += e.rankValue(rank) * (1000 / (i + 1))
		}
	}

	return &HandResult{
		Rank:        HighCard,
		Description: "High Card",
		Value:       value,
		Cards:       cards,
		Kickers:     kickers,
	}
}

// Helper functions
func (e *HandEvaluator) isFlush(cards poker.Cards) bool {
	if len(cards) < 5 {
		return false
	}
	suit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

func (e *HandEvaluator) isRoyalStraight(cards poker.Cards) bool {
	ranks := []poker.Rank{poker.RankAce, poker.RankKing, poker.RankQueen, poker.RankJack, poker.RankTen}
	rankSet := make(map[poker.Rank]bool)
	for _, card := range cards {
		rankSet[card.Rank] = true
	}
	for _, rank := range ranks {
		if !rankSet[rank] {
			return false
		}
	}
	return true
}

func (e *HandEvaluator) getStraightHighCard(cards poker.Cards) poker.Rank {
	if len(cards) < 5 {
		return poker.RankNone
	}

	ranks := make([]int, 0)
	rankSet := make(map[int]bool)

	for _, card := range cards {
		rank := e.rankValue(card.Rank)
		if !rankSet[rank] {
			ranks = append(ranks, rank)
			rankSet[rank] = true
		}
	}

	sort.Ints(ranks)

	// Check for regular straight
	if len(ranks) >= 5 {
		for i := len(ranks) - 5; i >= 0; i-- {
			if ranks[i+4]-ranks[i] == 4 {
				return e.valueToRank(ranks[i+4])
			}
		}
	}

	// Check for A-2-3-4-5 straight (wheel)
	if rankSet[14] && rankSet[2] && rankSet[3] && rankSet[4] && rankSet[5] {
		return poker.RankFive // 5-high straight
	}

	return poker.RankNone
}

func (e *HandEvaluator) getRankCounts(cards poker.Cards) map[poker.Rank]int {
	counts := make(map[poker.Rank]int)
	for _, card := range cards {
		counts[card.Rank]++
	}
	return counts
}

func (e *HandEvaluator) rankValue(rank poker.Rank) int {
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

func (e *HandEvaluator) valueToRank(value int) poker.Rank {
	switch value {
	case 14:
		return poker.RankAce
	case 13:
		return poker.RankKing
	case 12:
		return poker.RankQueen
	case 11:
		return poker.RankJack
	case 10:
		return poker.RankTen
	case 9:
		return poker.RankNine
	case 8:
		return poker.RankEight
	case 7:
		return poker.RankSeven
	case 6:
		return poker.RankSix
	case 5:
		return poker.RankFive
	case 4:
		return poker.RankFour
	case 3:
		return poker.RankThree
	case 2:
		return poker.RankTwo
	default:
		return poker.RankNone
	}
}

func (e *HandEvaluator) compareKickers(kickers1, kickers2 []poker.Rank) int {
	maxLen := len(kickers1)
	if len(kickers2) > maxLen {
		maxLen = len(kickers2)
	}

	for i := 0; i < maxLen; i++ {
		val1 := 0
		val2 := 0

		if i < len(kickers1) {
			val1 = e.rankValue(kickers1[i])
		}
		if i < len(kickers2) {
			val2 = e.rankValue(kickers2[i])
		}

		if val1 > val2 {
			return 1
		}
		if val1 < val2 {
			return -1
		}
	}
	return 0
}

func (e *HandEvaluator) generateCombinations(cards poker.Cards, k int, callback func(poker.Cards)) {
	n := len(cards)
	if k > n {
		return
	}

	indices := make([]int, k)
	for i := range indices {
		indices[i] = i
	}

	for {
		combination := make(poker.Cards, k)
		for i, idx := range indices {
			combination[i] = cards[idx]
		}
		callback(combination)

		// Generate next combination
		i := k - 1
		for i >= 0 && indices[i] == n-k+i {
			i--
		}
		if i < 0 {
			break
		}
		indices[i]++
		for j := i + 1; j < k; j++ {
			indices[j] = indices[j-1] + 1
		}
	}
}

// HandRankToString converts hand rank to string
func HandRankToString(rank HandRank) string {
	switch rank {
	case HighCard:
		return "High Card"
	case OnePair:
		return "One Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}
