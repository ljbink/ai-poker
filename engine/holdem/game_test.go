package holdem

import (
	"fmt"
	"testing"

	"github.com/ljbink/ai-poker/engine/poker"
)

func TestNewGame(t *testing.T) {
	smallBlind := 10
	bigBlind := 20
	game := NewGame(smallBlind, bigBlind)

	if game == nil {
		t.Fatal("NewGame returned nil")
	}

	if game.GetSmallBlind() != smallBlind {
		t.Errorf("Expected small blind %d, got %d", smallBlind, game.GetSmallBlind())
	}

	if game.GetBigBlind() != bigBlind {
		t.Errorf("Expected big blind %d, got %d", bigBlind, game.GetBigBlind())
	}

	if game.GetCurrentPhase() != PhasePreflop {
		t.Errorf("Expected initial phase to be PhasePreflop, got %d", game.GetCurrentPhase())
	}

	if len(game.GetCommunityCards()) != 0 {
		t.Errorf("Expected no community cards initially, got %d cards", len(game.GetCommunityCards()))
	}

	// Check action logs are initialized
	userActions := game.GetUserActions()
	if len(userActions.Preflop) != 0 || len(userActions.Flop) != 0 || len(userActions.Turn) != 0 || len(userActions.River) != 0 {
		t.Error("Expected empty user action logs initially")
	}
}

func TestPlayerSit(t *testing.T) {
	game := NewGame(10, 20)
	player := NewPlayer(1, "Test Player", 1000)

	// Test successful sitting
	err := game.PlayerSit(player, 0)
	if err != nil {
		t.Errorf("Unexpected error sitting player: %v", err)
	}

	// Verify player is seated
	seatedPlayer, err := game.GetPlayerBySit(0)
	if err != nil {
		t.Errorf("Error getting player by sit: %v", err)
	}
	if seatedPlayer.GetID() != player.GetID() {
		t.Errorf("Expected player ID %d, got %d", player.GetID(), seatedPlayer.GetID())
	}

	// Test nil player
	err = game.PlayerSit(nil, 1)
	if err == nil {
		t.Error("Expected error when sitting nil player")
	}

	// Test invalid sit number
	err = game.PlayerSit(player, -1)
	if err == nil {
		t.Error("Expected error for negative sit number")
	}

	err = game.PlayerSit(player, 10)
	if err == nil {
		t.Error("Expected error for sit number >= 10")
	}

	// Test sitting another player in occupied seat
	player2 := NewPlayer(2, "Test Player 2", 1000)
	err = game.PlayerSit(player2, 0)
	if err == nil {
		t.Error("Expected error when trying to sit in occupied seat")
	}

	// Test sitting same player in same seat (should succeed)
	err = game.PlayerSit(player, 0)
	if err != nil {
		t.Errorf("Unexpected error when sitting same player in same seat: %v", err)
	}
}

func TestPlayerLeave(t *testing.T) {
	game := NewGame(10, 20)
	player := NewPlayer(1, "Test Player", 1000)

	// Sit player first
	game.PlayerSit(player, 0)

	// Test successful leaving
	err := game.PlayerLeave(player)
	if err != nil {
		t.Errorf("Unexpected error when player leaves: %v", err)
	}

	// Verify player is no longer seated
	_, err = game.GetPlayerBySit(0)
	if err == nil {
		t.Error("Expected error when getting player from empty seat")
	}

	// Test leaving nil player
	err = game.PlayerLeave(nil)
	if err == nil {
		t.Error("Expected error when leaving nil player")
	}

	// Test leaving player not in game (should not error)
	player2 := NewPlayer(2, "Test Player 2", 1000)
	err = game.PlayerLeave(player2)
	if err != nil {
		t.Errorf("Unexpected error when leaving player not in game: %v", err)
	}
}

func TestGetPlayerByID(t *testing.T) {
	game := NewGame(10, 20)
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 1000)

	// Test getting player not in game
	_, err := game.GetPlayerByID(1)
	if err == nil {
		t.Error("Expected error when getting player not in game")
	}

	// Sit players
	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 5)

	// Test getting existing players
	foundPlayer, err := game.GetPlayerByID(1)
	if err != nil {
		t.Errorf("Unexpected error getting player by ID: %v", err)
	}
	if foundPlayer.GetID() != 1 {
		t.Errorf("Expected player ID 1, got %d", foundPlayer.GetID())
	}

	foundPlayer, err = game.GetPlayerByID(2)
	if err != nil {
		t.Errorf("Unexpected error getting player by ID: %v", err)
	}
	if foundPlayer.GetID() != 2 {
		t.Errorf("Expected player ID 2, got %d", foundPlayer.GetID())
	}

	// Test getting non-existent player
	_, err = game.GetPlayerByID(99)
	if err == nil {
		t.Error("Expected error when getting non-existent player")
	}
}

func TestGetPlayerBySit(t *testing.T) {
	game := NewGame(10, 20)
	player := NewPlayer(1, "Test Player", 1000)

	// Test getting from empty seat
	_, err := game.GetPlayerBySit(0)
	if err == nil {
		t.Error("Expected error when getting player from empty seat")
	}

	// Test invalid sit numbers
	_, err = game.GetPlayerBySit(-1)
	if err == nil {
		t.Error("Expected error for negative sit number")
	}

	_, err = game.GetPlayerBySit(10)
	if err == nil {
		t.Error("Expected error for sit number >= 10")
	}

	// Sit player and test successful retrieval
	game.PlayerSit(player, 3)
	foundPlayer, err := game.GetPlayerBySit(3)
	if err != nil {
		t.Errorf("Unexpected error getting player by sit: %v", err)
	}
	if foundPlayer.GetID() != player.GetID() {
		t.Errorf("Expected player ID %d, got %d", player.GetID(), foundPlayer.GetID())
	}
}

func TestGetPlayerSitByID(t *testing.T) {
	game := NewGame(10, 20)
	player := NewPlayer(1, "Test Player", 1000)

	// Test getting sit for player not in game
	_, err := game.GetPlayerSitByID(1)
	if err == nil {
		t.Error("Expected error when getting sit for player not in game")
	}

	// Sit player and test successful retrieval
	game.PlayerSit(player, 7)
	sit, err := game.GetPlayerSitByID(1)
	if err != nil {
		t.Errorf("Unexpected error getting player sit by ID: %v", err)
	}
	if sit != 7 {
		t.Errorf("Expected sit 7, got %d", sit)
	}
}

func TestGetCurrentPlayer(t *testing.T) {
	game := NewGame(10, 20)

	// Test with no players
	currentPlayer := game.GetCurrentPlayer()
	if currentPlayer != nil {
		t.Error("Expected nil current player when no players in game")
	}

	// Add players
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 1000)
	player3 := NewPlayer(3, "Player 3", 1000)

	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 3)
	game.PlayerSit(player3, 8)

	// Test getting first active player
	currentPlayer = game.GetCurrentPlayer()
	if currentPlayer == nil {
		t.Fatal("Expected non-nil current player")
	}
	if currentPlayer.GetID() != 1 {
		t.Errorf("Expected current player ID 1, got %d", currentPlayer.GetID())
	}

	// Fold first player and test getting next active player
	player1.Fold()
	currentPlayer = game.GetCurrentPlayer()
	if currentPlayer == nil {
		t.Fatal("Expected non-nil current player after first player folds")
	}
	if currentPlayer.GetID() != 2 {
		t.Errorf("Expected current player ID 2, got %d", currentPlayer.GetID())
	}

	// Fold all players
	player2.Fold()
	player3.Fold()
	currentPlayer = game.GetCurrentPlayer()
	if currentPlayer != nil {
		t.Error("Expected nil current player when all players are folded")
	}
}

func TestTakeAction(t *testing.T) {
	game := NewGame(10, 20)

	// Test action in preflop phase
	action := Action{PlayerID: 1, Type: ActionCall, Amount: 20}
	err := game.TakeAction(action)
	if err != nil {
		t.Errorf("Unexpected error taking action in preflop: %v", err)
	}

	userActions := game.GetUserActions()
	if len(userActions.Preflop) != 1 {
		t.Errorf("Expected 1 preflop action, got %d", len(userActions.Preflop))
	}
	if userActions.Preflop[0].PlayerID != 1 {
		t.Errorf("Expected action from player 1, got player %d", userActions.Preflop[0].PlayerID)
	}

	// Test action in different phase
	game.SetCurrentPhase(PhaseFlop)
	flopAction := Action{
		PlayerID: 2,
		Type:     ActionCall,
		Amount:   10,
	}

	err = game.TakeAction(flopAction)
	if err != nil {
		t.Errorf("Unexpected error taking flop action: %v", err)
	}

	userActions = game.GetUserActions()
	if len(userActions.Flop) != 1 {
		t.Errorf("Expected 1 flop action, got %d", len(userActions.Flop))
	}

	game.currentPhase = PhaseTurn
	action = Action{PlayerID: 3, Type: ActionCheck, Amount: 0}
	err = game.TakeAction(action)
	if err != nil {
		t.Errorf("Unexpected error taking action in turn: %v", err)
	}

	game.currentPhase = PhaseRiver
	action = Action{PlayerID: 4, Type: ActionFold, Amount: 0}
	err = game.TakeAction(action)
	if err != nil {
		t.Errorf("Unexpected error taking action in river: %v", err)
	}

	// Test invalid phase
	game.currentPhase = GamePhase(99)
	action = Action{PlayerID: 5, Type: ActionCall, Amount: 20}
	err = game.TakeAction(action)
	if err == nil {
		t.Error("Expected error for invalid game phase")
	}
}

func TestGameCommunityCards(t *testing.T) {
	game := NewGame(10, 20)

	// Test initial empty community cards
	cards := game.GetCommunityCards()
	if len(cards) != 0 {
		t.Errorf("Expected 0 community cards initially, got %d", len(cards))
	}

	// Add community cards manually for testing
	testCard1 := &poker.Card{Rank: poker.RankAce, Suit: poker.SuitSpade}
	testCard2 := &poker.Card{Rank: poker.RankKing, Suit: poker.SuitHeart}
	game.communityCards = poker.Cards{testCard1, testCard2}

	cards = game.GetCommunityCards()
	if len(cards) != 2 {
		t.Errorf("Expected 2 community cards, got %d", len(cards))
	}
	if cards[0].Rank != poker.RankAce || cards[0].Suit != poker.SuitSpade {
		t.Errorf("Expected Ace of Spades, got %s", cards[0].String())
	}
	if cards[1].Rank != poker.RankKing || cards[1].Suit != poker.SuitHeart {
		t.Errorf("Expected King of Hearts, got %s", cards[1].String())
	}
}

func TestGamePhaseProgression(t *testing.T) {
	game := NewGame(10, 20)

	// Test initial phase
	if game.GetCurrentPhase() != PhasePreflop {
		t.Errorf("Expected initial phase PhasePreflop, got %d", game.GetCurrentPhase())
	}

	// Test phase changes
	game.currentPhase = PhaseFlop
	if game.GetCurrentPhase() != PhaseFlop {
		t.Errorf("Expected phase PhaseFlop, got %d", game.GetCurrentPhase())
	}

	game.currentPhase = PhaseTurn
	if game.GetCurrentPhase() != PhaseTurn {
		t.Errorf("Expected phase PhaseTurn, got %d", game.GetCurrentPhase())
	}

	game.currentPhase = PhaseRiver
	if game.GetCurrentPhase() != PhaseRiver {
		t.Errorf("Expected phase PhaseRiver, got %d", game.GetCurrentPhase())
	}

	game.currentPhase = PhaseShowdown
	if game.GetCurrentPhase() != PhaseShowdown {
		t.Errorf("Expected phase PhaseShowdown, got %d", game.GetCurrentPhase())
	}
}

func TestShuffleDeck(t *testing.T) {
	game := NewGame(5, 10)

	// Get initial deck order
	initialDeck := make([]poker.Card, len(game.deck))
	for i, card := range game.deck {
		initialDeck[i] = *card
	}

	// Shuffle deck
	game.ShuffleDeck()

	// Check deck has same cards but potentially different order
	if len(game.deck) != 52 {
		t.Errorf("Expected 52 cards, got %d", len(game.deck))
	}

	// Count different positions (shuffle should change order most of the time)
	differentPositions := 0
	for i, card := range game.deck {
		if i < len(initialDeck) && (*card != initialDeck[i]) {
			differentPositions++
		}
	}

	// Expect at least some cards to be in different positions
	if differentPositions < 10 {
		t.Errorf("Shuffle should change at least 10 card positions, only changed %d", differentPositions)
	}
}

func TestDealHoleCards(t *testing.T) {
	game := NewGame(5, 10)

	// Add players
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 1000)
	player3 := NewPlayer(3, "Player 3", 1000)

	err := game.PlayerSit(player1, 0)
	if err != nil {
		t.Errorf("Error seating player1: %v", err)
	}
	err = game.PlayerSit(player2, 1)
	if err != nil {
		t.Errorf("Error seating player2: %v", err)
	}
	err = game.PlayerSit(player3, 2)
	if err != nil {
		t.Errorf("Error seating player3: %v", err)
	}

	// Deal hole cards
	err = game.DealHoleCards()
	if err != nil {
		t.Errorf("Error dealing hole cards: %v", err)
	}

	// Check each player has 2 cards
	if len(player1.GetHandCards()) != 2 {
		t.Errorf("Player1 expected 2 cards, got %d", len(player1.GetHandCards()))
	}
	if len(player2.GetHandCards()) != 2 {
		t.Errorf("Player2 expected 2 cards, got %d", len(player2.GetHandCards()))
	}
	if len(player3.GetHandCards()) != 2 {
		t.Errorf("Player3 expected 2 cards, got %d", len(player3.GetHandCards()))
	}

	// Check deck has 46 cards remaining (52 - 6 dealt)
	if len(game.deck) != 46 {
		t.Errorf("Expected 46 cards in deck, got %d", len(game.deck))
	}

	// Check all dealt cards are different
	allCards := make(map[string]bool)
	for _, player := range []IPlayer{player1, player2, player3} {
		for _, card := range player.GetHandCards() {
			cardKey := fmt.Sprintf("%d-%d", card.Suit, card.Rank)
			if allCards[cardKey] {
				t.Errorf("Card %s dealt twice", cardKey)
			}
			allCards[cardKey] = true
		}
	}
}

func TestDealHoleCardsInsufficientPlayers(t *testing.T) {
	game := NewGame(5, 10)

	// Add only one player
	player1 := NewPlayer(1, "Player 1", 1000)
	err := game.PlayerSit(player1, 0)
	if err != nil {
		t.Errorf("Error seating player: %v", err)
	}

	// Try to deal hole cards
	err = game.DealHoleCards()
	if err == nil {
		t.Error("Expected error for insufficient players, got none")
	}
	if err != nil && !contains(err.Error(), "need at least 2 players") {
		t.Errorf("Expected 'need at least 2 players' error, got %v", err)
	}
}

func TestDealFlop(t *testing.T) {
	game := NewGame(5, 10)

	// Deal flop
	err := game.DealFlop()
	if err != nil {
		t.Errorf("Error dealing flop: %v", err)
	}

	// Check phase changed
	if game.GetCurrentPhase() != PhaseFlop {
		t.Errorf("Expected PhaseFlop, got %v", game.GetCurrentPhase())
	}

	// Check 3 community cards
	communityCards := game.GetCommunityCards()
	if len(communityCards) != 3 {
		t.Errorf("Expected 3 community cards, got %d", len(communityCards))
	}

	// Check deck has 48 cards remaining (52 - 1 burn - 3 community)
	if len(game.deck) != 48 {
		t.Errorf("Expected 48 cards in deck, got %d", len(game.deck))
	}
}

func TestDealTurn(t *testing.T) {
	game := NewGame(5, 10)

	// Deal flop first
	err := game.DealFlop()
	if err != nil {
		t.Errorf("Error dealing flop: %v", err)
	}

	// Deal turn
	err = game.DealTurn()
	if err != nil {
		t.Errorf("Error dealing turn: %v", err)
	}

	// Check phase changed
	if game.GetCurrentPhase() != PhaseTurn {
		t.Errorf("Expected PhaseTurn, got %v", game.GetCurrentPhase())
	}

	// Check 4 community cards
	communityCards := game.GetCommunityCards()
	if len(communityCards) != 4 {
		t.Errorf("Expected 4 community cards, got %d", len(communityCards))
	}

	// Check deck has 46 cards remaining (52 - 2 burns - 4 community)
	if len(game.deck) != 46 {
		t.Errorf("Expected 46 cards in deck, got %d", len(game.deck))
	}
}

func TestDealRiver(t *testing.T) {
	game := NewGame(5, 10)

	// Deal flop and turn first
	err := game.DealFlop()
	if err != nil {
		t.Errorf("Error dealing flop: %v", err)
	}
	err = game.DealTurn()
	if err != nil {
		t.Errorf("Error dealing turn: %v", err)
	}

	// Deal river
	err = game.DealRiver()
	if err != nil {
		t.Errorf("Error dealing river: %v", err)
	}

	// Check phase changed
	if game.GetCurrentPhase() != PhaseRiver {
		t.Errorf("Expected PhaseRiver, got %v", game.GetCurrentPhase())
	}

	// Check 5 community cards
	communityCards := game.GetCommunityCards()
	if len(communityCards) != 5 {
		t.Errorf("Expected 5 community cards, got %d", len(communityCards))
	}

	// Check deck has 44 cards remaining (52 - 3 burns - 5 community)
	if len(game.deck) != 44 {
		t.Errorf("Expected 44 cards in deck, got %d", len(game.deck))
	}
}

func TestSetCurrentPhase(t *testing.T) {
	game := NewGame(5, 10)

	// Initial phase should be preflop
	if game.GetCurrentPhase() != PhasePreflop {
		t.Errorf("Expected PhasePreflop initially, got %v", game.GetCurrentPhase())
	}

	// Set to flop
	game.SetCurrentPhase(PhaseFlop)
	if game.GetCurrentPhase() != PhaseFlop {
		t.Errorf("Expected PhaseFlop, got %v", game.GetCurrentPhase())
	}

	// Set to turn
	game.SetCurrentPhase(PhaseTurn)
	if game.GetCurrentPhase() != PhaseTurn {
		t.Errorf("Expected PhaseTurn, got %v", game.GetCurrentPhase())
	}

	// Set to river
	game.SetCurrentPhase(PhaseRiver)
	if game.GetCurrentPhase() != PhaseRiver {
		t.Errorf("Expected PhaseRiver, got %v", game.GetCurrentPhase())
	}
}

func TestGetAllPlayers(t *testing.T) {
	game := NewGame(5, 10)

	// Initially no players
	players := game.GetAllPlayers()
	if len(players) != 0 {
		t.Errorf("Expected 0 players initially, got %d", len(players))
	}

	// Add players
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 1000)

	err := game.PlayerSit(player1, 0)
	if err != nil {
		t.Errorf("Error seating player1: %v", err)
	}
	err = game.PlayerSit(player2, 3)
	if err != nil {
		t.Errorf("Error seating player2: %v", err)
	}

	// Check all players returned
	players = game.GetAllPlayers()
	if len(players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(players))
	}

	// Check players are correct
	foundPlayer1, foundPlayer2 := false, false
	for _, player := range players {
		if player.GetID() == 1 {
			foundPlayer1 = true
		}
		if player.GetID() == 2 {
			foundPlayer2 = true
		}
	}
	if !foundPlayer1 {
		t.Error("Player 1 not found in returned players")
	}
	if !foundPlayer2 {
		t.Error("Player 2 not found in returned players")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
