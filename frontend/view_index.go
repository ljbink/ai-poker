package frontend

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// IndexView represents the main menu/welcome screen
type IndexView struct {
	model    *Model
	selected int
	options  []string
}

// NewIndexView creates a new index view
func NewIndexView(model *Model) *IndexView {
	return &IndexView{
		model:    model,
		selected: 0,
		options:  []string{"Start Game", "Settings", "Quit"},
	}
}

// Update handles input for the index view
func (v *IndexView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		switch v.selected {
		case 0: // Start Game
			v.model.currentView = ViewInputUserProfile
		case 1: // Settings
			v.model.currentView = ViewSettings
		case 2: // Quit
			return v.model, tea.Quit
		}
	}
	return v.model, nil
}

// Render renders the index view
func (v *IndexView) Render(width, height int) string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("🃏 Texas Hold'em Poker")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Menu options
	for i, option := range v.options {
		if i == v.selected {
			b.WriteString(selectedItemStyle.Render("▶ " + option))
		} else {
			b.WriteString(itemStyle.Render("  " + option))
		}
		b.WriteString("\n")
	}

	// Help text
	helpText := helpStyle.Render("\nUse ↑/↓ or j/k to navigate, Enter to select, q to quit")
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
func (v *IndexView) GetType() ViewType {
	return ViewIndex
}
