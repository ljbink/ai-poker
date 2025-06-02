package poker

import (
	"testing"
)

func TestBasePlayer(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		playerName   string
		expectedID   int
		expectedName string
	}{
		{
			name:         "Valid player with positive ID",
			id:           1,
			playerName:   "Alice",
			expectedID:   1,
			expectedName: "Alice",
		},
		{
			name:         "Valid player with zero ID",
			id:           0,
			playerName:   "Bob",
			expectedID:   0,
			expectedName: "Bob",
		},
		{
			name:         "Valid player with negative ID",
			id:           -1,
			playerName:   "Charlie",
			expectedID:   -1,
			expectedName: "Charlie",
		},
		{
			name:         "Player with empty name",
			id:           5,
			playerName:   "",
			expectedID:   5,
			expectedName: "",
		},
		{
			name:         "Player with long name",
			id:           100,
			playerName:   "This is a very long player name with spaces and numbers 123",
			expectedID:   100,
			expectedName: "This is a very long player name with spaces and numbers 123",
		},
		{
			name:         "Player with special characters in name",
			id:           42,
			playerName:   "Player@#$%^&*()",
			expectedID:   42,
			expectedName: "Player@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test struct literal creation
			player := BasePlayer{
				ID:   tt.id,
				Name: tt.playerName,
			}

			// Verify ID field
			if player.ID != tt.expectedID {
				t.Errorf("Expected ID %d, got %d", tt.expectedID, player.ID)
			}

			// Verify Name field
			if player.Name != tt.expectedName {
				t.Errorf("Expected Name %q, got %q", tt.expectedName, player.Name)
			}
		})
	}
}

func TestBasePlayer_ZeroValue(t *testing.T) {
	// Test zero value initialization
	var player BasePlayer

	if player.ID != 0 {
		t.Errorf("Expected zero value ID to be 0, got %d", player.ID)
	}

	if player.Name != "" {
		t.Errorf("Expected zero value Name to be empty string, got %q", player.Name)
	}
}

func TestBasePlayer_FieldAssignment(t *testing.T) {
	// Test field assignment after creation
	player := BasePlayer{}

	// Assign ID
	player.ID = 999
	if player.ID != 999 {
		t.Errorf("Expected ID assignment to work, got %d", player.ID)
	}

	// Assign Name
	player.Name = "Assigned Name"
	if player.Name != "Assigned Name" {
		t.Errorf("Expected Name assignment to work, got %q", player.Name)
	}
}

func TestBasePlayer_Copy(t *testing.T) {
	// Test struct copying
	original := BasePlayer{ID: 10, Name: "Original"}
	copy := original

	// Modify copy
	copy.ID = 20
	copy.Name = "Copy"

	// Verify original is unchanged
	if original.ID != 10 {
		t.Errorf("Original ID should be unchanged, got %d", original.ID)
	}
	if original.Name != "Original" {
		t.Errorf("Original Name should be unchanged, got %q", original.Name)
	}

	// Verify copy is changed
	if copy.ID != 20 {
		t.Errorf("Copy ID should be changed, got %d", copy.ID)
	}
	if copy.Name != "Copy" {
		t.Errorf("Copy Name should be changed, got %q", copy.Name)
	}
}
