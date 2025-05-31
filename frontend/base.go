package frontend

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewType represents different screens in the app
type ViewType int

const (
	ViewIndex ViewType = iota
	ViewInputUserProfile
	ViewSettings
	ViewGame
)

// Model represents the main application state
type Model struct {
	currentView          ViewType
	indexView            View
	inputUserProfileView View
	settingsView         View

	width  int
	height int
}

// NewModel creates a new TUI model
func NewModel() *Model {
	model := &Model{
		currentView: ViewIndex,
	}

	// Initialize views with the model reference
	model.indexView = NewIndexView(model)
	model.inputUserProfileView = NewUserProfileView(model)
	model.settingsView = NewSettingsView(model)

	return model
}

// GetPlayerName returns the player name from the centralized data store
func (m *Model) GetPlayerName() string {
	return GetData().GetPlayerName()
}

// Init initializes the model (required by Bubble Tea)
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles all messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// Route to appropriate view handler using interface
		switch m.currentView {
		case ViewIndex:
			return m.indexView.Update(msg)
		case ViewInputUserProfile:
			return m.inputUserProfileView.Update(msg)
		case ViewSettings:
			return m.settingsView.Update(msg)
		case ViewGame:
			return m.updateGameView(msg)
		}
	}

	return m, nil
}

// View renders the current view
func (m *Model) View() string {
	switch m.currentView {
	case ViewIndex:
		return m.indexView.Render(m.width, m.height)
	case ViewInputUserProfile:
		return m.inputUserProfileView.Render(m.width, m.height)
	case ViewSettings:
		return m.settingsView.Render(m.width, m.height)
	case ViewGame:
		return m.renderGameView()
	default:
		return "Unknown view"
	}
}

// RunTUI starts the Bubble Tea application
func RunTUI() error {
	model := NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Common styles
var (
	// Main container style
	containerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")) // Purple

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A78BFA")). // Light purple
			Align(lipgloss.Center).
			Padding(1, 0)

	// Menu item styles
	itemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Margin(0, 1).
			Foreground(lipgloss.Color("#E5E7EB")) // Light gray

	selectedItemStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Margin(0, 1).
				Background(lipgloss.Color("#7C3AED")). // Purple background
				Foreground(lipgloss.Color("#FFFFFF")). // White text
				Bold(true)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")). // Medium gray
			Margin(1, 0)
)

// updateGameView handles updates for the game screen (placeholder for now)
func (m *Model) updateGameView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Go back to index for now
		m.currentView = ViewIndex
	}
	return m, nil
}

// renderGameView renders the game screen (placeholder for now)
func (m *Model) renderGameView() string {
	content := containerStyle.Render(
		titleStyle.Render("🎮 Game View") + "\n\n" +
			"Welcome, " + m.GetPlayerName() + "!\n\n" +
			"Game logic will be implemented here.\n\n" +
			helpStyle.Render("Press Esc to go back to main menu, q to quit"),
	)

	if m.width > 0 {
		content = lipgloss.Place(
			m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	return content
}
