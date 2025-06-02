package holdem

import (
	"testing"
)

func TestActionTypes(t *testing.T) {
	// Test all action type constants
	expectedActions := []ActionType{
		ActionFold,
		ActionCheck,
		ActionCall,
		ActionRaise,
		ActionAllIn,
	}

	// Verify they have different values
	seen := make(map[ActionType]bool)
	for _, action := range expectedActions {
		if seen[action] {
			t.Errorf("Duplicate action type value: %d", action)
		}
		seen[action] = true
	}

	// Verify they are sequential (starting from 0)
	for i, action := range expectedActions {
		if int(action) != i {
			t.Errorf("Expected action type %d to have value %d, got %d", i, i, int(action))
		}
	}
}

func TestActionStruct(t *testing.T) {
	// Test creating an action
	playerID := 123
	actionType := ActionRaise
	amount := 100

	action := Action{
		PlayerID: playerID,
		Type:     actionType,
		Amount:   amount,
	}

	if action.PlayerID != playerID {
		t.Errorf("Expected PlayerID %d, got %d", playerID, action.PlayerID)
	}
	if action.Type != actionType {
		t.Errorf("Expected Type %d, got %d", actionType, action.Type)
	}
	if action.Amount != amount {
		t.Errorf("Expected Amount %d, got %d", amount, action.Amount)
	}
}

func TestActionStructZeroValues(t *testing.T) {
	// Test zero value action
	action := Action{}

	if action.PlayerID != 0 {
		t.Errorf("Expected zero PlayerID, got %d", action.PlayerID)
	}
	if action.Type != ActionFold {
		t.Errorf("Expected zero Type (ActionFold), got %d", action.Type)
	}
	if action.Amount != 0 {
		t.Errorf("Expected zero Amount, got %d", action.Amount)
	}
}

func TestActionWithNegativeValues(t *testing.T) {
	// Test action with negative values (should be allowed structurally)
	action := Action{
		PlayerID: -1,
		Type:     ActionType(-1),
		Amount:   -100,
	}

	if action.PlayerID != -1 {
		t.Errorf("Expected PlayerID -1, got %d", action.PlayerID)
	}
	if action.Type != ActionType(-1) {
		t.Errorf("Expected Type -1, got %d", action.Type)
	}
	if action.Amount != -100 {
		t.Errorf("Expected Amount -100, got %d", action.Amount)
	}
}

func TestIsValidActionType(t *testing.T) {
	// Test valid player action types
	validPlayerActions := []ActionType{
		ActionFold,
		ActionCheck,
		ActionCall,
		ActionRaise,
		ActionAllIn,
	}

	for _, action := range validPlayerActions {
		if !IsValidActionType(action) {
			t.Errorf("Expected player action type %d to be valid", action)
		}
	}

	// Test valid system action types
	validSystemActions := []ActionType{
		ActionSystemShuffle,
		ActionSystemDealHole,
		ActionSystemDealFlop,
		ActionSystemDealTurn,
		ActionSystemDealRiver,
		ActionSystemPhaseChange,
	}

	for _, action := range validSystemActions {
		if !IsValidActionType(action) {
			t.Errorf("Expected system action type %d to be valid", action)
		}
	}

	// Test invalid action types (values beyond our defined constants)
	invalidActions := []ActionType{
		ActionType(-1),
		ActionType(20),
		ActionType(50),
		ActionType(100),
	}

	for _, action := range invalidActions {
		if IsValidActionType(action) {
			t.Errorf("Expected action type %d to be invalid", action)
		}
	}
}

func TestActionTypeToString(t *testing.T) {
	testCases := []struct {
		actionType ActionType
		expected   string
	}{
		// Player actions
		{ActionFold, "Fold"},
		{ActionCheck, "Check"},
		{ActionCall, "Call"},
		{ActionRaise, "Raise"},
		{ActionAllIn, "All-In"},
		// System actions
		{ActionSystemShuffle, "System: Shuffle"},
		{ActionSystemDealHole, "System: Deal Hole Cards"},
		{ActionSystemDealFlop, "System: Deal Flop"},
		{ActionSystemDealTurn, "System: Deal Turn"},
		{ActionSystemDealRiver, "System: Deal River"},
		{ActionSystemPhaseChange, "System: Phase Change"},
	}

	for _, tc := range testCases {
		result := ActionTypeToString(tc.actionType)
		if result != tc.expected {
			t.Errorf("Expected ActionTypeToString(%d) = %s, got %s", tc.actionType, tc.expected, result)
		}
	}

	// Test invalid action type
	invalidAction := ActionType(99)
	result := ActionTypeToString(invalidAction)
	expected := "Unknown"
	if result != expected {
		t.Errorf("Expected ActionTypeToString(%d) = %s, got %s", invalidAction, expected, result)
	}
}

func TestActionStructEquality(t *testing.T) {
	// Test action equality
	action1 := Action{PlayerID: 1, Type: ActionRaise, Amount: 100}
	action2 := Action{PlayerID: 1, Type: ActionRaise, Amount: 100}
	action3 := Action{PlayerID: 2, Type: ActionRaise, Amount: 100}

	if action1 != action2 {
		t.Error("Expected identical actions to be equal")
	}

	if action1 == action3 {
		t.Error("Expected actions with different PlayerIDs to be unequal")
	}
}

func TestActionStructCopy(t *testing.T) {
	// Test copying actions
	original := Action{PlayerID: 1, Type: ActionRaise, Amount: 100}
	copy := original

	// Modify copy
	copy.PlayerID = 2
	copy.Amount = 200

	// Original should be unchanged
	if original.PlayerID != 1 {
		t.Errorf("Expected original PlayerID to remain 1, got %d", original.PlayerID)
	}
	if original.Amount != 100 {
		t.Errorf("Expected original Amount to remain 100, got %d", original.Amount)
	}

	// Copy should be changed
	if copy.PlayerID != 2 {
		t.Errorf("Expected copy PlayerID to be 2, got %d", copy.PlayerID)
	}
	if copy.Amount != 200 {
		t.Errorf("Expected copy Amount to be 200, got %d", copy.Amount)
	}
}

func TestActionSlice(t *testing.T) {
	// Test working with slices of actions
	actions := []Action{
		{PlayerID: 1, Type: ActionFold, Amount: 0},
		{PlayerID: 2, Type: ActionCall, Amount: 50},
		{PlayerID: 3, Type: ActionRaise, Amount: 100},
	}

	if len(actions) != 3 {
		t.Errorf("Expected 3 actions in slice, got %d", len(actions))
	}

	// Verify each action
	if actions[0].PlayerID != 1 || actions[0].Type != ActionFold {
		t.Error("First action not as expected")
	}
	if actions[1].PlayerID != 2 || actions[1].Type != ActionCall {
		t.Error("Second action not as expected")
	}
	if actions[2].PlayerID != 3 || actions[2].Type != ActionRaise {
		t.Error("Third action not as expected")
	}

	// Test appending
	newAction := Action{PlayerID: 4, Type: ActionAllIn, Amount: 200}
	actions = append(actions, newAction)

	if len(actions) != 4 {
		t.Errorf("Expected 4 actions after append, got %d", len(actions))
	}
	if actions[3].PlayerID != 4 || actions[3].Type != ActionAllIn {
		t.Error("Appended action not as expected")
	}
}

func TestSystemActionsInGame(t *testing.T) {
	// Test that system actions are properly logged during game operations
	game := NewGame(10, 20)
	player1 := NewPlayer(1, "Player 1", 1000)
	player2 := NewPlayer(2, "Player 2", 1000)

	// Add players
	game.PlayerSit(player1, 0)
	game.PlayerSit(player2, 1)

	// Deal hole cards (should log system actions)
	err := game.DealHoleCards()
	if err != nil {
		t.Errorf("Unexpected error dealing hole cards: %v", err)
	}

	// Check system actions were logged
	systemActions := game.GetSystemActions()
	if len(systemActions.Preflop) == 0 {
		t.Error("Expected system actions to be logged for hole card dealing")
	}

	// Verify the logged actions
	foundShuffle := false
	foundDealHole := false
	for _, action := range systemActions.Preflop {
		if action.Type == ActionSystemShuffle {
			foundShuffle = true
			if action.PlayerID != SystemPlayerID {
				t.Errorf("Expected system action to have SystemPlayerID (%d), got %d", SystemPlayerID, action.PlayerID)
			}
		}
		if action.Type == ActionSystemDealHole {
			foundDealHole = true
			if action.Amount != 4 { // 2 players * 2 cards each
				t.Errorf("Expected deal hole action to have amount 4, got %d", action.Amount)
			}
		}
	}

	if !foundShuffle {
		t.Error("Expected shuffle system action to be logged")
	}
	if !foundDealHole {
		t.Error("Expected deal hole cards system action to be logged")
	}

	// Test phase change
	game.SetCurrentPhase(PhaseFlop)

	// Should log phase change
	if len(systemActions.Preflop) == 0 {
		t.Error("Expected phase change to be logged")
	}

	// Deal flop
	err = game.DealFlop()
	if err != nil {
		t.Errorf("Unexpected error dealing flop: %v", err)
	}

	systemActions = game.GetSystemActions()
	if len(systemActions.Flop) == 0 {
		t.Error("Expected flop system actions to be logged")
	}

	// Verify flop action
	foundFlopDeal := false
	for _, action := range systemActions.Flop {
		if action.Type == ActionSystemDealFlop {
			foundFlopDeal = true
			if action.Amount != 3 { // 3 flop cards
				t.Errorf("Expected deal flop action to have amount 3, got %d", action.Amount)
			}
		}
	}

	if !foundFlopDeal {
		t.Error("Expected deal flop system action to be logged")
	}

	// Verify that user actions are still separate
	userActions := game.GetUserActions()
	if len(userActions.Preflop) != 0 {
		t.Error("Expected no user actions yet")
	}

	// Take a user action
	playerAction := Action{
		PlayerID: 1,
		Type:     ActionCheck,
		Amount:   0,
	}
	err = game.TakeAction(playerAction)
	if err != nil {
		t.Errorf("Unexpected error taking user action: %v", err)
	}

	// Check that user action was logged separately
	userActions = game.GetUserActions()
	if len(userActions.Flop) != 1 {
		t.Errorf("Expected 1 user action in flop, got %d", len(userActions.Flop))
	}

	if userActions.Flop[0].Type != ActionCheck {
		t.Errorf("Expected user action to be check, got %d", userActions.Flop[0].Type)
	}

	// Verify system actions weren't affected by user action
	systemActions = game.GetSystemActions()
	systemActionCount := len(systemActions.Preflop) + len(systemActions.Flop)
	if systemActionCount < 3 { // Should have shuffle, deal hole, deal flop at minimum
		t.Errorf("Expected at least 3 system actions, got %d", systemActionCount)
	}
}
