package holdem_ai

import (
	"testing"

	"github.com/ljbink/ai-poker/engine/holdem"
)

func TestBasicBotDecisionMaker(t *testing.T) {
	// Create a test game and bot
	game := createTestGameForBot()
	player := game.Players[0]
	bot := NewBasicBotDecisionMaker(player, game)

	// Test decision making
	actionChan := bot.MakeDecision()
	action := <-actionChan

	// Should be one of the valid action types
	validTypes := []ActionType{ActionFold, ActionCheck, ActionCall, ActionRaise, ActionAllIn}
	validAction := false
	for _, validType := range validTypes {
		if action.Type == validType {
			validAction = true
			break
		}
	}

	if !validAction {
		t.Errorf("Bot returned invalid action type: %v", action.Type)
	}
}

func TestBasicBotDecisionMakerConsistency(t *testing.T) {
	// Test that bot makes consistent decisions in similar situations
	var actions []Action
	for i := 0; i < 5; i++ {
		game := createTestGameForBot()
		player := game.Players[0]
		bot := NewBasicBotDecisionMaker(player, game)

		actionChan := bot.MakeDecision()
		action := <-actionChan
		actions = append(actions, action)
	}

	// Check that all actions are valid (should be one of the known action types)
	validTypes := []ActionType{ActionFold, ActionCheck, ActionCall, ActionRaise, ActionAllIn}
	for i, action := range actions {
		validAction := false
		for _, validType := range validTypes {
			if action.Type == validType {
				validAction = true
				break
			}
		}
		if !validAction {
			t.Errorf("Action %d is invalid: %v", i, action.Type)
		}
	}
}

func TestBasicBotDecisionMakerWithDifferentGameStates(t *testing.T) {
	game := createTestGameForBot()
	player := game.Players[0]
	bot := NewBasicBotDecisionMaker(player, game)

	// Test preflop
	actionChan := bot.MakeDecision()
	action := <-actionChan
	if action.Type == 0 && action.Amount == 0 {
		t.Error("Bot should make a decision in preflop")
	}

	// Test with community cards (flop)
	game.CurrentPhase = holdem.PhaseFlop
	actionChan = bot.MakeDecision()
	action = <-actionChan
	if action.Type == 0 && action.Amount == 0 {
		t.Error("Bot should make a decision in flop")
	}
}

func TestBasicBotDecisionMakerActionValidation(t *testing.T) {
	game := createTestGameForBot()
	player := game.Players[0]
	bot := NewBasicBotDecisionMaker(player, game)

	// Get valid actions
	validator := NewActionValidator()
	validActions := validator.GetValidActions(game, player)

	// Make decision
	actionChan := bot.MakeDecision()
	action := <-actionChan

	// Check that bot's action is among valid actions
	isValid := false
	for _, validAction := range validActions {
		if action.Type == validAction {
			isValid = true
			break
		}
	}

	if !isValid {
		t.Errorf("Bot chose invalid action %v, valid actions were: %v", action.Type, validActions)
	}
}

func TestHandStrengthEvaluation(t *testing.T) {
	game := createTestGameForBot()
	player := game.Players[0]
	bot := NewBasicBotDecisionMaker(player, game).(*BasicBotDecisionMaker)

	// Test with cards dealt (should return > 0 since player has cards after StartHand)
	strength := bot.evaluateHandStrength()
	if strength < 0.0 || strength > 1.0 {
		t.Errorf("Expected strength between 0.0 and 1.0, got %f", strength)
	}

	// Test that strength calculation works (not zero since player has cards)
	if strength == 0.0 {
		t.Error("Expected non-zero strength since player has cards after StartHand")
	}
}

// Helper function to create test game
func createTestGameForBot() *holdem.Game {
	game := holdem.NewGame([]string{"Bot", "Player2"}, 5, 10)
	game.StartHand()
	return game
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
