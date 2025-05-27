package holdem

import "testing"

func TestNewGame(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)

	if game == nil {
		t.Fatal("Expected game to be created")
	}

	if len(game.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(game.Players))
	}

	if game.SmallBlind != 5 {
		t.Errorf("Expected small blind 5, got %d", game.SmallBlind)
	}

	if game.BigBlind != 10 {
		t.Errorf("Expected big blind 10, got %d", game.BigBlind)
	}

	// Check player initialization
	for i, player := range game.Players {
		if player.GetID() != i {
			t.Errorf("Expected player ID %d, got %d", i, player.GetID())
		}
		if player.GetName() != playerNames[i] {
			t.Errorf("Expected player name %s, got %s", playerNames[i], player.GetName())
		}
		if player.GetChips() != 1000 {
			t.Errorf("Expected 1000 chips, got %d", player.GetChips())
		}
	}
}

func TestStartHand(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)

	game.StartHand()

	// Check that blinds were posted
	if game.Pot != 15 { // 5 + 10
		t.Errorf("Expected pot 15, got %d", game.Pot)
	}

	// Check that players have hole cards
	for i, player := range game.Players {
		cards := player.GetHandCards()
		if len(cards) != 2 {
			t.Errorf("Player %d should have 2 cards, got %d", i, len(cards))
		}
	}

	// Check current bet is big blind
	if game.CurrentBet != 10 {
		t.Errorf("Expected current bet 10, got %d", game.CurrentBet)
	}
}

func TestPlayerActions(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	initialPot := game.Pot
	currentPlayer := game.GetCurrentPlayer()
	initialChips := currentPlayer.GetChips()

	// Test call
	game.Call()

	expectedCallAmount := 5 // Bob needs to call 5 more to match Alice's 10
	if currentPlayer.GetChips() != initialChips-expectedCallAmount {
		t.Errorf("Expected chips %d, got %d", initialChips-expectedCallAmount, currentPlayer.GetChips())
	}

	if game.Pot != initialPot+expectedCallAmount {
		t.Errorf("Expected pot %d, got %d", initialPot+expectedCallAmount, game.Pot)
	}
}

func TestFold(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	currentPlayer := game.GetCurrentPlayer()

	// Test fold
	game.Fold()

	if !currentPlayer.IsFolded() {
		t.Error("Player should be folded")
	}
}

func TestNextPhase(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	// Should start at preflop
	if game.CurrentPhase != PhasePreflop {
		t.Errorf("Expected preflop phase, got %d", game.CurrentPhase)
	}

	// Advance to flop
	game.NextPhase()
	if game.CurrentPhase != PhaseFlop {
		t.Errorf("Expected flop phase, got %d", game.CurrentPhase)
	}

	// Should have 3 community cards
	if len(game.CommunityCards) != 3 {
		t.Errorf("Expected 3 community cards, got %d", len(game.CommunityCards))
	}
}

func TestBettingRoundCompletion(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	// Alice needs to call or raise (current bet is 10, Alice has 5)
	handComplete := game.Call()
	if handComplete {
		t.Error("Hand should not be complete after first call")
	}

	// Bob can check (both have called the big blind)
	handComplete = game.Check()
	if !handComplete {
		t.Error("Betting round should be complete and should advance to flop")
	}

	// Should be at flop now
	if game.CurrentPhase != PhaseFlop {
		t.Errorf("Expected flop phase, got %d", game.CurrentPhase)
	}

	if len(game.CommunityCards) != 3 {
		t.Errorf("Expected 3 community cards, got %d", len(game.CommunityCards))
	}
}

func TestAutomaticPhaseProgression(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	// Complete preflop betting
	game.Call()  // Alice calls
	game.Check() // Bob checks

	// Should be at flop
	if game.CurrentPhase != PhaseFlop {
		t.Errorf("Expected flop phase, got %d", game.CurrentPhase)
	}

	// Complete flop betting
	game.Check() // First player checks
	game.Check() // Second player checks

	// Should be at turn
	if game.CurrentPhase != PhaseTurn {
		t.Errorf("Expected turn phase, got %d", game.CurrentPhase)
	}

	if len(game.CommunityCards) != 4 {
		t.Errorf("Expected 4 community cards, got %d", len(game.CommunityCards))
	}

	// Complete turn betting
	game.Check()
	game.Check()

	// Should be at river
	if game.CurrentPhase != PhaseRiver {
		t.Errorf("Expected river phase, got %d", game.CurrentPhase)
	}

	if len(game.CommunityCards) != 5 {
		t.Errorf("Expected 5 community cards, got %d", len(game.CommunityCards))
	}

	// Complete river betting
	game.Check()
	game.Check()

	// Should be at showdown
	if game.CurrentPhase != PhaseShowdown {
		t.Errorf("Expected showdown phase, got %d", game.CurrentPhase)
	}

	if !game.IsHandComplete() {
		t.Error("Hand should be complete")
	}
}

func TestEarlyHandEnd(t *testing.T) {
	playerNames := []string{"Alice", "Bob", "Charlie"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	// Two players fold, leaving only one
	game.Fold() // First player folds
	game.Fold() // Second player folds

	// Hand should end immediately
	if game.CurrentPhase != PhaseShowdown {
		t.Errorf("Expected showdown phase after all but one fold, got %d", game.CurrentPhase)
	}

	if !game.IsHandComplete() {
		t.Error("Hand should be complete when only one player remains")
	}

	// Remaining player should have won the pot
	activePlayers := game.GetActivePlayers()
	if len(activePlayers) != 1 {
		t.Errorf("Expected 1 active player, got %d", len(activePlayers))
	}
}

func TestShowdownWinnerDetermination(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	// Play through to showdown
	game.Call()
	game.Check()
	game.Check()
	game.Check()
	game.Check()
	game.Check()
	game.Check()
	game.Check()

	// Should be at showdown
	if game.CurrentPhase != PhaseShowdown {
		t.Errorf("Expected showdown phase, got %d", game.CurrentPhase)
	}

	// Should have winners
	winners := game.GetWinners()
	if len(winners) == 0 {
		t.Error("Should have at least one winner")
	}

	// All winners should be active players
	for _, winner := range winners {
		if winner.IsFolded() {
			t.Error("Winner should not be folded")
		}
	}
}

func TestHandResultEvaluation(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	// Play to river
	game.Call()
	game.Check()
	game.Check()
	game.Check()
	game.Check()
	game.Check()

	// Should be able to get hand results at river
	for _, player := range game.Players {
		if !player.IsFolded() {
			result := game.GetHandResult(player)
			if result == nil {
				t.Errorf("Should have hand result for player %s", player.GetName())
			}
			if result.Rank < HighCard || result.Rank > RoyalFlush {
				t.Errorf("Invalid hand rank for player %s", player.GetName())
			}
		}
	}
}

func TestRaiseAction(t *testing.T) {
	playerNames := []string{"Alice", "Bob"}
	game := NewGame(playerNames, 5, 10)
	game.StartHand()

	initialPot := game.Pot
	currentPlayer := game.GetCurrentPlayer()
	initialChips := currentPlayer.GetChips()

	// Test raise
	raiseAmount := 10
	game.Raise(raiseAmount)

	// Current bet should be updated
	expectedBet := 20 // 10 (big blind) + 10 (raise)
	if game.CurrentBet != expectedBet {
		t.Errorf("Expected current bet %d, got %d", expectedBet, game.CurrentBet)
	}

	// Player chips should be reduced
	expectedChips := initialChips - (expectedBet - 5) // Player already had 5 in
	if currentPlayer.GetChips() != expectedChips {
		t.Errorf("Expected chips %d, got %d", expectedChips, currentPlayer.GetChips())
	}

	// Pot should be increased
	expectedPot := initialPot + (expectedBet - 5)
	if game.Pot != expectedPot {
		t.Errorf("Expected pot %d, got %d", expectedPot, game.Pot)
	}
}
