package frontend

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ljbink/ai-poker/frontend/component"
)

// GameSetupKeyMap defines keybindings for the game setup view
type GameSetupKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Continue key.Binding
	Back     key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k GameSetupKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Continue, k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k GameSetupKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Continue},
		{k.Back, k.Quit},
	}
}

var gameSetupKeys = GameSetupKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Continue: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "start game"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// GameSetupView represents the game setup screen
type GameSetupView struct {
	model           *Model
	focused         int // which input field is focused (0=small blind, 1=big blind, 2=num bots)
	smallBlindInput textinput.Model
	bigBlindInput   textinput.Model
	numBotsInput    textinput.Model
	keys            GameSetupKeyMap
	help            help.Model

	// Components
	header *component.HeaderComponent
	helper *component.HelperComponent
}

// NewGameSetupView creates a new game setup view
func NewGameSetupView(model *Model) *GameSetupView {
	// Get current settings
	settings := GetData().GetSettings()

	// Small blind input
	smallBlind := textinput.New()
	smallBlind.Placeholder = "5"
	smallBlind.Width = 15
	smallBlind.Prompt = "$ "
	smallBlind.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")) // Green
	smallBlind.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))
	smallBlind.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA"))
	smallBlind.SetValue(strconv.Itoa(settings.SmallBlind)) // Load from settings
	smallBlind.Focus()

	// Big blind input
	bigBlind := textinput.New()
	bigBlind.Placeholder = "10"
	bigBlind.Width = 15
	bigBlind.Prompt = "$ "
	bigBlind.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")) // Green
	bigBlind.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))
	bigBlind.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA"))
	bigBlind.SetValue(strconv.Itoa(settings.BigBlind)) // Load from settings

	// Number of bots input
	numBots := textinput.New()
	numBots.Placeholder = "3"
	numBots.Width = 15
	numBots.Prompt = "🤖 "
	numBots.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED")) // Purple
	numBots.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))
	numBots.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA"))
	numBots.SetValue(strconv.Itoa(settings.NumBots)) // Load from settings

	// Create help component
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))  // Purple
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")) // Medium gray

	return &GameSetupView{
		model:           model,
		focused:         0,
		smallBlindInput: smallBlind,
		bigBlindInput:   bigBlind,
		numBotsInput:    numBots,
		keys:            gameSetupKeys,
		help:            h,

		// Initialize components with default width (will be updated in Render)
		header: component.NewHeaderComponent("🎲 Game Setup", 80),
		helper: component.NewHelperComponent(gameSetupKeys, 80),
	}
}

// Update handles input for the game setup view
func (v *GameSetupView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, v.keys.Continue):
		if v.validateInputs() {
			// Store game settings
			v.saveGameSettings()
			// Move to game view
			v.model.currentView = ViewGame
			return v.model, nil
		}
	case key.Matches(msg, v.keys.Back):
		// Go back to login
		v.model.currentView = ViewLogin
		return v.model, nil
	case key.Matches(msg, v.keys.Up):
		v.focused--
		if v.focused < 0 {
			v.focused = 2
		}
		v.updateFocus()
	case key.Matches(msg, v.keys.Down):
		v.focused++
		if v.focused > 2 {
			v.focused = 0
		}
		v.updateFocus()
	case key.Matches(msg, v.keys.Quit):
		return v.model, tea.Quit
	}

	// Handle text input updates based on focused field
	switch v.focused {
	case 0:
		v.smallBlindInput, cmd = v.smallBlindInput.Update(msg)
		// Auto-update big blind to be 2x small blind
		if val, err := strconv.Atoi(v.smallBlindInput.Value()); err == nil && val > 0 {
			v.bigBlindInput.SetValue(strconv.Itoa(val * 2))
		}
	case 1:
		v.bigBlindInput, cmd = v.bigBlindInput.Update(msg)
	case 2:
		v.numBotsInput, cmd = v.numBotsInput.Update(msg)
	}

	return v.model, cmd
}

// updateFocus sets focus on the appropriate input field
func (v *GameSetupView) updateFocus() {
	v.smallBlindInput.Blur()
	v.bigBlindInput.Blur()
	v.numBotsInput.Blur()

	switch v.focused {
	case 0:
		v.smallBlindInput.Focus()
	case 1:
		v.bigBlindInput.Focus()
	case 2:
		v.numBotsInput.Focus()
	}
}

// validateInputs checks if all inputs are valid
func (v *GameSetupView) validateInputs() bool {
	smallBlind, err1 := strconv.Atoi(strings.TrimSpace(v.smallBlindInput.Value()))
	bigBlind, err2 := strconv.Atoi(strings.TrimSpace(v.bigBlindInput.Value()))
	numBots, err3 := strconv.Atoi(strings.TrimSpace(v.numBotsInput.Value()))

	return err1 == nil && err2 == nil && err3 == nil &&
		smallBlind > 0 && bigBlind > smallBlind &&
		numBots >= 1 && numBots <= 8
}

// saveGameSettings stores the game configuration
func (v *GameSetupView) saveGameSettings() {
	smallBlind, _ := strconv.Atoi(strings.TrimSpace(v.smallBlindInput.Value()))
	bigBlind, _ := strconv.Atoi(strings.TrimSpace(v.bigBlindInput.Value()))
	numBots, _ := strconv.Atoi(strings.TrimSpace(v.numBotsInput.Value()))

	// Store in centralized data store (we might need to add these methods)
	data := GetData()
	data.UpdateSetting("small_blind", smallBlind)
	data.UpdateSetting("big_blind", bigBlind)
	data.UpdateSetting("num_bots", numBots)
}

// Render renders the game setup view
func (v *GameSetupView) Render(width, height int) string {
	// Update component widths for current screen size
	v.header.SetWidth(width)
	v.helper.SetWidth(width)

	var b strings.Builder

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D1D5DB")).
		Render("Configure your poker game:")
	b.WriteString(instructions)
	b.WriteString("\n\n")

	// Small Blind section
	smallBlindLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5E7EB")).
		Bold(true).
		Render("Small Blind:")
	b.WriteString(smallBlindLabel)
	b.WriteString("\n")

	smallBlindBox := v.createInputBox(v.smallBlindInput, v.focused == 0)
	b.WriteString(smallBlindBox)
	b.WriteString("\n\n")

	// Big Blind section
	bigBlindLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5E7EB")).
		Bold(true).
		Render("Big Blind:")
	b.WriteString(bigBlindLabel)
	b.WriteString("\n")

	bigBlindBox := v.createInputBox(v.bigBlindInput, v.focused == 1)
	b.WriteString(bigBlindBox)
	b.WriteString("\n\n")

	// Number of Bots section
	numBotsLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5E7EB")).
		Bold(true).
		Render("Number of Bots (1-8):")
	b.WriteString(numBotsLabel)
	b.WriteString("\n")

	numBotsBox := v.createInputBox(v.numBotsInput, v.focused == 2)
	b.WriteString(numBotsBox)
	b.WriteString("\n\n")

	// Validation status
	if v.validateInputs() {
		statusMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")). // Green
			Render("✓ Ready to start game")
		b.WriteString(statusMsg)
	} else {
		statusMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")). // Red
			Render("⚠ Please check your input values")
		b.WriteString(statusMsg)
	}

	// Title at the top using header component
	titleAtTop := v.header.Render()

	// Help view at the bottom using helper component
	helpAtBottom := v.helper.Render()

	// Calculate actual space used by header and helper
	headerHeight := lipgloss.Height(titleAtTop)
	helperHeight := lipgloss.Height(helpAtBottom)
	availableHeight := height - headerHeight - helperHeight

	// Center the form content in the middle of available space
	content := b.String()
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

// createInputBox creates a styled input box
func (v *GameSetupView) createInputBox(input textinput.Model, focused bool) string {
	borderColor := "#6B7280" // Gray
	if focused {
		borderColor = "#7C3AED" // Purple when focused
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Render(input.View())
}

// GetType returns the view type
func (v *GameSetupView) GetType() ViewType {
	return ViewGameSetup
}

// ShortHelp returns keybindings to be shown in the mini help view
func (v *GameSetupView) ShortHelp() []key.Binding {
	return v.keys.ShortHelp()
}

// FullHelp returns keybindings for the expanded help view
func (v *GameSetupView) FullHelp() [][]key.Binding {
	return v.keys.FullHelp()
}
