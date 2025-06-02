package holdem_ai

import (
	"testing"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
)

func TestCreateBasicBot(t *testing.T) {
	bot := CreateBasicBot()

	if bot == nil {
		t.Fatal("CreateBasicBot returned nil")
	}

	// Verify it's a BasicBotDecisionMaker
	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateBasicBot did not return BasicBotDecisionMaker")
	}

	// Check default values
	if basicBot.Aggressiveness != 0.5 {
		t.Errorf("Expected aggressiveness 0.5, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.1 {
		t.Errorf("Expected bluff frequency 0.1, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateConservativeBot(t *testing.T) {
	bot := CreateConservativeBot()

	if bot == nil {
		t.Fatal("CreateConservativeBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateConservativeBot did not return BasicBotDecisionMaker")
	}

	// Conservative should have low aggressiveness and bluff frequency
	if basicBot.Aggressiveness != 0.2 {
		t.Errorf("Expected aggressiveness 0.2, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.05 {
		t.Errorf("Expected bluff frequency 0.05, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateAggressiveBot(t *testing.T) {
	bot := CreateAggressiveBot()

	if bot == nil {
		t.Fatal("CreateAggressiveBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateAggressiveBot did not return BasicBotDecisionMaker")
	}

	// Aggressive should have high aggressiveness and moderate bluff frequency
	if basicBot.Aggressiveness != 0.8 {
		t.Errorf("Expected aggressiveness 0.8, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.25 {
		t.Errorf("Expected bluff frequency 0.25, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateTightBot(t *testing.T) {
	bot := CreateTightBot()

	if bot == nil {
		t.Fatal("CreateTightBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateTightBot did not return BasicBotDecisionMaker")
	}

	// Tight should have very low aggressiveness and almost no bluffs
	if basicBot.Aggressiveness != 0.1 {
		t.Errorf("Expected aggressiveness 0.1, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.01 {
		t.Errorf("Expected bluff frequency 0.01, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateLooseBot(t *testing.T) {
	bot := CreateLooseBot()

	if bot == nil {
		t.Fatal("CreateLooseBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateLooseBot did not return BasicBotDecisionMaker")
	}

	// Loose should have very high aggressiveness and frequent bluffs
	if basicBot.Aggressiveness != 0.9 {
		t.Errorf("Expected aggressiveness 0.9, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.4 {
		t.Errorf("Expected bluff frequency 0.4, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateRandomBot(t *testing.T) {
	// Create multiple random bots to test variability
	bot1 := CreateRandomBot()
	bot2 := CreateRandomBot()
	bot3 := CreateRandomBot()

	if bot1 == nil || bot2 == nil || bot3 == nil {
		t.Fatal("CreateRandomBot returned nil")
	}

	basicBot1, ok1 := bot1.(*BasicBotDecisionMaker)
	basicBot2, ok2 := bot2.(*BasicBotDecisionMaker)
	basicBot3, ok3 := bot3.(*BasicBotDecisionMaker)

	if !ok1 || !ok2 || !ok3 {
		t.Fatal("CreateRandomBot did not return BasicBotDecisionMaker")
	}

	// Check ranges are valid
	bots := []*BasicBotDecisionMaker{basicBot1, basicBot2, basicBot3}
	for i, bot := range bots {
		// Aggressiveness should be between 0.3 and 0.9
		if bot.Aggressiveness < 0.3 || bot.Aggressiveness > 0.9 {
			t.Errorf("Bot %d aggressiveness %f out of range [0.3, 0.9]", i+1, bot.Aggressiveness)
		}

		// Bluff frequency should be between 0.05 and 0.3
		if bot.BluffFrequency < 0.05 || bot.BluffFrequency > 0.3 {
			t.Errorf("Bot %d bluff frequency %f out of range [0.05, 0.3]", i+1, bot.BluffFrequency)
		}
	}

	// Test that we get some variability (at least one different value)
	allSameAggression := (basicBot1.Aggressiveness == basicBot2.Aggressiveness &&
		basicBot2.Aggressiveness == basicBot3.Aggressiveness)
	allSameBluff := (basicBot1.BluffFrequency == basicBot2.BluffFrequency &&
		basicBot2.BluffFrequency == basicBot3.BluffFrequency)

	if allSameAggression && allSameBluff {
		t.Log("Warning: All random bots had identical parameters - this is unlikely but possible")
	}
}

func TestCreateCustomBot(t *testing.T) {
	customAggression := 0.65
	customBluff := 0.18

	bot := CreateCustomBot(customAggression, customBluff)

	if bot == nil {
		t.Fatal("CreateCustomBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateCustomBot did not return BasicBotDecisionMaker")
	}

	// Should match exactly what we specified
	if basicBot.Aggressiveness != customAggression {
		t.Errorf("Expected aggressiveness %f, got %f", customAggression, basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != customBluff {
		t.Errorf("Expected bluff frequency %f, got %f", customBluff, basicBot.BluffFrequency)
	}
}

func TestCreateNitBot(t *testing.T) {
	bot := CreateNitBot()

	if bot == nil {
		t.Fatal("CreateNitBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateNitBot did not return BasicBotDecisionMaker")
	}

	// Nit should be extremely tight
	if basicBot.Aggressiveness != 0.05 {
		t.Errorf("Expected aggressiveness 0.05, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.0 {
		t.Errorf("Expected bluff frequency 0.0, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateManiacBot(t *testing.T) {
	bot := CreateManiacBot()

	if bot == nil {
		t.Fatal("CreateManiacBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateManiacBot did not return BasicBotDecisionMaker")
	}

	// Maniac should be extremely aggressive
	if basicBot.Aggressiveness != 0.95 {
		t.Errorf("Expected aggressiveness 0.95, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.5 {
		t.Errorf("Expected bluff frequency 0.5, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateBalancedBot(t *testing.T) {
	bot := CreateBalancedBot()

	if bot == nil {
		t.Fatal("CreateBalancedBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateBalancedBot did not return BasicBotDecisionMaker")
	}

	// Balanced should have moderate values
	if basicBot.Aggressiveness != 0.6 {
		t.Errorf("Expected aggressiveness 0.6, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.15 {
		t.Errorf("Expected bluff frequency 0.15, got %f", basicBot.BluffFrequency)
	}
}

func TestCreateCallingStationBot(t *testing.T) {
	bot := CreateCallingStationBot()

	if bot == nil {
		t.Fatal("CreateCallingStationBot returned nil")
	}

	basicBot, ok := bot.(*BasicBotDecisionMaker)
	if !ok {
		t.Fatal("CreateCallingStationBot did not return BasicBotDecisionMaker")
	}

	// Calling station should be passive (low aggression, low bluffs)
	if basicBot.Aggressiveness != 0.3 {
		t.Errorf("Expected aggressiveness 0.3, got %f", basicBot.Aggressiveness)
	}
	if basicBot.BluffFrequency != 0.02 {
		t.Errorf("Expected bluff frequency 0.02, got %f", basicBot.BluffFrequency)
	}
}

func TestFactoryBotsImplementInterface(t *testing.T) {
	// Test that all factory functions return objects that implement IDecisionMaker
	factories := []func() IDecisionMaker{
		CreateBasicBot,
		CreateConservativeBot,
		CreateAggressiveBot,
		CreateTightBot,
		CreateLooseBot,
		CreateRandomBot,
		CreateNitBot,
		CreateManiacBot,
		CreateBalancedBot,
		CreateCallingStationBot,
	}

	for i, factory := range factories {
		bot := factory()
		if bot == nil {
			t.Errorf("Factory %d returned nil", i)
			continue
		}

		// Test that it implements the interface by calling MakeDecision
		game, player, _ := createTestGameSetup()
		ch := bot.MakeDecision(game, player)

		if ch == nil {
			t.Errorf("Factory %d bot returned nil channel", i)
			continue
		}

		// Verify we get a decision within reasonable time
		select {
		case action := <-ch:
			if !holdem.IsValidActionType(action.Type) {
				t.Errorf("Factory %d bot returned invalid action type: %d", i, action.Type)
			}
		case <-time.After(3 * time.Second):
			t.Errorf("Factory %d bot timed out", i)
		}
	}
}

func TestFactoryBotsBehavioralDifferences(t *testing.T) {
	// Test that different factory bots behave differently
	conservativeBot := CreateConservativeBot()
	aggressiveBot := CreateAggressiveBot()

	game, player, _ := createTestGameSetup()
	dealTestCards(game, player)

	// Run multiple trials to see behavioral differences
	conservativeFolds := 0
	aggressiveFolds := 0
	trials := 5 // Smaller number for factory test

	for i := 0; i < trials; i++ {
		// Reset game state
		game = holdem.NewGame(10, 20)
		player = holdem.NewPlayer(1, "Test", 1000)
		game.PlayerSit(player, 0)
		dealTestCards(game, player)

		// Conservative bot decision
		ch1 := conservativeBot.MakeDecision(game, player)
		action1 := <-ch1
		if action1.Type == holdem.ActionFold {
			conservativeFolds++
		}

		// Aggressive bot decision
		ch2 := aggressiveBot.MakeDecision(game, player)
		action2 := <-ch2
		if action2.Type == holdem.ActionFold {
			aggressiveFolds++
		}
	}

	// This is just a behavioral validation - we expect some difference but allow variance
	t.Logf("Conservative folds: %d, Aggressive folds: %d (out of %d trials)",
		conservativeFolds, aggressiveFolds, trials)
}

func TestFactoryBotsWithExtremeParameters(t *testing.T) {
	// Test custom bot with extreme parameters
	extremeBot1 := CreateCustomBot(0.0, 0.0) // Extremely passive
	extremeBot2 := CreateCustomBot(1.0, 1.0) // Extremely aggressive

	if extremeBot1 == nil || extremeBot2 == nil {
		t.Fatal("CreateCustomBot with extreme parameters returned nil")
	}

	// Both should still function and make decisions
	game, player, _ := createTestGameSetup()

	ch1 := extremeBot1.MakeDecision(game, player)
	ch2 := extremeBot2.MakeDecision(game, player)

	select {
	case action := <-ch1:
		if !holdem.IsValidActionType(action.Type) {
			t.Error("Extreme passive bot returned invalid action")
		}
	case <-time.After(3 * time.Second):
		t.Error("Extreme passive bot timed out")
	}

	select {
	case action := <-ch2:
		if !holdem.IsValidActionType(action.Type) {
			t.Error("Extreme aggressive bot returned invalid action")
		}
	case <-time.After(3 * time.Second):
		t.Error("Extreme aggressive bot timed out")
	}
}

func TestFactoryBotsParameterConsistency(t *testing.T) {
	// Test that creating the same bot type multiple times gives consistent parameters
	bot1 := CreateTightBot()
	bot2 := CreateTightBot()

	basicBot1, ok1 := bot1.(*BasicBotDecisionMaker)
	basicBot2, ok2 := bot2.(*BasicBotDecisionMaker)

	if !ok1 || !ok2 {
		t.Fatal("Failed to cast to BasicBotDecisionMaker")
	}

	if basicBot1.Aggressiveness != basicBot2.Aggressiveness {
		t.Errorf("TightBot aggressiveness inconsistent: %f vs %f",
			basicBot1.Aggressiveness, basicBot2.Aggressiveness)
	}

	if basicBot1.BluffFrequency != basicBot2.BluffFrequency {
		t.Errorf("TightBot bluff frequency inconsistent: %f vs %f",
			basicBot1.BluffFrequency, basicBot2.BluffFrequency)
	}
}
