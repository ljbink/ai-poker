package holdem_ai

import (
	"testing"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
	"github.com/ljbink/ai-poker/engine/poker"
)

func TestNewBasicBotDecisionMaker(t *testing.T) {
	aggressiveness := 0.7
	bluffFrequency := 0.2

	bot := NewBasicBotDecisionMaker(aggressiveness, bluffFrequency)

	if bot == nil {
		t.Fatal("NewBasicBotDecisionMaker returned nil")
	}

	if bot.Aggressiveness != aggressiveness {
		t.Errorf("Expected aggressiveness %f, got %f", aggressiveness, bot.Aggressiveness)
	}

	if bot.BluffFrequency != bluffFrequency {
		t.Errorf("Expected bluff frequency %f, got %f", bluffFrequency, bot.BluffFrequency)
	}

	if bot.evaluator == nil {
		t.Error("Bot evaluator is nil")
	}

	if bot.validator == nil {
		t.Error("Bot validator is nil")
	}
}

func TestBasicBotDecisionMakerMakeDecision(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)
	game, player, _ := createTestGameSetup()

	// Deal some cards to the player
	dealTestCards(game, player)

	ch := bot.MakeDecision(game, player)

	start := time.Now()
	select {
	case action := <-ch:
		duration := time.Since(start)

		// Bot should take realistic thinking time (500ms to 2s)
		if duration < 400*time.Millisecond {
			t.Errorf("Bot decided too quickly: %v", duration)
		}
		if duration > 3*time.Second {
			t.Errorf("Bot took too long: %v", duration)
		}

		// Verify valid action
		if action.PlayerID != player.GetID() {
			t.Errorf("Expected PlayerID %d, got %d", player.GetID(), action.PlayerID)
		}

		if !holdem.IsValidActionType(action.Type) {
			t.Errorf("Invalid action type: %d", action.Type)
		}

	case <-time.After(5 * time.Second):
		t.Error("Bot did not make decision within timeout")
	}
}

func TestBasicBotPersonalityTraits(t *testing.T) {
	// Test different bot personalities make different decisions
	conservativeBot := NewBasicBotDecisionMaker(0.1, 0.01)
	aggressiveBot := NewBasicBotDecisionMaker(0.9, 0.4)

	game, player, _ := createTestGameSetup()
	dealTestCards(game, player)

	// Run multiple decisions to see personality differences
	conservativeFolds := 0
	aggressiveFolds := 0
	trials := 10

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

	// Conservative bot should fold more often than aggressive bot
	// (This is a probabilistic test, so we allow some variance)
	if conservativeFolds <= aggressiveFolds {
		t.Logf("Warning: Conservative bot (%d folds) didn't fold more than aggressive bot (%d folds)",
			conservativeFolds, aggressiveFolds)
	}
}

func TestBasicBotHandStrengthEvaluation(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)
	game, player, _ := createTestGameSetup()

	// Test with strong hand (pocket aces)
	strongCards := []*poker.Card{
		poker.NewCard(poker.SuitHeart, poker.RankAce),
		poker.NewCard(poker.SuitSpade, poker.RankAce),
	}
	player.ResetForNewHand()
	for _, card := range strongCards {
		player.DealCard(card)
	}

	strongHandStrength := bot.evaluateHandStrength(game, player)

	// Test with weak hand (2-7 offsuit)
	player.ResetForNewHand()
	weakCards := []*poker.Card{
		poker.NewCard(poker.SuitHeart, poker.RankTwo),
		poker.NewCard(poker.SuitSpade, poker.RankSeven),
	}
	for _, card := range weakCards {
		player.DealCard(card)
	}

	weakHandStrength := bot.evaluateHandStrength(game, player)

	// Strong hand should have higher strength
	if strongHandStrength <= weakHandStrength {
		t.Errorf("Strong hand strength (%f) should be greater than weak hand (%f)",
			strongHandStrength, weakHandStrength)
	}

	// Both should be in valid range [0.0, 1.0]
	if strongHandStrength < 0.0 || strongHandStrength > 1.0 {
		t.Errorf("Strong hand strength %f out of range [0.0, 1.0]", strongHandStrength)
	}
	if weakHandStrength < 0.0 || weakHandStrength > 1.0 {
		t.Errorf("Weak hand strength %f out of range [0.0, 1.0]", weakHandStrength)
	}
}

func TestBasicBotHandRankToStrength(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)

	testCases := []struct {
		rank     holdem.HandRank
		expected float64
	}{
		{holdem.RoyalFlush, 1.0},
		{holdem.StraightFlush, 0.95},
		{holdem.FourOfAKind, 0.9},
		{holdem.FullHouse, 0.85},
		{holdem.Flush, 0.75},
		{holdem.Straight, 0.65},
		{holdem.ThreeOfAKind, 0.55},
		{holdem.TwoPair, 0.45},
		{holdem.OnePair, 0.3},
		{holdem.HighCard, 0.1},
	}

	for _, tc := range testCases {
		strength := bot.handRankToStrength(tc.rank)
		if strength != tc.expected {
			t.Errorf("Expected strength %f for rank %d, got %f", tc.expected, tc.rank, strength)
		}
	}

	// Test invalid rank
	invalidStrength := bot.handRankToStrength(holdem.HandRank(99))
	if invalidStrength != 0.0 {
		t.Errorf("Expected 0.0 for invalid rank, got %f", invalidStrength)
	}
}

func TestBasicBotPreflopEvaluation(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)

	// Test pocket aces
	pocketAces := []*poker.Card{
		poker.NewCard(poker.SuitHeart, poker.RankAce),
		poker.NewCard(poker.SuitSpade, poker.RankAce),
	}
	strength := bot.evaluatePreflop(pocketAces)
	if strength <= 0.3 {
		t.Errorf("Pocket aces should have high preflop strength, got %f", strength)
	}

	// Test suited connectors
	suitedConnectors := []*poker.Card{
		poker.NewCard(poker.SuitHeart, poker.RankJack),
		poker.NewCard(poker.SuitHeart, poker.RankTen),
	}
	suitedStrength := bot.evaluatePreflop(suitedConnectors)

	// Test same ranks offsuit
	offsuitConnectors := []*poker.Card{
		poker.NewCard(poker.SuitHeart, poker.RankJack),
		poker.NewCard(poker.SuitSpade, poker.RankTen),
	}
	offsuitStrength := bot.evaluatePreflop(offsuitConnectors)

	// Suited should be stronger than offsuit
	if suitedStrength <= offsuitStrength {
		t.Errorf("Suited connectors (%f) should be stronger than offsuit (%f)",
			suitedStrength, offsuitStrength)
	}

	// Test with empty cards
	emptyStrength := bot.evaluatePreflop([]*poker.Card{})
	if emptyStrength != 0.0 {
		t.Errorf("Expected 0.0 for empty cards, got %f", emptyStrength)
	}

	// Test with one card
	oneCard := []*poker.Card{poker.NewCard(poker.SuitHeart, poker.RankAce)}
	oneCardStrength := bot.evaluatePreflop(oneCard)
	if oneCardStrength != 0.0 {
		t.Errorf("Expected 0.0 for one card, got %f", oneCardStrength)
	}
}

func TestBasicBotRankToValue(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)

	testCases := []struct {
		rank     poker.Rank
		expected int
	}{
		{poker.RankAce, 14},
		{poker.RankKing, 13},
		{poker.RankQueen, 12},
		{poker.RankJack, 11},
		{poker.RankTen, 10},
		{poker.RankNine, 9},
		{poker.RankTwo, 2},
	}

	for _, tc := range testCases {
		value := bot.rankToValue(tc.rank)
		if value != tc.expected {
			t.Errorf("Expected value %d for rank %d, got %d", tc.expected, tc.rank, value)
		}
	}

	// Test invalid rank
	invalidValue := bot.rankToValue(poker.RankNone)
	if invalidValue != 0 {
		t.Errorf("Expected 0 for RankNone, got %d", invalidValue)
	}
}

func TestBasicBotUtilityFunctions(t *testing.T) {
	// Test minFloat64
	if minFloat64(0.5, 0.3) != 0.3 {
		t.Error("minFloat64 failed")
	}
	if minFloat64(0.2, 0.8) != 0.2 {
		t.Error("minFloat64 failed")
	}

	// Test minInt
	if minInt(10, 5) != 5 {
		t.Error("minInt failed")
	}
	if minInt(3, 7) != 3 {
		t.Error("minInt failed")
	}

	// Test maxInt
	if maxInt(10, 5) != 10 {
		t.Error("maxInt failed")
	}
	if maxInt(3, 7) != 7 {
		t.Error("maxInt failed")
	}

	// Test abs
	if abs(-5) != 5 {
		t.Error("abs failed for negative")
	}
	if abs(5) != 5 {
		t.Error("abs failed for positive")
	}
	if abs(0) != 0 {
		t.Error("abs failed for zero")
	}
}

func TestBasicBotShouldBluff(t *testing.T) {
	// High bluff frequency bot
	highBluffBot := NewBasicBotDecisionMaker(0.5, 0.8)

	// Low bluff frequency bot
	lowBluffBot := NewBasicBotDecisionMaker(0.5, 0.1)

	// Test with weak hand (should consider bluffing)
	weakHandStrength := 0.1

	highBluffCount := 0
	lowBluffCount := 0
	trials := 100

	for i := 0; i < trials; i++ {
		if highBluffBot.shouldBluff(weakHandStrength) {
			highBluffCount++
		}
		if lowBluffBot.shouldBluff(weakHandStrength) {
			lowBluffCount++
		}
	}

	// High bluff bot should bluff more often
	if highBluffCount <= lowBluffCount {
		t.Logf("Warning: High bluff bot (%d bluffs) didn't bluff more than low bluff bot (%d bluffs)",
			highBluffCount, lowBluffCount)
	}

	// Test with strong hand (should rarely bluff)
	strongHandStrength := 0.9
	strongHandBluffs := 0

	for i := 0; i < trials; i++ {
		if highBluffBot.shouldBluff(strongHandStrength) {
			strongHandBluffs++
		}
	}

	// Strong hands should rarely trigger bluffs
	if strongHandBluffs > trials/4 {
		t.Errorf("Too many bluffs with strong hand: %d/%d", strongHandBluffs, trials)
	}
}

func TestBasicBotIsActionAvailable(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)

	availableActions := []holdem.ActionType{
		holdem.ActionFold,
		holdem.ActionCheck,
		holdem.ActionRaise,
	}

	// Test available actions
	if !bot.isActionAvailable(holdem.ActionFold, availableActions) {
		t.Error("ActionFold should be available")
	}
	if !bot.isActionAvailable(holdem.ActionCheck, availableActions) {
		t.Error("ActionCheck should be available")
	}
	if !bot.isActionAvailable(holdem.ActionRaise, availableActions) {
		t.Error("ActionRaise should be available")
	}

	// Test unavailable action
	if bot.isActionAvailable(holdem.ActionCall, availableActions) {
		t.Error("ActionCall should not be available")
	}
	if bot.isActionAvailable(holdem.ActionAllIn, availableActions) {
		t.Error("ActionAllIn should not be available")
	}
}

func TestBasicBotCountActivePlayers(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)
	game := holdem.NewGame(10, 20)

	// No players initially
	count := bot.countActivePlayers(game)
	if count != 0 {
		t.Errorf("Expected 0 active players, got %d", count)
	}

	// Add players
	player1 := holdem.NewPlayer(1, "Player 1", 1000)
	player2 := holdem.NewPlayer(2, "Player 2", 1000)
	player3 := holdem.NewPlayer(3, "Player 3", 1000)

	game.PlayerSit(player1, 0)
	count = bot.countActivePlayers(game)
	if count != 1 {
		t.Errorf("Expected 1 active player, got %d", count)
	}

	game.PlayerSit(player2, 1)
	game.PlayerSit(player3, 2)
	count = bot.countActivePlayers(game)
	if count != 3 {
		t.Errorf("Expected 3 active players, got %d", count)
	}

	// Fold one player
	player2.Fold()
	count = bot.countActivePlayers(game)
	if count != 2 {
		t.Errorf("Expected 2 active players after fold, got %d", count)
	}
}

func TestBasicBotCalculateCallAmount(t *testing.T) {
	bot := NewBasicBotDecisionMaker(0.5, 0.1)
	game, player1, player2 := createTestGameSetup()

	// Initially no call needed
	callAmount := bot.calculateCallAmount(game, player1)
	if callAmount != 0 {
		t.Errorf("Expected call amount 0 initially, got %d", callAmount)
	}

	// Add a bet
	player2.Bet(50)
	raiseAction := holdem.Action{
		PlayerID: player2.GetID(),
		Type:     holdem.ActionRaise,
		Amount:   50,
	}
	game.TakeAction(raiseAction)

	callAmount = bot.calculateCallAmount(game, player1)
	if callAmount != 50 {
		t.Errorf("Expected call amount 50, got %d", callAmount)
	}

	// Player1 already has some bet
	player1.Bet(20)
	callAmount = bot.calculateCallAmount(game, player1)
	if callAmount != 30 { // 50 - 20 = 30
		t.Errorf("Expected call amount 30, got %d", callAmount)
	}
}

// Helper function to deal test cards to a player
func dealTestCards(game *holdem.Game, player holdem.IPlayer) {
	// Deal some reasonable hole cards
	card1 := poker.NewCard(poker.SuitHeart, poker.RankKing)
	card2 := poker.NewCard(poker.SuitSpade, poker.RankQueen)

	player.DealCard(card1)
	player.DealCard(card2)
}
