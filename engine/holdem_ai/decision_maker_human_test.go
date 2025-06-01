package holdem_ai

import (
	"testing"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
)

func TestNewHumanUserDecisionMaker(t *testing.T) {
	game := createTestGame()
	player := game.Players[0]
	dm := NewHumanUserDecisionMaker(player, game)

	if dm.timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", dm.timeout)
	}

	// Test that decision maker works correctly
	actions := dm.GetValidActions()
	if len(actions) == 0 {
		t.Error("Expected valid actions when game is active")
	}
}

func TestSendAction(t *testing.T) {
	game := createTestGame()
	player := game.Players[0]
	dm := NewHumanUserDecisionMaker(player, game)

	// Test valid action (use fold which is always valid)
	action := CreateFoldAction()
	err := dm.SendAction(action)
	if err != nil {
		t.Errorf("Unexpected error sending valid action: %v", err)
	}
}

func TestMakeDecisionWithAction(t *testing.T) {
	game := createTestGame()
	player := game.Players[0]
	dm := NewHumanUserDecisionMaker(player, game)

	// Set shorter timeout for testing
	dm.SetTimeout(100 * time.Millisecond)

	// Send action in separate goroutine (use fold which is always valid)
	expectedAction := CreateFoldAction()
	go func() {
		time.Sleep(10 * time.Millisecond) // Small delay
		dm.SendAction(expectedAction)
	}()

	// Make decision and wait for result
	actionChan := dm.MakeDecision()
	action := <-actionChan
	if action.Type != expectedAction.Type {
		t.Errorf("Expected action %s, got %s", expectedAction.Type, action.Type)
	}
}

func TestMakeDecisionTimeout(t *testing.T) {
	game := createTestGame()
	player := game.Players[0]
	dm := NewHumanUserDecisionMaker(player, game)

	// Set very short timeout
	dm.SetTimeout(10 * time.Millisecond)

	// Make decision without sending action (should timeout)
	actionChan := dm.MakeDecision()
	action := <-actionChan
	if action.Type != ActionFold {
		t.Errorf("Expected fold action on timeout, got %s", action.Type)
	}
}

func TestGetValidActions(t *testing.T) {
	game := createTestGame()
	player := game.Players[0]
	dm := NewHumanUserDecisionMaker(player, game)

	actions := dm.GetValidActions()
	if len(actions) == 0 {
		t.Error("Expected valid actions when game is active")
	}
}

func TestActionNeededCallback(t *testing.T) {
	game := createTestGame()
	player := game.Players[0]
	dm := NewHumanUserDecisionMaker(player, game)

	callbackCalled := false
	var receivedGame *holdem.Game
	var receivedPlayer holdem.IPlayer
	var receivedActions []ActionType

	// Set callback
	dm.SetActionNeededCallback(func(g *holdem.Game, p holdem.IPlayer, actions []ActionType) {
		callbackCalled = true
		receivedGame = g
		receivedPlayer = p
		receivedActions = actions
	})

	// Set short timeout and don't send action (will timeout)
	dm.SetTimeout(10 * time.Millisecond)

	// Make decision (will trigger callback)
	actionChan := dm.MakeDecision()
	<-actionChan // Wait for decision to complete

	// Give callback time to execute
	time.Sleep(20 * time.Millisecond)

	if !callbackCalled {
		t.Error("Expected callback to be called")
	}

	if receivedGame == nil {
		t.Error("Expected to receive game")
	}

	if receivedPlayer == nil {
		t.Error("Expected to receive player")
	}

	if receivedPlayer.GetID() != player.GetID() {
		t.Error("Expected callback to receive correct player")
	}

	if len(receivedActions) == 0 {
		t.Error("Expected to receive valid actions")
	}
}

func TestConvenienceActions(t *testing.T) {
	tests := []struct {
		name     string
		creator  func() Action
		expected ActionType
	}{
		{"Fold", CreateFoldAction, ActionFold},
		{"Check", CreateCheckAction, ActionCheck},
		{"Call", CreateCallAction, ActionCall},
		{"AllIn", CreateAllInAction, ActionAllIn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := tt.creator()
			if action.Type != tt.expected {
				t.Errorf("Expected action type %s, got %s", tt.expected, action.Type)
			}
		})
	}
}

func TestCreateRaiseAction(t *testing.T) {
	amount := 100
	action := CreateRaiseAction(amount)

	if action.Type != ActionRaise {
		t.Errorf("Expected action type %s, got %s", ActionRaise, action.Type)
	}

	if action.Amount != amount {
		t.Errorf("Expected amount %d, got %d", amount, action.Amount)
	}
}

func TestHelperFunctions(t *testing.T) {
	game := createTestGame()
	player := game.Players[0]

	// Test FormatGamePhase
	phase := FormatGamePhase(game.CurrentPhase)
	if phase != "Preflop" {
		t.Errorf("Expected phase 'Preflop', got '%s'", phase)
	}

	// Test FormatPlayerCards
	playerCards := FormatPlayerCards(player)
	if playerCards == "" {
		t.Error("Expected player cards to be formatted")
	}

	// Test FormatCommunityCards (should be empty in preflop)
	communityCards := FormatCommunityCards(game)
	if communityCards != "" {
		t.Errorf("Expected empty community cards in preflop, got '%s'", communityCards)
	}

	// Test CalculateCallAmount
	callAmount := CalculateCallAmount(game, player)
	if callAmount < 0 {
		t.Error("Call amount should not be negative")
	}

	// Test CalculateMinRaise
	minRaise := CalculateMinRaise(game)
	if minRaise != game.BigBlind {
		t.Errorf("Expected min raise %d, got %d", game.BigBlind, minRaise)
	}

	// Test CalculateMaxRaise
	maxRaise := CalculateMaxRaise(game, player)
	if maxRaise < 0 {
		t.Error("Max raise should not be negative")
	}
}

// Helper function to create test game
func createTestGame() *holdem.Game {
	game := holdem.NewGame([]string{"Player1", "Player2"}, 5, 10)
	game.StartHand()
	return game
}
