package poker_test

import (
	"github.com/ljbink/ai-poker/engine/poker"
	"testing"
)

func TestCardString(t *testing.T) {
	tests := []struct {
		card     poker.Card
		expected string
	}{
		{poker.Card{poker.SuitHeart, poker.RankAce}, "🂱"},
		{poker.Card{poker.SuitDiamond, poker.RankTen}, "🃊"},
		{poker.Card{poker.SuitClub, poker.RankKing}, "🃞"},
		{poker.Card{poker.SuitSpade, poker.RankQueen}, "🂽"},
		{poker.Card{poker.SuitNone, poker.RankNone}, "🂠"},
	}

	for _, test := range tests {
		result := test.card.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestNewCard(t *testing.T) {
	card := poker.NewCard(poker.SuitSpade, poker.RankAce)
	if card.Suit != poker.SuitSpade || card.Rank != poker.RankAce {
		t.Errorf("Expected Suit: %d, Rank: %d, got Suit: %d, Rank: %d", poker.SuitSpade, poker.RankAce, card.Suit, card.Rank)
	}
}
