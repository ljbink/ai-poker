package holdem

import (
	"testing"

	"github.com/ljbink/ai-poker/engine/poker"
)

func TestNewPlayer(t *testing.T) {
	id := 1
	name := "Test Player"
	chips := 1000

	player := NewPlayer(id, name, chips)

	if player == nil {
		t.Fatal("NewPlayer returned nil")
	}

	if player.GetID() != id {
		t.Errorf("Expected player ID %d, got %d", id, player.GetID())
	}

	if player.GetName() != name {
		t.Errorf("Expected player name %s, got %s", name, player.GetName())
	}

	if player.GetChips() != chips {
		t.Errorf("Expected player chips %d, got %d", chips, player.GetChips())
	}

	if player.GetBet() != 0 {
		t.Errorf("Expected initial bet 0, got %d", player.GetBet())
	}

	if player.GetTotalBet() != 0 {
		t.Errorf("Expected initial total bet 0, got %d", player.GetTotalBet())
	}

	if player.IsFolded() {
		t.Error("Expected player not to be folded initially")
	}

	if len(player.GetHandCards()) != 0 {
		t.Errorf("Expected no hand cards initially, got %d", len(player.GetHandCards()))
	}
}

func TestPlayerDealCard(t *testing.T) {
	player := NewPlayer(1, "Test Player", 1000)

	// Deal first card
	card1 := &poker.Card{Rank: poker.RankAce, Suit: poker.SuitSpade}
	result := player.DealCard(card1)

	// Test method chaining
	if result != player {
		t.Error("DealCard should return the player for method chaining")
	}

	cards := player.GetHandCards()
	if len(cards) != 1 {
		t.Errorf("Expected 1 card after dealing, got %d", len(cards))
	}
	if cards[0] != card1 {
		t.Error("Dealt card does not match expected card")
	}

	// Deal second card
	card2 := &poker.Card{Rank: poker.RankKing, Suit: poker.SuitHeart}
	player.DealCard(card2)

	cards = player.GetHandCards()
	if len(cards) != 2 {
		t.Errorf("Expected 2 cards after dealing second card, got %d", len(cards))
	}
	if cards[1] != card2 {
		t.Error("Second dealt card does not match expected card")
	}

	// Test dealing nil card (should still append)
	player.DealCard(nil)
	cards = player.GetHandCards()
	if len(cards) != 3 {
		t.Errorf("Expected 3 cards after dealing nil card, got %d", len(cards))
	}
}

func TestPlayerChips(t *testing.T) {
	initialChips := 1000
	player := NewPlayer(1, "Test Player", initialChips)

	// Test initial chips
	if player.GetChips() != initialChips {
		t.Errorf("Expected initial chips %d, got %d", initialChips, player.GetChips())
	}

	// Test granting chips
	grantAmount := 500
	result := player.GrandChips(grantAmount)

	// Test method chaining
	if result != player {
		t.Error("GrandChips should return the player for method chaining")
	}

	expectedChips := initialChips + grantAmount
	if player.GetChips() != expectedChips {
		t.Errorf("Expected chips %d after granting %d, got %d", expectedChips, grantAmount, player.GetChips())
	}

	// Test granting negative chips (should decrease)
	player.GrandChips(-200)
	expectedChips = expectedChips - 200
	if player.GetChips() != expectedChips {
		t.Errorf("Expected chips %d after granting -200, got %d", expectedChips, player.GetChips())
	}
}

func TestPlayerBetting(t *testing.T) {
	initialChips := 1000
	player := NewPlayer(1, "Test Player", initialChips)

	// Test initial betting state
	if player.GetBet() != 0 {
		t.Errorf("Expected initial bet 0, got %d", player.GetBet())
	}
	if player.GetTotalBet() != 0 {
		t.Errorf("Expected initial total bet 0, got %d", player.GetTotalBet())
	}

	// Test first bet
	betAmount := 100
	result := player.Bet(betAmount)

	// Test method chaining
	if result != player {
		t.Error("Bet should return the player for method chaining")
	}

	if player.GetBet() != betAmount {
		t.Errorf("Expected bet %d, got %d", betAmount, player.GetBet())
	}
	if player.GetTotalBet() != betAmount {
		t.Errorf("Expected total bet %d, got %d", betAmount, player.GetTotalBet())
	}
	if player.GetChips() != initialChips-betAmount {
		t.Errorf("Expected chips %d after betting %d, got %d", initialChips-betAmount, betAmount, player.GetChips())
	}

	// Test additional bet
	additionalBet := 50
	player.Bet(additionalBet)

	expectedBet := betAmount + additionalBet
	expectedTotalBet := betAmount + additionalBet
	expectedChips := initialChips - expectedTotalBet

	if player.GetBet() != expectedBet {
		t.Errorf("Expected bet %d after additional bet, got %d", expectedBet, player.GetBet())
	}
	if player.GetTotalBet() != expectedTotalBet {
		t.Errorf("Expected total bet %d after additional bet, got %d", expectedTotalBet, player.GetTotalBet())
	}
	if player.GetChips() != expectedChips {
		t.Errorf("Expected chips %d after additional bet, got %d", expectedChips, player.GetChips())
	}

	// Test reset bet
	resetResult := player.ResetBet()

	// Test method chaining
	if resetResult != player {
		t.Error("ResetBet should return the player for method chaining")
	}

	if player.GetBet() != 0 {
		t.Errorf("Expected bet 0 after reset, got %d", player.GetBet())
	}
	if player.GetTotalBet() != expectedTotalBet {
		t.Errorf("Expected total bet %d unchanged after reset, got %d", expectedTotalBet, player.GetTotalBet())
	}
	if player.GetChips() != expectedChips {
		t.Errorf("Expected chips %d unchanged after reset, got %d", expectedChips, player.GetChips())
	}
}

func TestPlayerFolding(t *testing.T) {
	player := NewPlayer(1, "Test Player", 1000)

	// Test initial folded state
	if player.IsFolded() {
		t.Error("Expected player not to be folded initially")
	}

	// Test folding
	result := player.Fold()

	// Test method chaining
	if result != player {
		t.Error("Fold should return the player for method chaining")
	}

	if !player.IsFolded() {
		t.Error("Expected player to be folded after calling Fold()")
	}
}

func TestPlayerResetForNewHand(t *testing.T) {
	player := NewPlayer(1, "Test Player", 1000)

	// Setup player state
	card1 := &poker.Card{Rank: poker.RankAce, Suit: poker.SuitSpade}
	card2 := &poker.Card{Rank: poker.RankKing, Suit: poker.SuitHeart}
	player.DealCard(card1)
	player.DealCard(card2)
	player.Bet(100)
	player.Fold()

	// Verify state before reset
	if len(player.GetHandCards()) != 2 {
		t.Error("Expected player to have cards before reset")
	}
	if player.GetBet() != 100 {
		t.Error("Expected player to have bet before reset")
	}
	if player.GetTotalBet() != 100 {
		t.Error("Expected player to have total bet before reset")
	}
	if !player.IsFolded() {
		t.Error("Expected player to be folded before reset")
	}

	// Test reset
	result := player.ResetForNewHand()

	// Test method chaining
	if result != player {
		t.Error("ResetForNewHand should return the player for method chaining")
	}

	// Verify state after reset
	if player.GetHandCards() != nil && len(player.GetHandCards()) != 0 {
		t.Errorf("Expected no hand cards after reset, got %d", len(player.GetHandCards()))
	}
	if player.GetBet() != 0 {
		t.Errorf("Expected bet 0 after reset, got %d", player.GetBet())
	}
	if player.GetTotalBet() != 0 {
		t.Errorf("Expected total bet 0 after reset, got %d", player.GetTotalBet())
	}
	if player.IsFolded() {
		t.Error("Expected player not to be folded after reset")
	}

	// Verify chips are unchanged
	if player.GetChips() != 900 { // 1000 - 100 from betting
		t.Errorf("Expected chips to remain at 900 after reset, got %d", player.GetChips())
	}
}

func TestPlayerMethodChaining(t *testing.T) {
	player := NewPlayer(1, "Test Player", 1000)

	// Test multiple operations chained together
	card := &poker.Card{Rank: poker.RankAce, Suit: poker.SuitSpade}
	result := player.DealCard(card).
		GrandChips(500).
		Bet(200).
		ResetBet().
		Fold().
		ResetForNewHand()

	if result != player {
		t.Error("Method chaining should return the same player instance")
	}

	// Verify final state
	if player.GetChips() != 1300 { // 1000 + 500 - 200
		t.Errorf("Expected chips 1300 after chained operations, got %d", player.GetChips())
	}
	if player.GetBet() != 0 {
		t.Errorf("Expected bet 0 after reset in chain, got %d", player.GetBet())
	}
	if player.GetTotalBet() != 0 {
		t.Errorf("Expected total bet 0 after reset in chain, got %d", player.GetTotalBet())
	}
	if player.IsFolded() {
		t.Error("Expected player not folded after reset in chain")
	}
}

func TestPlayerBettingEdgeCases(t *testing.T) {
	player := NewPlayer(1, "Test Player", 100)

	// Test betting 0
	player.Bet(0)
	if player.GetBet() != 0 {
		t.Errorf("Expected bet 0 after betting 0, got %d", player.GetBet())
	}
	if player.GetChips() != 100 {
		t.Errorf("Expected chips unchanged after betting 0, got %d", player.GetChips())
	}

	// Test betting negative amount (should still subtract from chips)
	player.Bet(-50)
	if player.GetBet() != -50 {
		t.Errorf("Expected bet -50 after betting negative, got %d", player.GetBet())
	}
	if player.GetChips() != 150 { // 100 - (-50) = 150
		t.Errorf("Expected chips 150 after betting -50, got %d", player.GetChips())
	}

	// Reset for next test
	player.ResetForNewHand()
	// Reset chips by granting the difference
	player.GrandChips(100 - player.GetChips())

	// Test betting more than available chips
	player.Bet(200)
	if player.GetBet() != 200 {
		t.Errorf("Expected bet 200 even if exceeding chips, got %d", player.GetBet())
	}
	if player.GetChips() != -100 {
		t.Errorf("Expected chips -100 after betting more than available, got %d", player.GetChips())
	}
}
