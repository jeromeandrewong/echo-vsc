package themepicker

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Label string
	Path  string
}

type Model struct {
	list   list.Model
	Choice Theme
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (t Theme) Title() string       { return t.Label }
func (t Theme) Description() string { return t.Path }
func (t Theme) FilterValue() string { return t.Label }

func New(themes []Theme) Model {
	items := make([]list.Item, len(themes))
	for i, theme := range themes {
		items[i] = theme
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "VSCode Themes"

	return Model{list: l}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			t, ok := m.list.SelectedItem().(Theme)
			if ok {
				m.Choice = t
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.list.View())
}
