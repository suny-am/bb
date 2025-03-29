package table2

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/suny-am/bb/internal/style"
	"golang.org/x/term"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var debugMode bool

type ColumnData struct {
	Key string
}

type RowData struct {
	Content []string
	Link    *string
}

type tableModel struct {
	table          table.Model
	rowData        []RowData
	firstRowInView int
	lastRowInView  int
	debug          string
	target         string
}

func (tm tableModel) Init() tea.Cmd { return nil }

func (tm tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			debugMode = !debugMode
			return tm, nil
		case "esc":
			if tm.table.Focused() {
				tm.table.Blur()
			} else {
				tm.table.Focus()
			}
		case "q", "ctrl+c":
			return tm, tea.Quit
		case "enter":
			openLink(*tm.rowData[tm.table.Cursor()].Link)
			tm.target = fmt.Sprintf("Let's go to %s!", tm.table.SelectedRow()[1])
		case "down":
			tm.MoveDownRelative()
			return tm, nil
		case "up":
			tm.MoveUpRelative()
			return tm, nil
		}
	case tea.MouseMsg:
		if tea.MouseEvent(msg).IsWheel() {
			switch tea.MouseButton(msg.Button) {
			case tea.MouseButtonWheelDown:
				tm.MoveUpRelative()
			case tea.MouseButtonWheelUp:
				tm.MoveDownRelative()
			}
		}

		switch tea.MouseAction(msg.Button) {
		case tea.MouseAction(tea.MouseButtonLeft):
			tm.CalculateClickPos(msg)
		}
	}
	tm.table, cmd = tm.table.Update(msg)
	return tm, cmd
}

func (tm *tableModel) CalculateClickPos(msg tea.MouseMsg) {
	if msg.Y >= 0 && msg.Y-3 < tm.table.Height() { // check if click is within tables visible rows area
		clicked := msg.Y - 3
		pos := tm.firstRowInView - 1 + clicked
		for i := tm.table.Height(); i > 0; i-- {
			if pos > tm.table.Cursor() {
				tm.table.MoveDown(1)
			}
			if pos < tm.table.Cursor() {
				tm.table.MoveUp(1)
			}
			if pos > tm.table.Cursor()+i {
				tm.table.MoveDown(i + 1)
			} else if pos < tm.table.Cursor()-i {
				tm.table.MoveUp(i + 1)
			}
		}
	}
}

func (tm *tableModel) MoveUpRelative() {
	tm.table.MoveUp(1)
	if tm.table.Cursor() <= tm.lastRowInView &&
		tm.table.Cursor() < tm.firstRowInView-1 {
		tm.firstRowInView--
		tm.lastRowInView--
	}
}

func (tm *tableModel) MoveDownRelative() {
	tm.table.MoveDown(1)
	if tm.table.Cursor() >= tm.lastRowInView {
		tm.firstRowInView++
		tm.lastRowInView++
	}
}

func (tm tableModel) View() string {
	if debugMode {
		tm.debug = style.BlockStyle.Render(fmt.Sprintf("%v\nHeight: %s\nRow: %s\nFirst row: %s\nLast row: %s",
			style.DebugHeaderStyle.Render("DEBUG"),
			style.PropValueStyle.Render(strconv.Itoa(tm.table.Height())),
			style.PropValueStyle.Render(strconv.Itoa(tm.table.Cursor()+1)),
			style.PropValueStyle.Render(strconv.Itoa(tm.firstRowInView)),
			style.PropValueStyle.Render(strconv.Itoa(tm.lastRowInView)),
		))
	}
	return baseStyle.Render(tm.table.View()) +
		"\n" +
		tm.debug +
		"\n" +
		style.BlockStyle.Render(tm.target)
}

func Draw(columnData []ColumnData, rowData []RowData) {
	columns := []table.Column{}
	rows := []table.Row{}

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error getting terminal size: %v", err)
	}

	width = width - 10

	height = int(math.Min(float64(len(rowData)), float64(height))) + 1

	columns = append(columns, table.Column{
		Title: "ID",
		Width: 4,
	})

	for _, cd := range columnData {
		columns = append(columns, table.Column{
			Title: cd.Key,
			Width: (width - 4) / len(columnData),
		})
	}

	for i, r := range rowData {
		row := table.Row{strconv.Itoa(i + 1)}
		row = append(row, r.Content...)
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := tableModel{
		t,
		rowData,
		1,
		49,
		"",
		"",
	}
	if _, err := tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func openLink(link string) {
	if err := exec.Command("open", link).Start(); err != nil {
		fmt.Println("Could not open link: ", link)
	}
}
