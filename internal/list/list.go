package list

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	listHeight   = 14
	defaultWidth = 20
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type (
	Item         string
	itemDelegate struct{}
)

func (i Item) FilterValue() string { return "" }

func DrawList(items []Item, updateFunc UpdateFunc) {
	var listItems []list.Item

	for _, i := range items {
		listItems = append(listItems, i)
	}

	l := list.New(listItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select an item to configure:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := ListModel{List: l, UpdateFunc: updateFunc}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type UpdateFunc func(m ListModel, key string) (tea.Model, tea.Cmd)

type ListModel struct {
	List       list.Model
	UpdateFunc UpdateFunc
	Choice     string
	Quitting   bool
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.Type.String() {
		case "ctrl+c", "esc", "q":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			m, cmd := m.UpdateFunc(m, msg.String())
			return m, cmd
		}
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	if m.Quitting {
		return quitTextStyle.Render("Finished? OK!")
	}

	if m.Choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s configured.", m.Choice))
	}

	return "\n" + m.List.View()
}
