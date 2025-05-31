package frontend

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a menu item for the list
type MenuItem struct {
	title       string
	description string
	action      ViewType
}

func (i MenuItem) FilterValue() string { return i.title }
func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }

// Custom item delegate for styling
type menuItemDelegate struct{}

func (d menuItemDelegate) Height() int                             { return 2 }
func (d menuItemDelegate) Spacing() int                            { return 1 }
func (d menuItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d menuItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(MenuItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.title)
	desc := fmt.Sprintf("%s", i.description)

	if index == m.Index() {
		// Selected item styling
		str = selectedItemStyle.Render("▶ " + str)
		desc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB")). // Light gray
			Render("  " + desc)
	} else {
		// Normal item styling
		str = itemStyle.Render("  " + str)
		desc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")). // Medium gray
			Render("  " + desc)
	}

	fmt.Fprint(w, str+"\n"+desc)
}

// IndexView represents the main menu/welcome screen
type IndexView struct {
	model *Model
	list  list.Model
}

// NewIndexView creates a new index view
func NewIndexView(model *Model) *IndexView {
	items := []list.Item{
		MenuItem{
			title:       "🎮 Start Game",
			description: "Begin a new poker game",
			action:      ViewInputUserProfile,
		},
		MenuItem{
			title:       "⚙️  Settings",
			description: "Configure game preferences",
			action:      ViewSettings,
		},
		MenuItem{
			title:       "🚪 Quit",
			description: "Exit the application",
			action:      ViewIndex, // Special case for quit
		},
	}

	l := list.New(items, menuItemDelegate{}, 0, 0)

	// Disable all list features and title to avoid any status indicators
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false) // Disable built-in title

	// Custom styling for the list
	l.Styles.NoItems = lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")) // Gray

	return &IndexView{
		model: model,
		list:  l,
	}
}

// Update handles input for the index view
func (v *IndexView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		selectedItem, ok := v.list.SelectedItem().(MenuItem)
		if ok {
			switch selectedItem.action {
			case ViewInputUserProfile:
				v.model.currentView = ViewInputUserProfile
			case ViewSettings:
				v.model.currentView = ViewSettings
			default: // Quit case
				return v.model, tea.Quit
			}
		}
		return v.model, nil
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v.model, cmd
}

// Render renders the index view
func (v *IndexView) Render(width, height int) string {
	// Update list dimensions
	v.list.SetWidth(width - 4)
	v.list.SetHeight(15) // Reduced height to accommodate manual title

	// Manual title rendering
	title := titleStyle.Render("🃏 Texas Hold'em Poker")

	// Render list
	listView := v.list.View()

	// Add custom help text
	helpText := helpStyle.Render(
		"\n" +
			"• ↑/↓ or j/k to navigate\n" +
			"• Enter/Space to select\n" +
			"• q to quit",
	)

	// Combine title, list, and help
	content := title + "\n\n" + listView + helpText

	// Center the content
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

// ShortHelp returns keybindings to be shown in the mini help view
func (v *IndexView) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter/space", "select"),
		),
		key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// FullHelp returns keybindings for the expanded help view
func (v *IndexView) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
		},
		{
			key.NewBinding(
				key.WithKeys("enter", " "),
				key.WithHelp("enter/space", "select"),
			),
			key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
	}
}
