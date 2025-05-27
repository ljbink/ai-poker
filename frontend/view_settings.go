package frontend

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsView represents the settings screen
type SettingsView struct {
	model    *Model
	selected int
	options  []SettingOption
}

// SettingOption represents a configurable setting
type SettingOption struct {
	Label       string
	Key         string
	ValueType   string // "bool", "int", "string"
	Description string
}

// NewSettingsView creates a new settings view
func NewSettingsView(model *Model) *SettingsView {
	return &SettingsView{
		model:    model,
		selected: 0,
		options: []SettingOption{
			{
				Label:       "Theme",
				Key:         "theme",
				ValueType:   "string",
				Description: "Application theme (dark/light)",
			},
			{
				Label:       "Sound Effects",
				Key:         "sound_enabled",
				ValueType:   "bool",
				Description: "Enable sound effects",
			},
			{
				Label:       "Animations",
				Key:         "animations_enabled",
				ValueType:   "bool",
				Description: "Enable UI animations",
			},
			{
				Label:       "Auto Save",
				Key:         "auto_save",
				ValueType:   "bool",
				Description: "Automatically save game progress",
			},
			{
				Label:       "Default Buy-in",
				Key:         "default_buy_in",
				ValueType:   "int",
				Description: "Default chip amount when starting a game",
			},
			{
				Label:       "Show Probabilities",
				Key:         "show_probabilities",
				ValueType:   "bool",
				Description: "Display hand probability information",
			},
		},
	}
}

// Update handles input for the settings view
func (v *SettingsView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if v.selected > 0 {
			v.selected--
		}
	case "down", "j":
		if v.selected < len(v.options)-1 {
			v.selected++
		}
	case "enter", " ":
		// Toggle or modify the selected setting
		v.toggleSetting(v.selected)
	case "esc":
		// Go back to index
		v.model.currentView = ViewIndex
	case "left", "h":
		// Decrease value for numeric settings
		v.adjustSetting(v.selected, -1)
	case "right", "l":
		// Increase value for numeric settings
		v.adjustSetting(v.selected, 1)
	}
	return v.model, nil
}

// Render renders the settings view
func (v *SettingsView) Render(width, height int) string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("⚙️  Settings")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Get current settings
	settings := GetData().GetSettings()

	// Settings options
	for i, option := range v.options {
		var line string
		var currentValue string

		// Get current value based on setting key
		switch option.Key {
		case "theme":
			currentValue = settings.Theme
		case "sound_enabled":
			currentValue = fmt.Sprintf("%t", settings.SoundEnabled)
		case "animations_enabled":
			currentValue = fmt.Sprintf("%t", settings.AnimationsEnabled)
		case "auto_save":
			currentValue = fmt.Sprintf("%t", settings.AutoSave)
		case "default_buy_in":
			currentValue = fmt.Sprintf("%d", settings.DefaultBuyIn)
		case "show_probabilities":
			currentValue = fmt.Sprintf("%t", settings.ShowProbabilities)
		}

		// Format the line
		line = fmt.Sprintf("%-20s: %s", option.Label, currentValue)

		if i == v.selected {
			b.WriteString(selectedItemStyle.Render("▶ " + line))
			b.WriteString("\n")
			// Show description for selected item
			description := helpStyle.Render("  " + option.Description)
			b.WriteString(description)
		} else {
			b.WriteString(itemStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}

	// Help text
	helpText := helpStyle.Render(
		"\n" +
			"• ↑/↓ or j/k to navigate\n" +
			"• Enter/Space to toggle boolean settings\n" +
			"• ←/→ or h/l to adjust numeric settings\n" +
			"• Esc to go back, q to quit",
	)
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
func (v *SettingsView) GetType() ViewType {
	return ViewSettings
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
