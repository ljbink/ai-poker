package frontend

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ljbink/ai-poker/frontend/component"
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

// IndexKeyMap defines keybindings for the index view
type IndexKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k IndexKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k IndexKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Select, k.Quit},
	}
}

var indexKeys = IndexKeyMap{
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
		key.WithHelp("enter/space", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// IndexView represents the main menu/welcome screen
type IndexView struct {
	model *Model
	list  list.Model
	keys  IndexKeyMap
	help  help.Model

	// Components
	header *component.HeaderComponent
	helper *component.HelperComponent
}

// NewIndexView creates a new index view
func NewIndexView(model *Model) *IndexView {
	items := []list.Item{
		MenuItem{
			title:       "🎮 Start Game",
			description: "Begin a new poker game",
			action:      ViewLogin,
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

	// Create help component with matching SettingsView styling
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#7C3AED"))  // Purple
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")) // Medium gray

	return &IndexView{
		model: model,
		list:  l,
		keys:  indexKeys,
		help:  h,

		// Initialize components with default width (will be updated in Render)
		header: component.NewHeaderComponent("🃏 Texas Hold'em Poker", 80),
		helper: component.NewHelperComponent(indexKeys, 80),
	}
}

// Update handles input for the index view
func (v *IndexView) Update(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, v.keys.Select):
		selectedItem, ok := v.list.SelectedItem().(MenuItem)
		if ok {
			switch selectedItem.action {
			case ViewLogin:
				v.model.currentView = ViewLogin
			case ViewSettings:
				v.model.currentView = ViewSettings
			default: // Quit case
				return v.model, tea.Quit
			}
		}
		return v.model, nil
	case key.Matches(msg, v.keys.Quit):
		return v.model, tea.Quit
	}

	var cmd tea.Cmd
	v.list, cmd = v.list.Update(msg)
	return v.model, cmd
}

// Render renders the index view
func (v *IndexView) Render(width, height int) string {
	// Update component widths for current screen size
	v.header.SetWidth(width)
	v.helper.SetWidth(width)

	// Title at the top using header component
	titleAtTop := v.header.Render()

	// Help view at the bottom using helper component
	helpAtBottom := v.helper.Render()

	// Calculate actual space used by header and helper
	headerHeight := lipgloss.Height(titleAtTop)
	helperHeight := lipgloss.Height(helpAtBottom)
	availableHeight := height - headerHeight - helperHeight

	// Update list dimensions to use remaining space
	v.list.SetWidth(width - 8)
	v.list.SetHeight(availableHeight - 2) // Small buffer for list margins

	// Render list content for center area
	listView := v.list.View()

	// Center the list content in the middle of available space
	centeredContent := lipgloss.Place(
		width, availableHeight,
		lipgloss.Center, lipgloss.Center,
		listView,
	)

	// Combine title, content, and help without extra spacing
	fullContent := titleAtTop + centeredContent + helpAtBottom

	// Apply full screen style
	fullScreenContainer := GetFullScreenStyle(width, height)
	return fullScreenContainer.Render(fullContent)
}

// GetType returns the view type
func (v *IndexView) GetType() ViewType {
	return ViewIndex
}

// ShortHelp returns keybindings to be shown in the mini help view
func (v *IndexView) ShortHelp() []key.Binding {
	return v.keys.ShortHelp()
}

// FullHelp returns keybindings for the expanded help view
func (v *IndexView) FullHelp() [][]key.Binding {
	return v.keys.FullHelp()
}
