package holdem_ai

import (
	"testing"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
)

// TestIDecisionMakerInterfaceCompliance ensures all implementations satisfy the interface
func TestIDecisionMakerInterfaceCompliance(t *testing.T) {
	var _ IDecisionMaker = &HumanDecisionMaker{}
	var _ IDecisionMaker = &BasicBotDecisionMaker{}
}

// TestDecisionMakerChannelBehavior tests the channel-based decision making
func TestDecisionMakerChannelBehavior(t *testing.T) {
	// Test that all decision makers return proper channels
	game := holdem.NewGame(10, 20)
	player := holdem.NewPlayer(1, "Test Player", 1000)
	game.PlayerSit(player, 0)

	testCases := []struct {
		name  string
		maker IDecisionMaker
	}{
		{"HumanDecisionMaker", NewHumanDecisionMaker()},
		{"BasicBotDecisionMaker", NewBasicBotDecisionMaker(0.5, 0.1)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ch := tc.maker.MakeDecision(game, player)

			// Verify channel is not nil
			if ch == nil {
				t.Error("MakeDecision returned nil channel")
			}

			// Verify channel type
			select {
			case action := <-ch:
				// Verify action has valid player ID
				if action.PlayerID != player.GetID() {
					t.Errorf("Expected action PlayerID %d, got %d", player.GetID(), action.PlayerID)
				}

				// Verify action type is valid
				if !holdem.IsValidActionType(action.Type) {
					t.Errorf("Invalid action type: %d", action.Type)
				}
			case <-time.After(5 * time.Second):
				// For bots, this should not timeout in normal cases
				if tc.name != "HumanDecisionMaker" {
					t.Error("Bot decision timed out unexpectedly")
				}
			}
		})
	}
}

// TestDecisionMakerWithEmptyGame tests behavior with minimal game state
func TestDecisionMakerWithEmptyGame(t *testing.T) {
	game := holdem.NewGame(10, 20)
	player := holdem.NewPlayer(1, "Test Player", 1000)

	// Don't sit the player - test with invalid game state
	bot := NewBasicBotDecisionMaker(0.5, 0.1)

	ch := bot.MakeDecision(game, player)

	select {
	case action := <-ch:
		// Should get a fold action when no valid actions available
		if action.Type != holdem.ActionFold {
			t.Errorf("Expected fold action for invalid game state, got %d", action.Type)
		}
	case <-time.After(2 * time.Second):
		t.Error("Bot did not make decision within timeout")
	}
}

// TestDecisionMakerWithNilInputs tests error handling with nil inputs
func TestDecisionMakerWithNilInputs(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)

	// Test with nil game
	ch := bot.MakeDecision(nil, holdem.NewPlayer(1, "Test", 1000))

	select {
	case action := <-ch:
		// Should handle gracefully and return fold
		if action.Type != holdem.ActionFold {
			t.Errorf("Expected fold action for nil game, got %d", action.Type)
		}
	case <-time.After(2 * time.Second):
		t.Error("Bot did not handle nil game within timeout")
	}

	// Test with nil player
	game := holdem.NewGame(10, 20)
	ch = bot.MakeDecision(game, nil)

	select {
	case action := <-ch:
		// Should handle gracefully
		if action.PlayerID != 0 {
			t.Errorf("Expected PlayerID 0 for nil player, got %d", action.PlayerID)
		}
	case <-time.After(2 * time.Second):
		t.Error("Bot did not handle nil player within timeout")
	}
}

// Helper function to create a basic game setup for testing
func createTestGameSetup() (*holdem.Game, holdem.IPlayer, holdem.IPlayer) {
	game := holdem.NewGame(10, 20)
	player1 := holdem.NewPlayer(1, "Player 1", 1000)
	player2 := holdem.NewPlayer(2, "Player 2", 1000)

	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 1)

	return game, player1, player2
}
