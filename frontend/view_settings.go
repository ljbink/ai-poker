package frontend

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ljbink/ai-poker/frontend/component"
)

// SettingsKeyMap defines keybindings for the settings view
type SettingsKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Left   key.Binding
	Right  key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k SettingsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Left, k.Right, k.Back}
}

// FullHelp returns keybindings for the expanded help view.
func (k SettingsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select},
		{k.Left, k.Right, k.Back, k.Quit},
	}
}

var settingsKeys = SettingsKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "toggle"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "decrease"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "increase"),
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

// SettingsView represents the settings screen
type SettingsView struct {
	model    *Model
	selected int
	options  []SettingOption
	keys     SettingsKeyMap
	help     help.Model

	// Components
	header *component.HeaderComponent
	helper *component.HelperComponent
}

// SettingOption represents a configurable setting
type SettingOption struct {
	Label       string
	Key         string
	ValueType   string // "bool", "int", "string"
	Description string
	Icon        string
}

// NewSettingsView creates a new settings view
func NewSettingsView(model *Model) *SettingsView {
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))  // Purple
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")) // Medium gray

	return &SettingsView{
		model:    model,
		selected: 0,
		keys:     settingsKeys,
		help:     h,

		// Initialize components with default width (will be updated in Render)
		header: component.NewHeaderComponent("⚙️  Settings", 80),
		helper: component.NewHelperComponent(settingsKeys, 80),
		options: []SettingOption{
			{
				Label:       "Theme",
				Key:         "theme",
				ValueType:   "string",
				Description: "Application theme (dark/light/auto)",
				Icon:        "🎨",
			},
			{
				Label:       "Sound Effects",
				Key:         "sound_enabled",
				ValueType:   "bool",
				Description: "Enable sound effects",
				Icon:        "🔊",
			},
			{
				Label:       "Animations",
				Key:         "animations_enabled",
				ValueType:   "bool",
				Description: "Enable UI animations",
				Icon:        "✨",
			},
			{
				Label:       "Auto Save",
				Key:         "auto_save",
				ValueType:   "bool",
				Description: "Automatically save game progress",
				Icon:        "💾",
			},
			{
				Label:       "Default Buy-in",
				Key:         "default_buy_in",
				ValueType:   "int",
				Description: "Default chip amount when starting a game",
				Icon:        "💰",
			},
			{
				Label:       "Show Probabilities",
				Key:         "show_probabilities",
				ValueType:   "bool",
				Description: "Display hand probability information",
				Icon:        "📊",
			},
		},
	}
}

// Update handles input for the settings view
func (v *SettingsView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, v.keys.Up):
		if v.selected > 0 {
			v.selected--
		}
	case key.Matches(msg, v.keys.Down):
		if v.selected < len(v.options)-1 {
			v.selected++
		}
	case key.Matches(msg, v.keys.Select):
		// Toggle or modify the selected setting
		v.toggleSetting(v.selected)
	case key.Matches(msg, v.keys.Back):
		// Go back to index
		v.model.currentView = ViewIndex
	case key.Matches(msg, v.keys.Left):
		// Decrease value for numeric settings
		v.adjustSetting(v.selected, -1)
	case key.Matches(msg, v.keys.Right):
		// Increase value for numeric settings
		v.adjustSetting(v.selected, 1)
	case key.Matches(msg, v.keys.Quit):
		return v.model, tea.Quit
	}
	return v.model, nil
}

// Render renders the settings view
func (v *SettingsView) Render(width, height int) string {
	// Update component widths for current screen size
	v.header.SetWidth(width)
	v.helper.SetWidth(width)

	var b strings.Builder

	// Get current settings
	settings := GetData().GetSettings()

	// Settings options with enhanced styling
	for i, option := range v.options {
		var line string
		var currentValue string
		var valueStyle lipgloss.Style

		// Get current value based on setting key
		switch option.Key {
		case "theme":
			currentValue = settings.Theme
			valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA")).Bold(true) // Light purple
		case "sound_enabled":
			if settings.SoundEnabled {
				currentValue = "✓ enabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")) // Green
			} else {
				currentValue = "✗ disabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")) // Red
			}
		case "animations_enabled":
			if settings.AnimationsEnabled {
				currentValue = "✓ enabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")) // Green
			} else {
				currentValue = "✗ disabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")) // Red
			}
		case "auto_save":
			if settings.AutoSave {
				currentValue = "✓ enabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")) // Green
			} else {
				currentValue = "✗ disabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")) // Red
			}
		case "default_buy_in":
			currentValue = fmt.Sprintf("%d chips", settings.DefaultBuyIn)
			valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B")) // Yellow/Orange
		case "show_probabilities":
			if settings.ShowProbabilities {
				currentValue = "✓ enabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")) // Green
			} else {
				currentValue = "✗ disabled"
				valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")) // Red
			}
		}

		// Format the line with icon
		line = fmt.Sprintf("%s %-18s: %s", option.Icon, option.Label, valueStyle.Render(currentValue))

		if i == v.selected {
			// Selected item styling with border
			selectedStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")). // White text
				Background(lipgloss.Color("#7C3AED")). // Purple background
				Padding(0, 1).
				Bold(true)
			b.WriteString(selectedStyle.Render("▶ " + line))
			b.WriteString("\n")
		} else {
			b.WriteString(itemStyle.Render("  " + line))
			b.WriteString("\n")
		}

		// Show description for selected item
		description := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")). // Medium gray
			Italic(true).
			Render("  " + option.Description)
		b.WriteString(description)
		b.WriteString("\n")
		b.WriteString("\n")
	}

	// Title at the top using header component
	titleAtTop := v.header.Render()

	// Help view at the bottom using helper component
	helpAtBottom := v.helper.Render()

	// Calculate actual space used by header and helper
	headerHeight := lipgloss.Height(titleAtTop)
	helperHeight := lipgloss.Height(helpAtBottom)
	availableHeight := height - headerHeight - helperHeight

	// Center the settings content in the middle of available space
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
func (v *SettingsView) GetType() ViewType {
	return ViewSettings
}

// ShortHelp returns keybindings to be shown in the mini help view
func (v *SettingsView) ShortHelp() []key.Binding {
	return v.keys.ShortHelp()
}

// FullHelp returns keybindings for the expanded help view
func (v *SettingsView) FullHelp() [][]key.Binding {
	return v.keys.FullHelp()
}

// toggleSetting toggles a boolean setting or cycles string settings
func (v *SettingsView) toggleSetting(index int) {
	if index >= 0 && index < len(v.options) {
		option := v.options[index]
		settings := GetData().GetSettings()

		switch option.Key {
		case "theme":
			// Cycle through themes
			switch settings.Theme {
			case "dark":
				GetData().UpdateSetting("theme", "light")
			case "light":
				GetData().UpdateSetting("theme", "auto")
			case "auto":
				GetData().UpdateSetting("theme", "dark")
			default:
				GetData().UpdateSetting("theme", "dark")
			}
		case "sound_enabled":
			GetData().UpdateSetting("sound_enabled", !settings.SoundEnabled)
		case "animations_enabled":
			GetData().UpdateSetting("animations_enabled", !settings.AnimationsEnabled)
		case "auto_save":
			GetData().UpdateSetting("auto_save", !settings.AutoSave)
		case "show_probabilities":
			GetData().UpdateSetting("show_probabilities", !settings.ShowProbabilities)
		}
	}
}

// adjustSetting adjusts numeric settings
func (v *SettingsView) adjustSetting(index int, delta int) {
	if index >= 0 && index < len(v.options) {
		option := v.options[index]
		if option.ValueType == "int" && option.Key == "default_buy_in" {
			settings := GetData().GetSettings()
			newValue := settings.DefaultBuyIn + (delta * 100) // Adjust by 100 chips
			if newValue >= 100 && newValue <= 10000 {         // Reasonable limits
				GetData().UpdateSetting("default_buy_in", newValue)
			}
		}
	}
}
