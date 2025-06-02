package frontend

import (
	tea "github.com/charmbracelet/bubbletea"
)

// View represents a screen/view in the application
type View interface {
	// Update handles input messages and returns updated model and commands
	Update(msg tea.KeyMsg) (tea.Model, tea.Cmd)

	// Render returns the string representation of the view
	Render(width, height int) string

	// GetType returns the view type for navigation
	GetType() ViewType
}
