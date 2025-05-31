package frontend

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// UserProfileView represents the user profile input screen
type UserProfileView struct {
	model     *Model
	textInput textinput.Model
}

// NewUserProfileView creates a new user profile input view
func NewUserProfileView(model *Model) *UserProfileView {
	ti := textinput.New()
	ti.Placeholder = "Enter your name..."
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 30
	ti.Prompt = "➤ "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))
	ti.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA"))

	return &UserProfileView{
		model:     model,
		textInput: ti,
	}
}

// Update handles input for the user profile view
func (v *UserProfileView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		if strings.TrimSpace(v.textInput.Value()) != "" {
			// Store the player name in the centralized data store
			GetData().SetPlayerName(v.textInput.Value())
			// Move to game view once we have a player name
			v.model.currentView = ViewGame
			return v.model, nil
		}
	case "esc":
		// Go back to index
		v.model.currentView = ViewIndex
		return v.model, nil
	}

	// Handle textinput updates
	v.textInput, cmd = v.textInput.Update(msg)
	return v.model, cmd
}

// Render renders the user profile input view
func (v *UserProfileView) Render(width, height int) string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("🎭 Player Setup")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D1D5DB")).
		Render("What should we call you?")
	b.WriteString(instructions)
	b.WriteString("\n\n")

	// Text input field
	inputBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		Padding(0, 1).
		Render(v.textInput.View())

	b.WriteString(inputBox)
	b.WriteString("\n\n")

	// Instructions and help
	var helpLines []string

	if strings.TrimSpace(v.textInput.Value()) != "" {
		helpLines = append(helpLines, "✓ Press Enter to continue")
	} else {
		helpLines = append(helpLines, "")
	}
	helpLines = append(helpLines, "• Press Esc to go back")
	helpLines = append(helpLines, "• Press q to quit")

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

// ShortHelp returns keybindings to be shown in the mini help view
func (v *UserProfileView) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "continue"),
		),
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// FullHelp returns keybindings for the expanded help view
func (v *UserProfileView) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "continue"),
			),
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "back"),
			),
		},
		{
			key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
	}
}
