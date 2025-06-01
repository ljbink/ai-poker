package frontend

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ljbink/ai-poker/frontend/component"
)

// LoginKeyMap defines keybindings for the login view
type LoginKeyMap struct {
	Continue key.Binding
	Back     key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k LoginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Continue, k.Back, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k LoginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Continue, k.Back},
		{k.Quit},
	}
}

var loginKeys = LoginKeyMap{
	Continue: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "continue"),
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

// LoginView represents the login screen
type LoginView struct {
	model     *Model
	textInput textinput.Model
	keys      LoginKeyMap
	help      help.Model

	// Components
	header *component.HeaderComponent
	helper *component.HelperComponent
}

// NewLoginView creates a new login view
func NewLoginView(model *Model) *LoginView {
	ti := textinput.New()
	ti.Placeholder = "Enter your name..."
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 30
	ti.Prompt = "➤ "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F4F6"))
	ti.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA"))

	// Create help component with matching SettingsView styling
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))  // Purple
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")) // Medium gray

	return &LoginView{
		model:     model,
		textInput: ti,
		keys:      loginKeys,
		help:      h,

		// Initialize components with default width (will be updated in Render)
		header: component.NewHeaderComponent("🔑 Login", 80),
		helper: component.NewHelperComponent(loginKeys, 80),
	}
}

// Update handles input for the login view
func (v *LoginView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, v.keys.Continue):
		if strings.TrimSpace(v.textInput.Value()) != "" {
			// Store the player name in the centralized data store
			GetData().SetPlayerName(v.textInput.Value())
			// Move to game setup view to configure the game
			v.model.currentView = ViewGameSetup
			return v.model, nil
		}
	case key.Matches(msg, v.keys.Back):
		// Go back to index
		v.model.currentView = ViewIndex
		return v.model, nil
	case key.Matches(msg, v.keys.Quit):
		return v.model, tea.Quit
	}

	// Handle textinput updates
	v.textInput, cmd = v.textInput.Update(msg)
	return v.model, cmd
}

// Render renders the login view
func (v *LoginView) Render(width, height int) string {
	// Update component widths for current screen size
	v.header.SetWidth(width)
	v.helper.SetWidth(width)

	var b strings.Builder

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

	// Status message
	if strings.TrimSpace(v.textInput.Value()) != "" {
		statusMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")). // Green
			Render("✓ Ready to continue")
		b.WriteString(statusMsg)
	} else {
		statusMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")). // Medium gray
			Render("Enter your player name")
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

// GetType returns the view type
func (v *LoginView) GetType() ViewType {
	return ViewLogin
}

// ShortHelp returns keybindings to be shown in the mini help view
func (v *LoginView) ShortHelp() []key.Binding {
	return v.keys.ShortHelp()
}

// FullHelp returns keybindings for the expanded help view
func (v *LoginView) FullHelp() [][]key.Binding {
	return v.keys.FullHelp()
}
