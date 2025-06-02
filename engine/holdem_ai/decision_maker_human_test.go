package holdem_ai

import (
	"testing"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
)

func TestNewHumanDecisionMaker(t *testing.T) {
	human := NewHumanDecisionMaker()

	if human == nil {
		t.Fatal("NewHumanDecisionMaker returned nil")
	}

	if human.validator == nil {
		t.Error("HumanDecisionMaker validator is nil")
	}

	if human.actionChannel == nil {
		t.Error("HumanDecisionMaker actionChannel is nil")
	}

	// Test channel capacity
	if cap(human.actionChannel) != 1 {
		t.Errorf("Expected actionChannel capacity 1, got %d", cap(human.actionChannel))
	}
}

func TestHumanDecisionMakerTimeout(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player, _ := createTestGameSetup()

	// Don't provide any action - should timeout
	ch := human.MakeDecision(game, player)

	start := time.Now()
	select {
	case action := <-ch:
		duration := time.Since(start)

		// Should timeout around 60 seconds, but we'll be lenient for testing
		if duration < 50*time.Second {
			t.Errorf("Expected timeout around 60s, got %v", duration)
		}

		// Should return fold action on timeout
		if action.Type != holdem.ActionFold {
			t.Errorf("Expected fold action on timeout, got %d", action.Type)
		}

		if action.PlayerID != player.GetID() {
			t.Errorf("Expected PlayerID %d, got %d", player.GetID(), action.PlayerID)
		}

		if action.Amount != 0 {
			t.Errorf("Expected Amount 0 for fold, got %d", action.Amount)
		}

	case <-time.After(65 * time.Second):
		t.Error("Decision maker did not timeout within expected time")
	}
}

func TestHumanDecisionMakerValidAction(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player, _ := createTestGameSetup()

	// Provide a valid action
	validAction := holdem.Action{
		PlayerID: player.GetID(),
		Type:     holdem.ActionCheck,
		Amount:   0,
	}

	ch := human.MakeDecision(game, player)

	// Send action after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		human.SetAction(validAction)
	}()

	select {
	case receivedAction := <-ch:
		if receivedAction.PlayerID != validAction.PlayerID {
			t.Errorf("Expected PlayerID %d, got %d", validAction.PlayerID, receivedAction.PlayerID)
		}
		if receivedAction.Type != validAction.Type {
			t.Errorf("Expected Type %d, got %d", validAction.Type, receivedAction.Type)
		}
		if receivedAction.Amount != validAction.Amount {
			t.Errorf("Expected Amount %d, got %d", validAction.Amount, receivedAction.Amount)
		}
	case <-time.After(2 * time.Second):
		t.Error("Did not receive action within timeout")
	}
}

func TestHumanDecisionMakerInvalidAction(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player, _ := createTestGameSetup()

	// Provide an invalid action (wrong player ID)
	invalidAction := holdem.Action{
		PlayerID: 999, // Wrong player ID
		Type:     holdem.ActionCheck,
		Amount:   0,
	}

	ch := human.MakeDecision(game, player)

	// Send invalid action
	go func() {
		time.Sleep(100 * time.Millisecond)
		human.SetAction(invalidAction)
	}()

	select {
	case receivedAction := <-ch:
		// Should receive a fold action as fallback
		if receivedAction.Type != holdem.ActionFold {
			t.Errorf("Expected fold action as fallback, got %d", receivedAction.Type)
		}
		if receivedAction.PlayerID != player.GetID() {
			t.Errorf("Expected PlayerID %d, got %d", player.GetID(), receivedAction.PlayerID)
		}
		if receivedAction.Amount != 0 {
			t.Errorf("Expected Amount 0 for fold, got %d", receivedAction.Amount)
		}
	case <-time.After(2 * time.Second):
		t.Error("Did not receive fallback action within timeout")
	}
}

func TestHumanDecisionMakerSetActionChannelBehavior(t *testing.T) {
	human := NewHumanDecisionMaker()

	action := holdem.Action{
		PlayerID: 1,
		Type:     holdem.ActionRaise,
		Amount:   100,
	}

	// Test successful action setting
	human.SetAction(action)

	// Channel should have the action
	select {
	case receivedAction := <-human.actionChannel:
		if receivedAction.PlayerID != action.PlayerID {
			t.Errorf("Expected PlayerID %d, got %d", action.PlayerID, receivedAction.PlayerID)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Action was not set in channel")
	}

	// Test channel full behavior - should not block
	human.SetAction(action) // Fill the channel
	human.SetAction(action) // This should not block (channel full case)

	// Channel should still have one action
	select {
	case <-human.actionChannel:
		// Expected
	default:
		t.Error("Channel should have had an action")
	}

	// Channel should now be empty
	select {
	case <-human.actionChannel:
		t.Error("Channel should be empty after reading")
	default:
		// Expected
	}
}

func TestHumanDecisionMakerGetAvailableActions(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player, _ := createTestGameSetup()

	actions := human.GetAvailableActions(game, player)

	// Should return some actions for a valid game state
	if len(actions) == 0 {
		t.Error("Expected some available actions")
	}

	// All returned actions should be valid
	for _, action := range actions {
		if !holdem.IsValidActionType(action) {
			t.Errorf("Invalid action type returned: %d", action)
		}
	}
}

func TestHumanDecisionMakerRaiseAmounts(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player, _ := createTestGameSetup()

	minRaise := human.GetMinRaiseAmount(game, player)
	maxRaise := human.GetMaxRaiseAmount(game, player)

	// Min raise should be positive for most cases
	if minRaise < 0 {
		t.Errorf("Min raise should not be negative, got %d", minRaise)
	}

	// Max raise should be at least the player's chips
	if maxRaise < 0 {
		t.Errorf("Max raise should not be negative, got %d", maxRaise)
	}

	// Max should be >= min (unless both are 0)
	if maxRaise > 0 && minRaise > maxRaise {
		t.Errorf("Min raise (%d) should not exceed max raise (%d)", minRaise, maxRaise)
	}
}

func TestHumanDecisionMakerValidateAction(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player, _ := createTestGameSetup()

	// Test valid action
	validAction := holdem.Action{
		PlayerID: player.GetID(),
		Type:     holdem.ActionCheck,
		Amount:   0,
	}

	err := human.ValidateAction(game, player, validAction)
	if err != nil {
		t.Errorf("Expected no error for valid action, got: %v", err)
	}

	// Test invalid action
	invalidAction := holdem.Action{
		PlayerID: 999, // Wrong player ID
		Type:     holdem.ActionCheck,
		Amount:   0,
	}

	err = human.ValidateAction(game, player, invalidAction)
	if err == nil {
		t.Error("Expected error for invalid action")
	}
}

func TestHumanDecisionMakerGetCallAmount(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player1, player2 := createTestGameSetup()

	// Initially, call amount should be 0 (no bets)
	callAmount := human.GetCallAmount(game, player1)
	if callAmount != 0 {
		t.Errorf("Expected call amount 0 initially, got %d", callAmount)
	}

	// Add a bet from player2
	player2.Bet(50)
	raiseAction := holdem.Action{
		PlayerID: player2.GetID(),
		Type:     holdem.ActionRaise,
		Amount:   50,
	}
	game.TakeAction(raiseAction)

	// Now player1 should need to call 50
	callAmount = human.GetCallAmount(game, player1)
	if callAmount != 50 {
		t.Errorf("Expected call amount 50, got %d", callAmount)
	}

	// If player1 has already bet some amount
	player1.Bet(20)
	callAmount = human.GetCallAmount(game, player1)
	if callAmount != 30 { // 50 - 20 = 30
		t.Errorf("Expected call amount 30, got %d", callAmount)
	}
}

func TestHumanDecisionMakerGetCurrentPhaseActions(t *testing.T) {
	human := NewHumanDecisionMaker()
	game, player, _ := createTestGameSetup()

	// Test different phases
	phases := []holdem.GamePhase{
		holdem.PhasePreflop,
		holdem.PhaseFlop,
		holdem.PhaseTurn,
		holdem.PhaseRiver,
	}

	for _, phase := range phases {
		game.SetCurrentPhase(phase)

		// Add an action in this phase
		action := holdem.Action{
			PlayerID: player.GetID(),
			Type:     holdem.ActionCheck,
			Amount:   0,
		}
		game.TakeAction(action)

		// Get actions for this phase
		actions := human.getCurrentPhaseActions(game)

		if len(actions) == 0 {
			t.Errorf("Expected actions in phase %d", phase)
		}

		if actions[0].Type != holdem.ActionCheck {
			t.Errorf("Expected check action in phase %d, got %d", phase, actions[0].Type)
		}
	}

	// Test invalid phase
	game.SetCurrentPhase(holdem.GamePhase(99))
	actions := human.getCurrentPhaseActions(game)
	if len(actions) != 0 {
		t.Errorf("Expected no actions for invalid phase, got %d", len(actions))
	}
}
