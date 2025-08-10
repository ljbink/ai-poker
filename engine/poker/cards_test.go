package poker_test

import (
	"strings"
	"testing"

	"github.com/ljbink/ai-poker/engine/poker"
)

func TestCardsLength(t *testing.T) {
	// Test empty deck
	emptyDeck := poker.Cards{}
	if emptyDeck.Length() != 0 {
		t.Errorf("Expected empty deck length of 0, got %d", emptyDeck.Length())
	}

	// Test full deck
	deck := poker.NewDeckCards()
	if deck.Length() != 54 {
		t.Errorf("Expected deck length of 54, got %d", deck.Length())
	}

	// Test deck with some cards
	smallDeck := poker.Cards{}
	smallDeck.Append(poker.NewCard(poker.SuitHeart, poker.RankAce))
	smallDeck.Append(poker.NewCard(poker.SuitSpade, poker.RankKing))
	if smallDeck.Length() != 2 {
		t.Errorf("Expected small deck length of 2, got %d", smallDeck.Length())
	}
}

func TestCardsShuffle(t *testing.T) {
	// Test shuffling with multiple attempts to ensure it's actually shuffling
	deck := poker.NewDeckCards()
	originalOrder := deck.String()

	shuffled := false
	// Try multiple shuffles since there's a very small chance the shuffle could result in the same order
	for i := 0; i < 10; i++ {
		deck.Shuffle()
		if deck.String() != originalOrder {
			shuffled = true
			break
		}
	}

	if !shuffled {
		t.Error("Deck was not shuffled after 10 attempts")
	}

	// Test shuffling empty deck (should not panic)
	emptyDeck := poker.Cards{}
	emptyDeck.Shuffle()
	if emptyDeck.Length() != 0 {
		t.Error("Empty deck length changed after shuffle")
	}

	// Test shuffling single card deck (should not panic)
	singleCardDeck := poker.Cards{}
	singleCardDeck.Append(poker.NewCard(poker.SuitHeart, poker.RankAce))
	singleCardDeck.Shuffle()
	if singleCardDeck.Length() != 1 {
		t.Error("Single card deck length changed after shuffle")
	}
}

func TestCardsAppend(t *testing.T) {
	// Test appending single card
	deck := poker.Cards{}
	card := poker.NewCard(poker.SuitHeart, poker.RankAce)
	deck.Append(card)
	if deck.Length() != 1 {
		t.Errorf("Expected deck length of 1, got %d", deck.Length())
	}

	// Test appending multiple cards at once
	card2 := poker.NewCard(poker.SuitSpade, poker.RankKing)
	card3 := poker.NewCard(poker.SuitDiamond, poker.RankQueen)
	deck.Append(card2, card3)
	if deck.Length() != 3 {
		t.Errorf("Expected deck length of 3, got %d", deck.Length())
	}

	// Test appending no cards (should not change length)
	originalLength := deck.Length()
	deck.Append()
	if deck.Length() != originalLength {
		t.Errorf("Expected deck length to remain %d, got %d", originalLength, deck.Length())
	}

	// Test appending nil cards
	deck.Append(nil)
	if deck.Length() != originalLength+1 {
		t.Errorf("Expected deck length to be %d after appending nil, got %d", originalLength+1, deck.Length())
	}
}

func TestCardsRemove(t *testing.T) {
	// Test removing from full deck
	deck1 := poker.NewDeckCards()
	for _, card := range deck1 {
		deck2 := poker.NewDeckCards()
		removed := deck2.Remove(card)
		if removed != 1 {
			t.Errorf("Expected to remove 1 card, removed %d", removed)
		}
		if deck2.Length() != 53 {
			t.Errorf("Expected deck length of 53 after removal, got %d", deck2.Length())
		}
	}

	// Test removing multiple cards at once
	deck := poker.NewDeckCards()
	card1 := poker.NewCard(poker.SuitHeart, poker.RankAce)
	card2 := poker.NewCard(poker.SuitSpade, poker.RankKing)
	removed := deck.Remove(card1, card2)
	if removed != 2 {
		t.Errorf("Expected to remove 2 cards, removed %d", removed)
	}
	if deck.Length() != 52 {
		t.Errorf("Expected deck length of 52 after removing 2 cards, got %d", deck.Length())
	}

	// Test removing non-existent card
	deck = poker.NewDeckCards()
	nonExistentCard := poker.NewCard(poker.Suit(99), poker.Rank(99))
	removed = deck.Remove(nonExistentCard)
	if removed != 0 {
		t.Errorf("Expected to remove 0 non-existent cards, removed %d", removed)
	}
	if deck.Length() != 54 {
		t.Errorf("Expected deck length to remain 54, got %d", deck.Length())
	}

	// Test removing from empty deck
	emptyDeck := poker.Cards{}
	card := poker.NewCard(poker.SuitHeart, poker.RankAce)
	removed = emptyDeck.Remove(card)
	if removed != 0 {
		t.Errorf("Expected to remove 0 cards from empty deck, removed %d", removed)
	}

	// Test removing no cards (should not change deck)
	deck = poker.NewDeckCards()
	originalLength := deck.Length()
	removed = deck.Remove()
	if removed != 0 {
		t.Errorf("Expected to remove 0 cards when removing nothing, removed %d", removed)
	}
	if deck.Length() != originalLength {
		t.Errorf("Expected deck length to remain %d, got %d", originalLength, deck.Length())
	}

	// Test removing nil card
	deck = poker.Cards{}
	deck.Append(poker.NewCard(poker.SuitHeart, poker.RankAce))
	removed = deck.Remove(nil)
	if removed != 0 {
		t.Errorf("Expected to remove 0 nil cards from deck with only non-nil, removed %d", removed)
	}

	// Test removing nil when deck contains nil entries
	deck = poker.Cards{}
	deck.Append(nil, poker.NewCard(poker.SuitHeart, poker.RankAce), nil)
	removed = deck.Remove(nil)
	if removed != 2 {
		t.Errorf("Expected to remove 2 nil cards, removed %d", removed)
	}
	// Ensure non-nil card remains
	if deck.Length() != 1 {
		t.Errorf("Expected deck length 1 after removing two nils, got %d", deck.Length())
	}

	// Test removing duplicate cards
	deck = poker.Cards{}
	card = poker.NewCard(poker.SuitHeart, poker.RankAce)
	deck.Append(card, card, card) // Add same card 3 times
	removed = deck.Remove(card)
	if removed != 3 {
		t.Errorf("Expected to remove 3 duplicate cards, removed %d", removed)
	}
	if deck.Length() != 0 {
		t.Errorf("Expected empty deck after removing all cards, got length %d", deck.Length())
	}
}

func TestCardsString(t *testing.T) {
	// Test empty deck string
	emptyDeck := poker.Cards{}
	if emptyDeck.String() != "" {
		t.Errorf("Expected empty string for empty deck, got %q", emptyDeck.String())
	}

	// Test single card string
	singleDeck := poker.Cards{}
	card := poker.NewCard(poker.SuitHeart, poker.RankAce)
	singleDeck.Append(card)
	expected := card.String()
	if singleDeck.String() != expected {
		t.Errorf("Expected %q, got %q", expected, singleDeck.String())
	}

	// Test multiple cards string (should be space-separated)
	multiDeck := poker.Cards{}
	card1 := poker.NewCard(poker.SuitHeart, poker.RankAce)
	card2 := poker.NewCard(poker.SuitSpade, poker.RankKing)
	card3 := poker.NewCard(poker.SuitDiamond, poker.RankQueen)
	multiDeck.Append(card1, card2, card3)

	result := multiDeck.String()
	expectedParts := []string{card1.String(), card2.String(), card3.String()}
	expected = strings.Join(expectedParts, " ")

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Verify it contains spaces between cards
	if !strings.Contains(result, " ") {
		t.Error("Expected spaces between cards in string representation")
	}

	// Test with nil card (should handle gracefully)
	deckWithNil := poker.Cards{}
	deckWithNil.Append(card1, nil, card2)
	resultWithNil := deckWithNil.String()
	// This should not panic and should include some representation
	if len(resultWithNil) == 0 {
		t.Error("Expected non-empty string even with nil card")
	}
}

func TestNewDeckCards(t *testing.T) {
	deck := poker.NewDeckCards()

	// Test total count (52 regular cards + 2 jokers)
	if deck.Length() != 54 {
		t.Errorf("Expected 54 cards in new deck, got %d", deck.Length())
	}

	// Test that it contains all expected cards
	expectedSuits := []poker.Suit{poker.SuitHeart, poker.SuitDiamond, poker.SuitClub, poker.SuitSpade}
	expectedRanks := []poker.Rank{
		poker.RankAce, poker.RankTwo, poker.RankThree, poker.RankFour,
		poker.RankFive, poker.RankSix, poker.RankSeven, poker.RankEight,
		poker.RankNine, poker.RankTen, poker.RankJack, poker.RankQueen,
		poker.RankKing,
	}

	// Count cards by suit and rank
	suitCounts := make(map[poker.Suit]int)
	rankCounts := make(map[poker.Rank]int)

	for _, card := range deck {
		if card != nil {
			suitCounts[card.Suit]++
			rankCounts[card.Rank]++
		}
	}

	// Each suit should have 13 cards (except SuitNone which has 2 jokers)
	for _, suit := range expectedSuits {
		if suitCounts[suit] != 13 {
			t.Errorf("Expected 13 cards of suit %d, got %d", suit, suitCounts[suit])
		}
	}

	// SuitNone should have 2 cards (jokers)
	if suitCounts[poker.SuitNone] != 2 {
		t.Errorf("Expected 2 joker cards, got %d", suitCounts[poker.SuitNone])
	}

	// Each rank should have 4 cards (except jokers)
	for _, rank := range expectedRanks {
		if rankCounts[rank] != 4 {
			t.Errorf("Expected 4 cards of rank %d, got %d", rank, rankCounts[rank])
		}
	}

	// Jokers should have 1 card each
	if rankCounts[poker.RankJoker] != 1 {
		t.Errorf("Expected 1 Joker card, got %d", rankCounts[poker.RankJoker])
	}
	if rankCounts[poker.RankColoredJoker] != 1 {
		t.Errorf("Expected 1 ColoredJoker card, got %d", rankCounts[poker.RankColoredJoker])
	}

	// Test that multiple calls return independent decks
	deck1 := poker.NewDeckCards()
	deck2 := poker.NewDeckCards()

	// Remove a card from deck1
	if deck1.Length() > 0 {
		cardToRemove := deck1[0]
		deck1.Remove(cardToRemove)
	}

	// deck2 should still be full
	if deck2.Length() != 54 {
		t.Error("NewDeckCards should return independent instances")
	}
}

func TestCardsEdgeCases(t *testing.T) {
	// Test operations on very large decks
	largeDeck := poker.Cards{}
	for i := 0; i < 1000; i++ {
		largeDeck.Append(poker.NewCard(poker.SuitHeart, poker.RankAce))
	}

	if largeDeck.Length() != 1000 {
		t.Errorf("Expected large deck to have 1000 cards, got %d", largeDeck.Length())
	}

	// Shuffle large deck (should not panic)
	largeDeck.Shuffle()
	if largeDeck.Length() != 1000 {
		t.Error("Large deck length changed after shuffle")
	}

	// Remove all cards from large deck
	cardToRemove := poker.NewCard(poker.SuitHeart, poker.RankAce)
	removed := largeDeck.Remove(cardToRemove)
	if removed != 1000 {
		t.Errorf("Expected to remove 1000 identical cards, removed %d", removed)
	}
	if largeDeck.Length() != 0 {
		t.Errorf("Expected empty deck after removing all cards, got %d", largeDeck.Length())
	}
}

func TestCardsBooleanAndCountsClosures(t *testing.T) {
	// Test the type aliases exist and can be used
	var boolClosure poker.CardBooleanClosure = func(c *poker.Card) bool {
		return c != nil && c.Suit == poker.SuitHeart
	}

	var countsClosure poker.CardCountsClosure = func(val poker.Rank, count int) bool {
		return count > 2
	}

	// Test that we can use these closures
	card := poker.NewCard(poker.SuitHeart, poker.RankAce)
	if !boolClosure(card) {
		t.Error("CardBooleanClosure should work with heart cards")
	}

	if !countsClosure(poker.RankAce, 5) {
		t.Error("CardCountsClosure should work with count > 2")
	}

	// Test with nil card
	if boolClosure(nil) {
		t.Error("CardBooleanClosure should handle nil gracefully")
	}
}
