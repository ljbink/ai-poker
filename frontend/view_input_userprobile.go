package frontend

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UserProfileView represents the user profile input screen
type UserProfileView struct {
	model      *Model
	playerName string
	focused    bool
}

// NewUserProfileView creates a new user profile input view
func NewUserProfileView(model *Model) *UserProfileView {
	return &UserProfileView{
		model:      model,
		playerName: "",
		focused:    true,
	}
}

// Update handles input for the user profile view
func (v *UserProfileView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if strings.TrimSpace(v.playerName) != "" {
			// Store the player name in the centralized data store
			GetData().SetPlayerName(v.playerName)
			// Move to game view once we have a player name
			v.model.currentView = ViewGame
		}
	case "esc":
		// Go back to index
		v.model.currentView = ViewIndex
	case "backspace":
		if len(v.playerName) > 0 {
			v.playerName = v.playerName[:len(v.playerName)-1]
		}
	default:
		// Add character to player name (basic input handling)
		if len(msg.String()) == 1 && len(v.playerName) < 20 {
			char := msg.String()
			// Allow letters, numbers, and spaces
			if (char >= "a" && char <= "z") ||
				(char >= "A" && char <= "Z") ||
				(char >= "0" && char <= "9") ||
				char == " " {
				v.playerName += char
			}
		}
	}
	return v.model, nil
}

// Render renders the user profile input view
func (v *UserProfileView) Render(width, height int) string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("Player Setup")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Instructions
	instructions := "Enter your player name:"
	b.WriteString(instructions)
	b.WriteString("\n\n")

	// Input field
	inputContent := v.playerName
	if v.focused {
		inputContent += "█" // Cursor
	}

	var inputBox string
	if v.focused {
		inputBox = focusedInputStyle.Render(inputContent)
	} else {
		inputBox = inputStyle.Render(inputContent)
	}

	b.WriteString(inputBox)
	b.WriteString("\n\n")

	// Instructions and help
	var helpLines []string

	if strings.TrimSpace(v.playerName) != "" {
		helpLines = append(helpLines, " • Press Enter to continue")
	} else {
		helpLines = append(helpLines, "")
	}
	helpLines = append(helpLines, " • Press Esc to go back")
	helpLines = append(helpLines, " • Press q to quit")
	helpText := helpStyle.Render(strings.Join(helpLines, "\n"))
	b.WriteString(helpText)

	// Center the content
	content := b.String()
	if width > 0 {
		content = lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			containerStyle.Render(content),
		)
	}

	return content
}

// GetType returns the view type
func (v *UserProfileView) GetType() ViewType {
	return ViewInputUserProfile
}

// GetPlayerName returns the entered player name
func (v *UserProfileView) GetPlayerName() string {
	return v.playerName
}
