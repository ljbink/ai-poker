package poker_test

import (
	"fmt"
	"testing"

	"github.com/ljbink/ai-poker/engine/poker"
)

func TestCardString(t *testing.T) {
	tests := []struct {
		card     poker.Card
		expected string
	}{
		// Heart cards
		{poker.Card{poker.SuitHeart, poker.RankAce}, "🂱"},
		{poker.Card{poker.SuitHeart, poker.RankTwo}, "🂲"},
		{poker.Card{poker.SuitHeart, poker.RankThree}, "🂳"},
		{poker.Card{poker.SuitHeart, poker.RankFour}, "🂴"},
		{poker.Card{poker.SuitHeart, poker.RankFive}, "🂵"},
		{poker.Card{poker.SuitHeart, poker.RankSix}, "🂶"},
		{poker.Card{poker.SuitHeart, poker.RankSeven}, "🂷"},
		{poker.Card{poker.SuitHeart, poker.RankEight}, "🂸"},
		{poker.Card{poker.SuitHeart, poker.RankNine}, "🂹"},
		{poker.Card{poker.SuitHeart, poker.RankTen}, "🂺"},
		{poker.Card{poker.SuitHeart, poker.RankJack}, "🂫"},
		{poker.Card{poker.SuitHeart, poker.RankQueen}, "🂭"},
		{poker.Card{poker.SuitHeart, poker.RankKing}, "🂮"},

		// Diamond cards
		{poker.Card{poker.SuitDiamond, poker.RankAce}, "🃁"},
		{poker.Card{poker.SuitDiamond, poker.RankTwo}, "🃂"},
		{poker.Card{poker.SuitDiamond, poker.RankThree}, "🃃"},
		{poker.Card{poker.SuitDiamond, poker.RankFour}, "🃄"},
		{poker.Card{poker.SuitDiamond, poker.RankFive}, "🃅"},
		{poker.Card{poker.SuitDiamond, poker.RankSix}, "🃆"},
		{poker.Card{poker.SuitDiamond, poker.RankSeven}, "🃇"},
		{poker.Card{poker.SuitDiamond, poker.RankEight}, "🃈"},
		{poker.Card{poker.SuitDiamond, poker.RankNine}, "🃉"},
		{poker.Card{poker.SuitDiamond, poker.RankTen}, "🃊"},
		{poker.Card{poker.SuitDiamond, poker.RankJack}, "🃋"},
		{poker.Card{poker.SuitDiamond, poker.RankQueen}, "🃍"},
		{poker.Card{poker.SuitDiamond, poker.RankKing}, "🃎"},

		// Club cards
		{poker.Card{poker.SuitClub, poker.RankAce}, "🃑"},
		{poker.Card{poker.SuitClub, poker.RankTwo}, "🃒"},
		{poker.Card{poker.SuitClub, poker.RankThree}, "🃓"},
		{poker.Card{poker.SuitClub, poker.RankFour}, "🃔"},
		{poker.Card{poker.SuitClub, poker.RankFive}, "🃕"},
		{poker.Card{poker.SuitClub, poker.RankSix}, "🃖"},
		{poker.Card{poker.SuitClub, poker.RankSeven}, "🃗"},
		{poker.Card{poker.SuitClub, poker.RankEight}, "🃘"},
		{poker.Card{poker.SuitClub, poker.RankNine}, "🃙"},
		{poker.Card{poker.SuitClub, poker.RankTen}, "🃚"},
		{poker.Card{poker.SuitClub, poker.RankJack}, "🃛"},
		{poker.Card{poker.SuitClub, poker.RankQueen}, "🃝"},
		{poker.Card{poker.SuitClub, poker.RankKing}, "🃞"},

		// Spade cards
		{poker.Card{poker.SuitSpade, poker.RankAce}, "🂡"},
		{poker.Card{poker.SuitSpade, poker.RankTwo}, "🂢"},
		{poker.Card{poker.SuitSpade, poker.RankThree}, "🂣"},
		{poker.Card{poker.SuitSpade, poker.RankFour}, "🂤"},
		{poker.Card{poker.SuitSpade, poker.RankFive}, "🂥"},
		{poker.Card{poker.SuitSpade, poker.RankSix}, "🂦"},
		{poker.Card{poker.SuitSpade, poker.RankSeven}, "🂧"},
		{poker.Card{poker.SuitSpade, poker.RankEight}, "🂨"},
		{poker.Card{poker.SuitSpade, poker.RankNine}, "🂩"},
		{poker.Card{poker.SuitSpade, poker.RankTen}, "🂪"},
		{poker.Card{poker.SuitSpade, poker.RankJack}, "🂻"},
		{poker.Card{poker.SuitSpade, poker.RankQueen}, "🂽"},
		{poker.Card{poker.SuitSpade, poker.RankKing}, "🂾"},

		// Special cards
		{poker.Card{poker.SuitNone, poker.RankNone}, "🂠"},
		{poker.Card{poker.SuitNone, poker.RankJoker}, "🃟"},
		{poker.Card{poker.SuitNone, poker.RankColoredJoker}, "🃏"},
	}

	for _, test := range tests {
		result := test.card.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s for card %+v", test.expected, result, test.card)
		}
	}
}

func TestCardStringFallback(t *testing.T) {
	// Test the fallback case when card is not in unicode map
	// We need to use an invalid combination that's not in the map
	invalidCard := poker.Card{Suit: poker.Suit(99), Rank: poker.Rank(99)}
	result := invalidCard.String()
	expected := fmt.Sprintf("%s-%s", "", "")
	if result != expected {
		t.Errorf("Expected fallback format %q, got %q", expected, result)
	}
}

func TestNewCard(t *testing.T) {
	testCases := []struct {
		name string
		suit poker.Suit
		rank poker.Rank
	}{
		{"Ace of Spades", poker.SuitSpade, poker.RankAce},
		{"King of Hearts", poker.SuitHeart, poker.RankKing},
		{"Two of Clubs", poker.SuitClub, poker.RankTwo},
		{"Queen of Diamonds", poker.SuitDiamond, poker.RankQueen},
		{"None card", poker.SuitNone, poker.RankNone},
		{"Joker", poker.SuitNone, poker.RankJoker},
		{"Colored Joker", poker.SuitNone, poker.RankColoredJoker},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card := poker.NewCard(tc.suit, tc.rank)

			if card == nil {
				t.Error("NewCard returned nil")
				return
			}

			if card.Suit != tc.suit {
				t.Errorf("Expected Suit: %d, got: %d", tc.suit, card.Suit)
			}

			if card.Rank != tc.rank {
				t.Errorf("Expected Rank: %d, got: %d", tc.rank, card.Rank)
			}
		})
	}
}

func TestSuitConstants(t *testing.T) {
	// Test all suit constants exist and have expected values
	suitTests := []struct {
		suit     poker.Suit
		expected uint8
	}{
		{poker.SuitNone, 0},
		{poker.SuitHeart, 1},
		{poker.SuitDiamond, 2},
		{poker.SuitClub, 3},
		{poker.SuitSpade, 4},
	}

	for _, test := range suitTests {
		if uint8(test.suit) != test.expected {
			t.Errorf("Expected suit %d to have value %d, got %d", test.suit, test.expected, uint8(test.suit))
		}
	}
}

func TestRankConstants(t *testing.T) {
	// Test all rank constants exist and have expected values
	rankTests := []struct {
		rank     poker.Rank
		expected uint8
	}{
		{poker.RankNone, 0},
		{poker.RankAce, 1},
		{poker.RankTwo, 2},
		{poker.RankThree, 3},
		{poker.RankFour, 4},
		{poker.RankFive, 5},
		{poker.RankSix, 6},
		{poker.RankSeven, 7},
		{poker.RankEight, 8},
		{poker.RankNine, 9},
		{poker.RankTen, 10},
		{poker.RankJack, 11},
		{poker.RankQueen, 12},
		{poker.RankKing, 13},
		{poker.RankJoker, 14},
		{poker.RankColoredJoker, 15},
	}

	for _, test := range rankTests {
		if uint8(test.rank) != test.expected {
			t.Errorf("Expected rank %d to have value %d, got %d", test.rank, test.expected, uint8(test.rank))
		}
	}
}

func TestRankMapAccess(t *testing.T) {
	// Test that RankMap can be accessed and contains expected entries
	allRanks := []poker.Rank{
		poker.RankNone, poker.RankAce, poker.RankTwo, poker.RankThree,
		poker.RankFour, poker.RankFive, poker.RankSix, poker.RankSeven,
		poker.RankEight, poker.RankNine, poker.RankTen, poker.RankJack,
		poker.RankQueen, poker.RankKing, poker.RankJoker, poker.RankColoredJoker,
	}

	for _, rank := range allRanks {
		if _, exists := poker.RankMap[rank]; !exists {
			t.Errorf("RankMap missing entry for rank %d", rank)
		}
	}

	// Test specific expected values
	expectedRankMappings := map[poker.Rank]string{
		poker.RankNone:         "",
		poker.RankAce:          "A",
		poker.RankTwo:          "2",
		poker.RankThree:        "3",
		poker.RankFour:         "4",
		poker.RankFive:         "5",
		poker.RankSix:          "6",
		poker.RankSeven:        "7",
		poker.RankEight:        "8",
		poker.RankNine:         "9",
		poker.RankTen:          "10",
		poker.RankJack:         "J",
		poker.RankQueen:        "Q",
		poker.RankKing:         "K",
		poker.RankJoker:        "Joker",
		poker.RankColoredJoker: "ColoredJoker",
	}

	for rank, expected := range expectedRankMappings {
		if actual, exists := poker.RankMap[rank]; !exists {
			t.Errorf("RankMap missing entry for rank %d", rank)
		} else if actual != expected {
			t.Errorf("Expected RankMap[%d] = %q, got %q", rank, expected, actual)
		}
	}
}

func TestCardCreationEdgeCases(t *testing.T) {
	// Test edge cases for card creation
	testCases := []struct {
		name string
		suit poker.Suit
		rank poker.Rank
	}{
		{"Max suit value", poker.Suit(255), poker.RankAce},
		{"Max rank value", poker.SuitSpade, poker.Rank(255)},
		{"Both max values", poker.Suit(255), poker.Rank(255)},
		{"Zero values", poker.Suit(0), poker.Rank(0)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			card := poker.NewCard(tc.suit, tc.rank)
			if card == nil {
				t.Error("NewCard should not return nil for any input")
			}
			if card.Suit != tc.suit || card.Rank != tc.rank {
				t.Errorf("Card fields not set correctly: expected %d/%d, got %d/%d",
					tc.suit, tc.rank, card.Suit, card.Rank)
			}
		})
	}
}
