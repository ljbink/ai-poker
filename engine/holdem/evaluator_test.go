package holdem

import (
	"testing"

	"github.com/ljbink/ai-poker/engine/poker"
)

func TestNewHandEvaluator(t *testing.T) {
	evaluator := NewHandEvaluator()
	if evaluator == nil {
		t.Fatal("NewHandEvaluator returned nil")
	}
}

func TestHandRankValues(t *testing.T) {
	// Test that hand ranks have correct ordering
	expectedRanks := []HandRank{
		HighCard,
		OnePair,
		TwoPair,
		ThreeOfAKind,
		Straight,
		Flush,
		FullHouse,
		FourOfAKind,
		StraightFlush,
		RoyalFlush,
	}

	for i, rank := range expectedRanks {
		if int(rank) != i {
			t.Errorf("Expected rank %d to have value %d, got %d", i, i, int(rank))
		}
	}
}

func TestEvaluateHandEmpty(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test with no cards
	result := evaluator.EvaluateHand([]*poker.Card{}, poker.Cards{})
	if result.Rank != HighCard {
		t.Errorf("Expected HighCard for no cards, got %d", result.Rank)
	}
	if result.Value != 0 {
		t.Errorf("Expected value 0 for no cards, got %d", result.Value)
	}

	// Test with one card only
	cards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
	}
	result = evaluator.EvaluateHand(cards, poker.Cards{})
	if result.Rank != HighCard {
		t.Errorf("Expected HighCard for one card, got %d", result.Rank)
	}
}

func TestEvaluateHighCard(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test high card hand: A-K-Q-J-9 of different suits
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankQueen, Suit: poker.SuitDiamond},
		{Rank: poker.RankJack, Suit: poker.SuitClub},
		{Rank: poker.RankNine, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != HighCard {
		t.Errorf("Expected HighCard, got %d", result.Rank)
	}
	if len(result.Cards) != 5 {
		t.Errorf("Expected 5 cards in result, got %d", len(result.Cards))
	}
}

func TestEvaluateOnePair(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test pair of Aces
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankAce, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankKing, Suit: poker.SuitDiamond},
		{Rank: poker.RankQueen, Suit: poker.SuitClub},
		{Rank: poker.RankJack, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != OnePair {
		t.Errorf("Expected OnePair, got %d", result.Rank)
	}
	if len(result.Cards) != 5 {
		t.Errorf("Expected 5 cards in result, got %d", len(result.Cards))
	}
}

func TestEvaluateTwoPair(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test two pair: Aces and Kings
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitDiamond},
		{Rank: poker.RankKing, Suit: poker.SuitClub},
		{Rank: poker.RankQueen, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != TwoPair {
		t.Errorf("Expected TwoPair, got %d", result.Rank)
	}
}

func TestEvaluateThreeOfAKind(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test three of a kind: Three Aces
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankAce, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitDiamond},
		{Rank: poker.RankKing, Suit: poker.SuitClub},
		{Rank: poker.RankQueen, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != ThreeOfAKind {
		t.Errorf("Expected ThreeOfAKind, got %d", result.Rank)
	}
}

func TestEvaluateStraight(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test straight: A-K-Q-J-T
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankQueen, Suit: poker.SuitDiamond},
		{Rank: poker.RankJack, Suit: poker.SuitClub},
		{Rank: poker.RankTen, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != Straight {
		t.Errorf("Expected Straight, got %d", result.Rank)
	}
}

func TestEvaluateWheelStraight(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test wheel straight: A-2-3-4-5 (Ace low)
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankTwo, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankThree, Suit: poker.SuitDiamond},
		{Rank: poker.RankFour, Suit: poker.SuitClub},
		{Rank: poker.RankFive, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != Straight {
		t.Errorf("Expected Straight (wheel), got %d", result.Rank)
	}
}

func TestEvaluateFlush(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test flush: All spades
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitSpade},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankQueen, Suit: poker.SuitSpade},
		{Rank: poker.RankJack, Suit: poker.SuitSpade},
		{Rank: poker.RankNine, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != Flush {
		t.Errorf("Expected Flush, got %d", result.Rank)
	}
}

func TestEvaluateFullHouse(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test full house: Three Aces and two Kings
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankAce, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitDiamond},
		{Rank: poker.RankKing, Suit: poker.SuitClub},
		{Rank: poker.RankKing, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != FullHouse {
		t.Errorf("Expected FullHouse, got %d", result.Rank)
	}
}

func TestEvaluateFourOfAKind(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test four of a kind: Four Aces
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankAce, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitDiamond},
		{Rank: poker.RankAce, Suit: poker.SuitClub},
		{Rank: poker.RankKing, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != FourOfAKind {
		t.Errorf("Expected FourOfAKind, got %d", result.Rank)
	}
}

func TestEvaluateStraightFlush(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test straight flush: 9-T-J-Q-K of spades
	holeCards := []*poker.Card{
		{Rank: poker.RankNine, Suit: poker.SuitSpade},
		{Rank: poker.RankTen, Suit: poker.SuitSpade},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankJack, Suit: poker.SuitSpade},
		{Rank: poker.RankQueen, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != StraightFlush {
		t.Errorf("Expected StraightFlush, got %d", result.Rank)
	}
}

func TestEvaluateRoyalFlush(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test royal flush: T-J-Q-K-A of spades
	holeCards := []*poker.Card{
		{Rank: poker.RankTen, Suit: poker.SuitSpade},
		{Rank: poker.RankJack, Suit: poker.SuitSpade},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankQueen, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitSpade},
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != RoyalFlush {
		t.Errorf("Expected RoyalFlush, got %d", result.Rank)
	}
}

func TestCompareHands(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Create different ranked hands
	highCard := &HandResult{Rank: HighCard, Value: 100}
	onePair := &HandResult{Rank: OnePair, Value: 200}
	twoPair := &HandResult{Rank: TwoPair, Value: 300}

	// Test higher rank wins
	if evaluator.CompareHands(onePair, highCard) != 1 {
		t.Error("OnePair should beat HighCard")
	}
	if evaluator.CompareHands(highCard, onePair) != -1 {
		t.Error("HighCard should lose to OnePair")
	}
	if evaluator.CompareHands(twoPair, onePair) != 1 {
		t.Error("TwoPair should beat OnePair")
	}

	// Test equal ranks with different values
	onePair1 := &HandResult{Rank: OnePair, Value: 200}
	onePair2 := &HandResult{Rank: OnePair, Value: 300}
	if evaluator.CompareHands(onePair2, onePair1) != 1 {
		t.Error("Higher value OnePair should win")
	}

	// Test equal hands
	onePair3 := &HandResult{Rank: OnePair, Value: 200, Kickers: []poker.Rank{poker.RankKing}}
	onePair4 := &HandResult{Rank: OnePair, Value: 200, Kickers: []poker.Rank{poker.RankKing}}
	if evaluator.CompareHands(onePair3, onePair4) != 0 {
		t.Error("Identical hands should be equal")
	}
}

func TestKickerComparison(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test kicker comparison with same rank and value
	hand1 := &HandResult{
		Rank:    OnePair,
		Value:   200,
		Kickers: []poker.Rank{poker.RankKing, poker.RankQueen},
	}
	hand2 := &HandResult{
		Rank:    OnePair,
		Value:   200,
		Kickers: []poker.Rank{poker.RankKing, poker.RankJack},
	}

	if evaluator.CompareHands(hand1, hand2) != 1 {
		t.Error("Hand with Queen kicker should beat hand with Jack kicker")
	}

	// Test first kicker difference
	hand3 := &HandResult{
		Rank:    OnePair,
		Value:   200,
		Kickers: []poker.Rank{poker.RankAce},
	}
	hand4 := &HandResult{
		Rank:    OnePair,
		Value:   200,
		Kickers: []poker.Rank{poker.RankKing},
	}

	if evaluator.CompareHands(hand3, hand4) != 1 {
		t.Error("Hand with Ace kicker should beat hand with King kicker")
	}
}

func TestHandRankToString(t *testing.T) {
	testCases := []struct {
		rank     HandRank
		expected string
	}{
		{HighCard, "High Card"},
		{OnePair, "One Pair"},
		{TwoPair, "Two Pair"},
		{ThreeOfAKind, "Three of a Kind"},
		{Straight, "Straight"},
		{Flush, "Flush"},
		{FullHouse, "Full House"},
		{FourOfAKind, "Four of a Kind"},
		{StraightFlush, "Straight Flush"},
		{RoyalFlush, "Royal Flush"},
	}

	for _, tc := range testCases {
		result := HandRankToString(tc.rank)
		if result != tc.expected {
			t.Errorf("Expected HandRankToString(%d) = %s, got %s", tc.rank, tc.expected, result)
		}
	}

	// Test unknown rank
	unknown := HandRank(99)
	result := HandRankToString(unknown)
	if result != "Unknown" {
		t.Errorf("Expected 'Unknown' for invalid rank, got %s", result)
	}
}

func TestEvaluateHandWithNilCards(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test with nil cards in the mix
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		nil,
	}
	communityCards := poker.Cards{
		{Rank: poker.RankKing, Suit: poker.SuitHeart},
		nil,
		{Rank: poker.RankQueen, Suit: poker.SuitDiamond},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	// Should handle nil cards gracefully
	if result == nil {
		t.Error("EvaluateHand should not return nil")
	}
}

func TestEvaluateSevenCards(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test with full 7 cards (2 hole + 5 community)
	// This should form a straight: A-K-Q-J-T, which is higher than two pair
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankAce, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankKing, Suit: poker.SuitDiamond},
		{Rank: poker.RankKing, Suit: poker.SuitClub},
		{Rank: poker.RankQueen, Suit: poker.SuitSpade},
		{Rank: poker.RankJack, Suit: poker.SuitHeart},
		{Rank: poker.RankTen, Suit: poker.SuitDiamond},
	}

	result := evaluator.EvaluateHand(holeCards, communityCards)
	// The best hand is A-K-Q-J-T straight, not two pair
	if result.Rank != Straight {
		t.Errorf("Expected Straight from 7 cards (A-K-Q-J-T), got %d", result.Rank)
	}
	if len(result.Cards) != 5 {
		t.Errorf("Expected 5 cards in result even with 7 input cards, got %d", len(result.Cards))
	}
}

func TestEvaluatePartialHands(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Test with only 2 cards (preflop)
	holeCards := []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankAce, Suit: poker.SuitHeart},
	}

	result := evaluator.EvaluateHand(holeCards, poker.Cards{})
	if result.Rank != OnePair {
		t.Errorf("Expected OnePair from pocket pair, got %d", result.Rank)
	}

	// Test with 3 cards (no community cards yet)
	holeCards = []*poker.Card{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitHeart},
	}
	communityCards := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitDiamond},
	}

	result = evaluator.EvaluateHand(holeCards, communityCards)
	if result.Rank != OnePair {
		t.Errorf("Expected OnePair from 3 cards with pair, got %d", result.Rank)
	}
}

func TestEvaluatePartialHandsMoreBranches(t *testing.T) {
	evaluator := NewHandEvaluator()

	// 4 cards with no pairs/trips/quads should yield HighCard via partial evaluator
	cards := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitHeart},
		{Rank: poker.RankNine, Suit: poker.SuitClub},
		{Rank: poker.RankFour, Suit: poker.SuitDiamond},
	}
	res := evaluator.evaluatePartialHand(cards)
	if res.Rank != HighCard {
		t.Errorf("Expected HighCard for 4 mixed cards, got %v", res.Rank)
	}

	// 3 cards making trips should be detected in partial evaluator
	trips := poker.Cards{
		{Rank: poker.RankTen, Suit: poker.SuitSpade},
		{Rank: poker.RankTen, Suit: poker.SuitHeart},
		{Rank: poker.RankTen, Suit: poker.SuitClub},
	}
	res = evaluator.evaluatePartialHand(trips)
	if res.Rank != ThreeOfAKind {
		t.Errorf("Expected ThreeOfAKind for three tens, got %v", res.Rank)
	}

	// 4 cards making quads should be detected in partial evaluator
	quads := poker.Cards{
		{Rank: poker.RankFive, Suit: poker.SuitSpade},
		{Rank: poker.RankFive, Suit: poker.SuitHeart},
		{Rank: poker.RankFive, Suit: poker.SuitDiamond},
		{Rank: poker.RankFive, Suit: poker.SuitClub},
	}
	res = evaluator.evaluatePartialHand(quads)
	if res.Rank != FourOfAKind {
		t.Errorf("Expected FourOfAKind for four fives, got %v", res.Rank)
	}
}

func TestRankValueAndValueToRankMappings(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Check rankValue mapping for all standard ranks
	rankToVal := map[poker.Rank]int{
		poker.RankTwo:   2,
		poker.RankThree: 3,
		poker.RankFour:  4,
		poker.RankFive:  5,
		poker.RankSix:   6,
		poker.RankSeven: 7,
		poker.RankEight: 8,
		poker.RankNine:  9,
		poker.RankTen:   10,
		poker.RankJack:  11,
		poker.RankQueen: 12,
		poker.RankKing:  13,
		poker.RankAce:   14,
	}
	for r, v := range rankToVal {
		if got := evaluator.rankValue(r); got != v {
			t.Errorf("rankValue(%v) = %d, want %d", r, got, v)
		}
	}

	// Check valueToRank mapping including out-of-range
	valToRank := map[int]poker.Rank{
		2:  poker.RankTwo,
		3:  poker.RankThree,
		4:  poker.RankFour,
		5:  poker.RankFive,
		6:  poker.RankSix,
		7:  poker.RankSeven,
		8:  poker.RankEight,
		9:  poker.RankNine,
		10: poker.RankTen,
		11: poker.RankJack,
		12: poker.RankQueen,
		13: poker.RankKing,
		14: poker.RankAce,
	}
	for v, r := range valToRank {
		if got := evaluator.valueToRank(v); got != r {
			t.Errorf("valueToRank(%d) = %v, want %v", v, got, r)
		}
	}
	// Out of range values
	if got := evaluator.valueToRank(1); got != poker.RankNone {
		t.Errorf("valueToRank(1) = %v, want RankNone", got)
	}
	if got := evaluator.valueToRank(15); got != poker.RankNone {
		t.Errorf("valueToRank(15) = %v, want RankNone", got)
	}
}

func TestIsFlushNegativeCases(t *testing.T) {
	evaluator := NewHandEvaluator()

	// Fewer than 5 cards should not be a flush
	notEnough := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitSpade},
		{Rank: poker.RankQueen, Suit: poker.SuitSpade},
		{Rank: poker.RankJack, Suit: poker.SuitSpade},
	}
	if evaluator.isFlush(notEnough) {
		t.Error("isFlush should be false for <5 cards")
	}

	// Five cards mixed suits should not be a flush
	mixed := poker.Cards{
		{Rank: poker.RankAce, Suit: poker.SuitSpade},
		{Rank: poker.RankKing, Suit: poker.SuitHeart},
		{Rank: poker.RankQueen, Suit: poker.SuitClub},
		{Rank: poker.RankJack, Suit: poker.SuitDiamond},
		{Rank: poker.RankTen, Suit: poker.SuitSpade},
	}
	if evaluator.isFlush(mixed) {
		t.Error("isFlush should be false for mixed suits")
	}
}
