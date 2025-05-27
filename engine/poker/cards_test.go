package poker_test

import (
	"testing"

	"github.com/ljbink/ai-poker/engine/poker"
)

func TestCardsLength(t *testing.T) {
	deck := poker.NewDeckCards()
	if deck.Length() != 54 {
		t.Errorf("Expected deck length of 54, got %d", deck.Length())
	}
}

func TestCardsShuffle(t *testing.T) {
	deck := poker.NewDeckCards()
	deckBeforeShuffle := deck.String()
	deck.Shuffle()
	deckAfterShuffle := deck.String()
	if deckBeforeShuffle == deckAfterShuffle {
		t.Errorf("Deck was not shuffled")
	}
}

func TestCardsAppend(t *testing.T) {
	deck := poker.Cards{}
	card := poker.NewCard(poker.SuitHeart, poker.RankAce)
	deck.Append(card)
	if deck.Length() != 1 {
		t.Errorf("Expected deck length of 1, got %d", deck.Length())
	}
}

func TestCardsRemove(t *testing.T) {
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

}
