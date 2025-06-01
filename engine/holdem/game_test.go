package holdem

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame([]string{"Alice", "Bob"}, 5, 10)

	if len(game.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(game.Players))
	}

	if game.SmallBlind != 5 {
		t.Errorf("Expected small blind 5, got %d", game.SmallBlind)
	}

	if game.BigBlind != 10 {
		t.Errorf("Expected big blind 10, got %d", game.BigBlind)
	}

	// Check players are created correctly
	for i, player := range game.Players {
		if player.GetID() != i {
			t.Errorf("Expected player ID %d, got %d", i, player.GetID())
		}
		expectedNames := []string{"Alice", "Bob"}
		if player.GetName() != expectedNames[i] {
			t.Errorf("Expected player name %s, got %s", expectedNames[i], player.GetName())
		}
		if player.GetChips() != 1000 {
			t.Errorf("Expected player chips 1000, got %d", player.GetChips())
		}
	}
}

func TestStartHand(t *testing.T) {
	game := NewGame([]string{"Alice", "Bob"}, 5, 10)
	game.StartHand()

	// Check that each player has 2 cards
	for i, player := range game.Players {
		cards := player.GetHandCards()
		if len(cards) != 2 {
			t.Errorf("Player %d should have 2 cards, got %d", i, len(cards))
		}
	}

	// Check that pot has blinds
	expectedPot := game.SmallBlind + game.BigBlind
	if game.Pot != expectedPot {
		t.Errorf("Expected pot %d, got %d", expectedPot, game.Pot)
	}

	// Check current bet is big blind
	if game.CurrentBet != game.BigBlind {
		t.Errorf("Expected current bet %d, got %d", game.BigBlind, game.CurrentBet)
	}
}

func TestGetWinners(t *testing.T) {
	game := NewGame([]string{"Alice", "Bob"}, 5, 10)
	game.StartHand()

	// Set game to showdown phase for testing
	game.CurrentPhase = PhaseShowdown

	winners := game.GetWinners()
	if len(winners) == 0 {
		t.Error("Expected at least one winner")
	}
}

func TestGameUtilities(t *testing.T) {
	game := NewGame([]string{"Alice", "Bob", "Charlie"}, 5, 10)

	// Test GetTotalPlayers
	if game.GetTotalPlayers() != 3 {
		t.Errorf("Expected 3 total players, got %d", game.GetTotalPlayers())
	}

	// Test GetNumActivePlayers (all active initially)
	if game.GetNumActivePlayers() != 3 {
		t.Errorf("Expected 3 active players, got %d", game.GetNumActivePlayers())
	}

	// Test GetPlayerPosition
	player := game.Players[1]
	position := game.GetPlayerPosition(player)
	if position != 1 {
		t.Errorf("Expected position 1, got %d", position)
	}

	// Test IsPlayerTurn
	game.CurrentPlayerPosition = 0
	if !game.IsPlayerTurn(game.Players[0]) {
		t.Error("Expected player 0 to be current player")
	}

	// Test CalculatePotOdds
	game.Pot = 100
	game.CurrentBet = 20
	game.Players[0].Bet(0) // Player hasn't bet yet
	potOdds := game.CalculatePotOdds(game.Players[0])
	expected := 20.0 / 120.0
	if potOdds != expected {
		t.Errorf("Expected pot odds %.3f, got %.3f", expected, potOdds)
	}
}
