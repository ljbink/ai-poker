package holdem

import (
	"testing"
)

func TestNewActionValidator(t *testing.T) {
	validator := NewActionValidator()
	if validator == nil {
		t.Fatal("NewActionValidator returned nil")
	}
}

func TestValidationErrorCode(t *testing.T) {
	// Test validation error codes have different values
	codes := []ValidationErrorCode{
		ErrorInvalidPlayer,
		ErrorInvalidAction,
		ErrorInsufficientChips,
		ErrorInvalidAmount,
		ErrorOutOfTurn,
		ErrorGameState,
		ErrorActionNotAllowed,
	}

	seen := make(map[ValidationErrorCode]bool)
	for _, code := range codes {
		if seen[code] {
			t.Errorf("Duplicate validation error code: %d", code)
		}
		seen[code] = true
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Message: "Test error",
		Code:    ErrorInvalidPlayer,
	}

	if err.Error() != "Test error" {
		t.Errorf("Expected error message 'Test error', got %s", err.Error())
	}
}

func TestValidateBasicAction(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0) // Need to sit player for turn validation

	// Test valid action
	action := Action{PlayerID: 1, Type: ActionFold, Amount: 0}
	err := validator.ValidateAction(game, player, action)
	if err != nil {
		t.Errorf("Unexpected error for valid action: %v", err)
	}

	// Test invalid player ID
	action = Action{PlayerID: 0, Type: ActionFold, Amount: 0}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for invalid player ID")
	}
	if err.Code != ErrorInvalidPlayer {
		t.Errorf("Expected ErrorInvalidPlayer, got %d", err.Code)
	}

	// Test negative amount
	action = Action{PlayerID: 1, Type: ActionRaise, Amount: -100}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for negative amount")
	}
	if err.Code != ErrorInvalidAmount {
		t.Errorf("Expected ErrorInvalidAmount, got %d", err.Code)
	}
}

func TestValidateNilInputs(t *testing.T) {
	validator := NewActionValidator()
	action := Action{PlayerID: 1, Type: ActionFold, Amount: 0}

	// Test nil game
	err := validator.ValidateAction(nil, NewPlayer(1, "Player 1", 1000), action)
	if err == nil {
		t.Error("Expected error for nil game")
	}

	// Test nil player
	err = validator.ValidateAction(NewGame(10, 20), nil, action)
	if err == nil {
		t.Error("Expected error for nil player")
	}
}

func TestValidatePlayerMismatch(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)

	// Test action with different player ID
	action := Action{PlayerID: 2, Type: ActionFold, Amount: 0}
	err := validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for player ID mismatch")
	}
	if err.Code != ErrorInvalidPlayer {
		t.Errorf("Expected ErrorInvalidPlayer, got %d", err.Code)
	}
}

func TestOutOfTurnValidation(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	// Seat two players; current player will be the first non-folded (player1)
	player1 := NewPlayer(1, "P1", 1000)
	player2 := NewPlayer(2, "P2", 1000)
	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 1)

	// Player2 tries to act out of turn
	action := Action{PlayerID: 2, Type: ActionCheck, Amount: 0}
	if err := validator.ValidateAction(game, player2, action); err == nil {
		t.Error("Expected out-of-turn error for player2")
	} else if err.Code != ErrorOutOfTurn {
		t.Errorf("Expected ErrorOutOfTurn, got %v", err.Code)
	}
}

func TestValidateFold(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test valid fold
	action := Action{PlayerID: 1, Type: ActionFold, Amount: 0}
	err := validator.ValidateAction(game, player, action)
	if err != nil {
		t.Errorf("Unexpected error for valid fold: %v", err)
	}

	// Test fold with amount (should be invalid - fold must have amount 0)
	action = Action{PlayerID: 1, Type: ActionFold, Amount: 100}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for fold with non-zero amount")
	}
	if err.Code != ErrorInvalidAmount {
		t.Errorf("Expected ErrorInvalidAmount, got %d", err.Code)
	}
}

func TestValidateCheck(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test valid check (no current bet)
	action := Action{PlayerID: 1, Type: ActionCheck, Amount: 0}
	err := validator.ValidateAction(game, player, action)
	if err != nil {
		t.Errorf("Unexpected error for valid check: %v", err)
	}

	// Test check with amount (should be invalid)
	action = Action{PlayerID: 1, Type: ActionCheck, Amount: 10}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for check with amount")
	}

	// Create a bet scenario to test invalid check
	game.TakeAction(Action{PlayerID: 2, Type: ActionRaise, Amount: 50})
	action = Action{PlayerID: 1, Type: ActionCheck, Amount: 0}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for check when there's a bet to call")
	}
}

func TestValidateCall(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 1000)
	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 1)

	// Setup a bet to call
	game.TakeAction(Action{PlayerID: 2, Type: ActionRaise, Amount: 50})

	// Test valid call
	action := Action{PlayerID: 1, Type: ActionCall, Amount: 50}
	err := validator.ValidateAction(game, player1, action)
	if err != nil {
		t.Errorf("Unexpected error for valid call: %v", err)
	}

	// Test calling when there's no bet to call
	game = NewGame(10, 20)
	player1 = NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player1, 0)
	action = Action{PlayerID: 1, Type: ActionCall, Amount: 0}
	err = validator.ValidateAction(game, player1, action)
	if err == nil {
		t.Error("Expected error for call when there's no bet to call")
	}

	// Test call with wrong amount
	game = NewGame(10, 20)
	player1 = NewPlayer(1, "Player 1", 1000)
	player2 = NewPlayer(2, "Player 2", 1000)
	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 1)
	game.TakeAction(Action{PlayerID: 2, Type: ActionRaise, Amount: 60})
	action = Action{PlayerID: 1, Type: ActionCall, Amount: 50}
	err = validator.ValidateAction(game, player1, action)
	if err == nil {
		t.Error("Expected error for call with incorrect amount")
	}

	// Test call with insufficient chips - make this player the current player
	poorPlayer := NewPlayer(3, "Poor Player", 25)
	game.PlayerSit(poorPlayer, 2)
	// Fold both player1 and player2 to make poorPlayer the current player
	player1.Fold()
	player2.Fold()
	action = Action{PlayerID: 3, Type: ActionCall, Amount: 60}
	err = validator.ValidateAction(game, poorPlayer, action)
	if err == nil {
		t.Error("Expected error for call with insufficient chips")
	}
	if err.Code != ErrorInsufficientChips {
		t.Errorf("Expected ErrorInsufficientChips, got %d", err.Code)
	}
}

func TestValidateRaise(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test valid raise
	action := Action{PlayerID: 1, Type: ActionRaise, Amount: 40}
	err := validator.ValidateAction(game, player, action)
	if err != nil {
		t.Errorf("Unexpected error for valid raise: %v", err)
	}

	// Test raise with insufficient chips
	action = Action{PlayerID: 1, Type: ActionRaise, Amount: 1500}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for raise with insufficient chips")
	}

	// Test raise below minimum
	action = Action{PlayerID: 1, Type: ActionRaise, Amount: 5}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for raise below minimum")
	}
}

func TestValidateAllIn(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test valid all-in
	action := Action{PlayerID: 1, Type: ActionAllIn, Amount: 1000}
	err := validator.ValidateAction(game, player, action)
	if err != nil {
		t.Errorf("Unexpected error for valid all-in: %v", err)
	}

	// Test all-in with incorrect amount
	action = Action{PlayerID: 1, Type: ActionAllIn, Amount: 500}
	err = validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for all-in with incorrect amount")
	}

	// Test all-in with no chips
	brokePlayer := NewPlayer(2, "Broke Player", 0)
	game.PlayerSit(brokePlayer, 1)
	action = Action{PlayerID: 2, Type: ActionAllIn, Amount: 0}
	err = validator.ValidateAction(game, brokePlayer, action)
	if err == nil {
		t.Error("Expected error for all-in with no chips")
	}
}

func TestGetAvailableActions(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test with no current bet (can check, raise, all-in, fold)
	actions := validator.GetAvailableActions(game, player)
	expectedActions := []ActionType{ActionFold, ActionCheck, ActionRaise, ActionAllIn}

	if len(actions) != len(expectedActions) {
		t.Errorf("Expected %d actions, got %d", len(expectedActions), len(actions))
	}

	for _, expected := range expectedActions {
		found := false
		for _, actual := range actions {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected action %d not found in available actions", expected)
		}
	}

	// Test with folded player (no actions available)
	player.Fold()
	actions = validator.GetAvailableActions(game, player)
	if len(actions) != 0 {
		t.Errorf("Expected no actions for folded player, got %d", len(actions))
	}

	// Test with nil inputs
	actions = validator.GetAvailableActions(nil, player)
	if len(actions) != 0 {
		t.Errorf("Expected no actions for nil game, got %d", len(actions))
	}

	actions = validator.GetAvailableActions(game, nil)
	if len(actions) != 0 {
		t.Errorf("Expected no actions for nil player, got %d", len(actions))
	}
}

func TestGetAvailableActionsWithBet(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 1000)
	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 1)

	// Setup a bet
	game.TakeAction(Action{PlayerID: 2, Type: ActionRaise, Amount: 50})

	// Test available actions with a bet to call
	actions := validator.GetAvailableActions(game, player1)
	expectedActions := []ActionType{ActionFold, ActionCall, ActionRaise, ActionAllIn}

	if len(actions) != len(expectedActions) {
		t.Errorf("Expected %d actions with bet, got %d", len(expectedActions), len(actions))
	}

	// Check that ActionCheck is not available
	for _, action := range actions {
		if action == ActionCheck {
			t.Error("ActionCheck should not be available when there's a bet to call")
		}
	}
}

func TestGetMinRaiseAmount(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test with no current bet (min raise = big blind)
	minRaise := validator.GetMinRaiseAmount(game, player)
	if minRaise != 20 {
		t.Errorf("Expected min raise 20 (big blind), got %d", minRaise)
	}

	// Test with nil inputs
	minRaise = validator.GetMinRaiseAmount(nil, player)
	if minRaise != 0 {
		t.Errorf("Expected 0 for nil game, got %d", minRaise)
	}

	minRaise = validator.GetMinRaiseAmount(game, nil)
	if minRaise != 0 {
		t.Errorf("Expected 0 for nil player, got %d", minRaise)
	}

	// Test with existing bet and previous raise (min raise = currentBet - prevBet)
	game = NewGame(10, 20)
	playerA := NewPlayer(1, "A", 1000)
	playerB := NewPlayer(2, "B", 1000)
	game.PlayerSit(playerA, 0)
	game.PlayerSit(playerB, 1)
	// Player A raises to 40, then Player B raises to 80; min next raise = 80-40 = 40
	game.TakeAction(Action{PlayerID: 1, Type: ActionRaise, Amount: 40})
	game.TakeAction(Action{PlayerID: 2, Type: ActionRaise, Amount: 80})
	got := validator.GetMinRaiseAmount(game, playerA)
	// Compute expected as callAmount + min increment (40)
	currentBet := validator.getCurrentBet(game)
	callAmount := currentBet - playerA.GetBet()
	if callAmount < 0 {
		callAmount = 0
	}
	expected := callAmount + 40
	if got != expected {
		t.Errorf("Expected min raise %d (callAmount %d + 40), got %d", expected, callAmount, got)
	}
}

func TestGetMaxRaiseAmount(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test max raise = player's chips
	maxRaise := validator.GetMaxRaiseAmount(game, player)
	if maxRaise != 1000 {
		t.Errorf("Expected max raise 1000 (player chips), got %d", maxRaise)
	}

	// Test with nil inputs
	maxRaise = validator.GetMaxRaiseAmount(nil, player)
	if maxRaise != 0 {
		t.Errorf("Expected 0 for nil game, got %d", maxRaise)
	}

	maxRaise = validator.GetMaxRaiseAmount(game, nil)
	if maxRaise != 0 {
		t.Errorf("Expected 0 for nil player, got %d", maxRaise)
	}
}

func TestUnknownActionType(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Test unknown action type
	action := Action{PlayerID: 1, Type: ActionType(99), Amount: 0}
	err := validator.ValidateAction(game, player, action)
	if err == nil {
		t.Error("Expected error for unknown action type")
	}
	if err.Code != ErrorInvalidAction {
		t.Errorf("Expected ErrorInvalidAction, got %d", err.Code)
	}
}

func TestValidationErrorCodeToString(t *testing.T) {
	testCases := []struct {
		code     ValidationErrorCode
		expected string
	}{
		{ErrorInvalidPlayer, "Invalid Player"},
		{ErrorInvalidAction, "Invalid Action"},
		{ErrorInsufficientChips, "Insufficient Chips"},
		{ErrorInvalidAmount, "Invalid Amount"},
		{ErrorOutOfTurn, "Out of Turn"},
		{ErrorGameState, "Game State"},
		{ErrorActionNotAllowed, "Action Not Allowed"},
	}

	for _, tc := range testCases {
		result := ValidationErrorCodeToString(tc.code)
		if result != tc.expected {
			t.Errorf("Expected ValidationErrorCodeToString(%d) = %s, got %s", tc.code, tc.expected, result)
		}
	}

	// Test unknown code
	unknown := ValidationErrorCode(99)
	result := ValidationErrorCodeToString(unknown)
	if result != "Unknown" {
		t.Errorf("Expected 'Unknown' for invalid code, got %s", result)
	}
}

func TestComplexGameScenario(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)

	// Setup multiple players
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 500)
	player3 := NewPlayer(3, "Player 3", 100)

	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 1)
	game.PlayerSit(player3, 2)

	// Player 1 raises (current player)
	action := Action{PlayerID: 1, Type: ActionRaise, Amount: 60}
	err := validator.ValidateAction(game, player1, action)
	if err != nil {
		t.Errorf("Unexpected error for player 1 raise: %v", err)
	}
	game.TakeAction(action)

	// Fold player1 to make player2 the current player
	player1.Fold()

	// Player 2 calls
	action = Action{PlayerID: 2, Type: ActionCall, Amount: 60}
	err = validator.ValidateAction(game, player2, action)
	if err != nil {
		t.Errorf("Unexpected error for player 2 call: %v", err)
	}

	// Fold player2 to make player3 the current player
	player2.Fold()

	// Player 3 goes all-in (they're the current player now)
	action = Action{PlayerID: 3, Type: ActionAllIn, Amount: 100}
	err = validator.ValidateAction(game, player3, action)
	if err != nil {
		t.Errorf("Unexpected error for player 3 all-in: %v", err)
	}
}

func TestPlayerBettingHistory(t *testing.T) {
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Player 1", 1000)
	game.PlayerSit(player, 0)

	// Player makes a bet
	player.Bet(30)

	// Test that current bet affects validation
	// This should set up the scenario for proper betting calculations
	action := Action{PlayerID: 1, Type: ActionRaise, Amount: 50}
	err := validator.ValidateAction(game, player, action)
	if err != nil {
		t.Errorf("Unexpected error for raise after bet: %v", err)
	}
}

func TestValidatorAllPhasesActions(t *testing.T) {
	// Test getCurrentPhaseActions across all game phases for better coverage
	validator := NewActionValidator()
	game := NewGame(10, 20)
	player := NewPlayer(1, "Test Player", 1000)
	game.PlayerSit(player, 0)

	// Helper function to add actions to different phases
	addAction := func(phase GamePhase, actionType ActionType) {
		game.SetCurrentPhase(phase)
		action := Action{
			PlayerID: 1,
			Type:     actionType,
			Amount:   10,
		}
		game.TakeAction(action)
	}

	// Add actions to each phase
	addAction(PhasePreflop, ActionRaise)
	addAction(PhaseFlop, ActionCall)
	addAction(PhaseTurn, ActionCheck) // This phase wasn't being tested
	addAction(PhaseRiver, ActionFold) // This phase wasn't being tested

	// Test getCurrentPhaseActions retrieval for all phases
	testCases := []struct {
		phase         GamePhase
		expectedCount int
		expectedType  ActionType
	}{
		{PhasePreflop, 1, ActionRaise},
		{PhaseFlop, 1, ActionCall},
		{PhaseTurn, 1, ActionCheck}, // Now testing Turn phase
		{PhaseRiver, 1, ActionFold}, // Now testing River phase
	}

	for _, tc := range testCases {
		game.SetCurrentPhase(tc.phase)
		actions := validator.getCurrentPhaseActions(game)

		if len(actions) != tc.expectedCount {
			t.Errorf("Expected %d actions in phase %d, got %d", tc.expectedCount, tc.phase, len(actions))
		}

		if len(actions) > 0 && actions[0].Type != tc.expectedType {
			t.Errorf("Expected action type %d in phase %d, got %d", tc.expectedType, tc.phase, actions[0].Type)
		}
	}

	// Test invalid phase (default case)
	game.SetCurrentPhase(GamePhase(99)) // Invalid phase
	actions := validator.getCurrentPhaseActions(game)
	if len(actions) != 0 {
		t.Errorf("Expected 0 actions for invalid phase, got %d", len(actions))
	}

	// Test validateGameState prevents actions during showdown
	game.SetCurrentPhase(PhaseShowdown)
	playerShow := NewPlayer(9, "Show", 1000)
	game.PlayerSit(playerShow, 9)
	action := Action{PlayerID: 9, Type: ActionCheck, Amount: 0}
	if err := validator.ValidateAction(game, playerShow, action); err == nil {
		t.Error("Expected error: no actions allowed during showdown phase")
	}
}

func TestValidateAllInEdgeCases(t *testing.T) {
	// Test edge cases for all-in validation to improve coverage
	validator := NewActionValidator()
	game := NewGame(10, 20)

	// Test player with 0 chips trying to go all-in
	playerZeroChips := NewPlayer(1, "Zero Chips", 0)

	action := Action{
		PlayerID: 1,
		Type:     ActionAllIn,
		Amount:   0,
	}

	// Test the validateAllIn function directly to focus on chip validation
	err := validator.validateAllIn(game, playerZeroChips, action)
	if err == nil {
		t.Error("Expected error when player with 0 chips tries all-in")
	}
	if err.Code != ErrorInsufficientChips {
		t.Errorf("Expected ErrorInsufficientChips, got %d", err.Code)
	}

	// Test all-in with incorrect amount
	playerWithChips := NewPlayer(2, "Has Chips", 100)

	wrongAction := Action{
		PlayerID: 2,
		Type:     ActionAllIn,
		Amount:   50, // Should be 100 (all chips)
	}

	// Test the validateAllIn function directly to focus on amount validation
	err = validator.validateAllIn(game, playerWithChips, wrongAction)
	if err == nil {
		t.Error("Expected error when all-in amount doesn't match chip count")
	}
	if err.Code != ErrorInvalidAmount {
		t.Errorf("Expected ErrorInvalidAmount, got %d", err.Code)
	}

	// Test valid all-in amount
	correctAction := Action{
		PlayerID: 2,
		Type:     ActionAllIn,
		Amount:   100,
	}

	// Test the validateAllIn function directly
	err = validator.validateAllIn(game, playerWithChips, correctAction)
	if err != nil {
		t.Errorf("Unexpected error for valid all-in amount: %v", err)
	}

	// Test player with negative chips (edge case)
	playerNegativeChips := NewPlayer(3, "Negative Chips", -50)
	negativeAction := Action{
		PlayerID: 3,
		Type:     ActionAllIn,
		Amount:   -50,
	}

	err = validator.validateAllIn(game, playerNegativeChips, negativeAction)
	if err == nil {
		t.Error("Expected error for player with negative chips")
	}
	if err.Code != ErrorInsufficientChips {
		t.Errorf("Expected ErrorInsufficientChips for negative chips, got %d", err.Code)
	}
}
