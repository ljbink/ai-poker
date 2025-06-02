package frontend

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ljbink/ai-poker/frontend/component"
)

// GameKeyMap defines keybindings for the game view
type GameKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k GameKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k GameKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Back, k.Quit},
	}
}

var gameKeys = GameKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// GameView represents the game screen
type GameView struct {
	model *Model
	keys  GameKeyMap
	help  help.Model

	// Components
	header *component.HeaderComponent
	helper *component.HelperComponent
}

// NewGameView creates a new game view
func NewGameView(model *Model) *GameView {
	// Create help component with matching styling
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))  // Purple
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")) // Medium gray

	return &GameView{
		model: model,
		keys:  gameKeys,
		help:  h,

		// Initialize components with default width (will be updated in Render)
		header: component.NewHeaderComponent("🎮 Game View", 80),
		helper: component.NewHelperComponent(gameKeys, 80),
	}
}

// Update handles input for the game view
func (v *GameView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, v.keys.Back):
		// Go back to index
		v.model.currentView = ViewIndex
	case key.Matches(msg, v.keys.Quit):
		return v.model, tea.Quit
	}
	return v.model, nil
}

// Render renders the game view
func (v *GameView) Render(width, height int) string {
	// Update component widths for current screen size
	v.header.SetWidth(width)
	v.helper.SetWidth(width)

	content := "Welcome, " + GetData().GetPlayerName() + "!\n\n" +
		"Game logic will be implemented here."

	// Title at the top using header component
	titleAtTop := v.header.Render()

	// Help view at the bottom using helper component
	helpAtBottom := v.helper.Render()

	// Calculate actual space used by header and helper
	headerHeight := lipgloss.Height(titleAtTop)
	helperHeight := lipgloss.Height(helpAtBottom)
	availableHeight := height - headerHeight - helperHeight

	// Center the game content in the middle of available space
	centeredContent := lipgloss.Place(
		width, availableHeight,
		lipgloss.Center, lipgloss.Center,
		content,
	)

	// Combine title, content, and help without extra spacing
	fullContent := titleAtTop + centeredContent + helpAtBottom

	// Apply full screen style
	fullScreenContainer := GetFullScreenStyle(width, height)
	return fullScreenContainer.Render(fullContent)
}

// GetType returns the view type
func (v *GameView) GetType() ViewType {
	return ViewGame
}

// ShortHelp returns keybindings to be shown in the mini help view
func (v *GameView) ShortHelp() []key.Binding {
	return v.keys.ShortHelp()
}

// FullHelp returns keybindings for the expanded help view
func (v *GameView) FullHelp() [][]key.Binding {
	return v.keys.FullHelp()
}
