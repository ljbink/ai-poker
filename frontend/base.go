package frontend

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewType represents different screens in the app
type ViewType int

const (
	ViewIndex ViewType = iota
	ViewLogin
	ViewGameSetup
	ViewSettings
	ViewGame
)

// Model represents the main application state
type Model struct {
	currentView   ViewType
	indexView     View
	loginView     View
	gameSetupView View
	settingsView  View
	gameView      View

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
	model.loginView = NewLoginView(model)
	model.gameSetupView = NewGameSetupView(model)
	model.settingsView = NewSettingsView(model)
	model.gameView = NewGameView(model)

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
		case ViewLogin:
			return m.loginView.Update(msg)
		case ViewGameSetup:
			return m.gameSetupView.Update(msg)
		case ViewSettings:
			return m.settingsView.Update(msg)
		case ViewGame:
			return m.gameView.Update(msg)
		}
	}

	return m, nil
}

// View renders the current view
func (m *Model) View() string {
	switch m.currentView {
	case ViewIndex:
		return m.indexView.Render(m.width, m.height)
	case ViewLogin:
		return m.loginView.Render(m.width, m.height)
	case ViewGameSetup:
		return m.gameSetupView.Render(m.width, m.height)
	case ViewSettings:
		return m.settingsView.Render(m.width, m.height)
	case ViewGame:
		return m.gameView.Render(m.width, m.height)
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
	// Menu item styles
	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB")) // Light gray

	selectedItemStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7C3AED")). // Purple background
				Foreground(lipgloss.Color("#FFFFFF")). // White text
				Bold(true)
)

// GetFullScreenStyle returns a style configured for the given dimensions
func GetFullScreenStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height)
}
