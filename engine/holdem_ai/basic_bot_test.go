package holdem_ai

import (
	"testing"

	"github.com/ljbink/ai-poker/engine/holdem"
	"github.com/ljbink/ai-poker/engine/poker"
)

func TestBasicBot(t *testing.T) {
	// Create a bot
	bot := NewBasicBot("TestBot")

	// Verify it implements the interface correctly
	if bot.GetName() != "TestBot" {
		t.Errorf("Expected name 'TestBot', got '%s'", bot.GetName())
	}

	if !bot.IsBot() {
		t.Error("Expected IsBot() to return true")
	}
}

func TestActionValidator(t *testing.T) {
	validator := NewActionValidator()

	// Create a simple game and player for testing
	game := holdem.NewGame([]string{"Player1", "Bot1"}, 5, 10)
	player := game.Players[0] // Get first player

	// Test fold action (should always be valid)
	foldAction := Action{Type: ActionFold}
	if !validator.IsValidAction(foldAction, game, player) {
		t.Error("Fold action should always be valid")
	}

	// Test check action when no bet to call
	checkAction := Action{Type: ActionCheck}
	if !validator.IsValidAction(checkAction, game, player) {
		t.Error("Check action should be valid when no bet to call")
	}
}

func TestGameAdapter(t *testing.T) {
	adapter := NewGameAdapter()

	// Create a simple game for testing
	game := holdem.NewGame([]string{"Player1", "Bot1"}, 5, 10)

	if game == nil {
		t.Fatal("Failed to create game")
	}

	// Test that adapter can be created
	if adapter == nil {
		t.Fatal("Failed to create game adapter")
	}
}

func TestPotOddsCalculation(t *testing.T) {
	game := holdem.NewGame([]string{"Player1", "Bot1"}, 5, 10)
	player := game.Players[0]

	// Set up a scenario with pot and bet
	game.Pot = 100
	game.CurrentBet = 20
	player.Bet(0) // Player hasn't bet yet

	potOdds := game.CalculatePotOdds(player)
	expected := 20.0 / 120.0 // 20 / (100 + 20)

	if potOdds != expected {
		t.Errorf("Expected pot odds %.3f, got %.3f", expected, potOdds)
	}

	// Test edge case - no bet to call
	game.CurrentBet = 0
	potOdds = game.CalculatePotOdds(player)
	if potOdds != 0.0 {
		t.Errorf("Expected pot odds 0.0 when no bet to call, got %.3f", potOdds)
	}
}

func TestHandStrengthEvaluation(t *testing.T) {
	bot := NewBasicBot("TestBot").(*BasicBot)
	game := holdem.NewGame([]string{"TestPlayer"}, 5, 10)
	player := game.Players[0]

	// Test pocket aces - manually set cards
	aceSpade := poker.NewCard(poker.SuitSpade, poker.RankAce)
	aceHeart := poker.NewCard(poker.SuitHeart, poker.RankAce)
	player.DealCard(aceSpade)
	player.DealCard(aceHeart)

	strength := bot.evaluateHandStrength(game, player)
	if strength < 0.8 { // Pocket aces should be very strong
		t.Errorf("Expected high strength for pocket aces, got %.2f", strength)
	}

	// Reset and test weak hand
	player.ResetForNewHand()
	twoSpade := poker.NewCard(poker.SuitSpade, poker.RankTwo)
	sevenHeart := poker.NewCard(poker.SuitHeart, poker.RankSeven)
	player.DealCard(twoSpade)
	player.DealCard(sevenHeart)

	strength = bot.evaluateHandStrength(game, player)
	if strength > 0.4 { // 2-7 offsuit should be weak
		t.Errorf("Expected low strength for 2-7 offsuit, got %.2f", strength)
	}
}
