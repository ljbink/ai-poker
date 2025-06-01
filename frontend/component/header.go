package component

import (
	"github.com/charmbracelet/lipgloss"
)

// HeaderComponent represents a reusable header component
type HeaderComponent struct {
	titleStyle lipgloss.Style
	title      string
	width      int
}

// NewHeaderComponent creates a new header component with consistent styling
func NewHeaderComponent(title string, width int) *HeaderComponent {
	// Title style matching the existing design - remove padding
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A78BFA")). // Light purple
		Align(lipgloss.Center)

	return &HeaderComponent{
		titleStyle: titleStyle,
		title:      title,
		width:      width,
	}
}

// Render renders the header using the stored title and width
func (h *HeaderComponent) Render() string {
	titleRendered := h.titleStyle.Render(h.title)

	return lipgloss.NewStyle().
		Width(h.width).
		Align(lipgloss.Center).
		Render(titleRendered)
}

// SetTitle updates the header title
func (h *HeaderComponent) SetTitle(title string) {
	h.title = title
}

// SetWidth updates the header width
func (h *HeaderComponent) SetWidth(width int) {
	h.width = width
}
