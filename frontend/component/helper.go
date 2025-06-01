package component

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// HelperComponent represents a reusable helper component
type HelperComponent struct {
	help   help.Model
	keyMap KeyMapInterface
	width  int
}

// KeyMapInterface defines the interface that key maps must implement
type KeyMapInterface interface {
	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
}

// NewHelperComponent creates a new helper component with consistent styling
func NewHelperComponent(keyMap KeyMapInterface, width int) *HelperComponent {
	// Create help component with matching styling
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))  // Purple
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")) // Medium gray

	return &HelperComponent{
		help:   h,
		keyMap: keyMap,
		width:  width,
	}
}

// Render renders the helper using the stored keyMap and width
func (h *HelperComponent) Render() string {
	helpView := h.help.View(h.keyMap)

	return lipgloss.NewStyle().
		Width(h.width).
		Align(lipgloss.Center).
		Render(helpView)
}

// SetKeyMap updates the helper keyMap
func (h *HelperComponent) SetKeyMap(keyMap KeyMapInterface) {
	h.keyMap = keyMap
}

// SetWidth updates the helper width
func (h *HelperComponent) SetWidth(width int) {
	h.width = width
}
