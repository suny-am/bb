package table

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/suny-am/bb/internal/style"
	"golang.org/x/term"
)

type ColumnData struct {
	Key string
}

type RowData struct {
	Content []string
	Link    *string
}

type TableModel struct {
	table          table.Model
	rowData        []RowData
	firstRowInView int
	lastRowInView  int
	debug          string
	target         string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var debugMode bool

func (tm TableModel) Init() tea.Cmd { return nil }

func (tm TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		tm.updateTableDimensions()

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

func (tm *TableModel) CalculateClickPos(msg tea.MouseMsg) {
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

func (tm *TableModel) MoveUpRelative() {
	tm.table.MoveUp(1)
	if tm.table.Cursor() <= tm.lastRowInView &&
		tm.table.Cursor() < tm.firstRowInView-1 {
		tm.firstRowInView--
		tm.lastRowInView--
	}
}

func (tm *TableModel) MoveDownRelative() {
	tm.table.MoveDown(1)
	if tm.table.Cursor() >= tm.lastRowInView {
		tm.firstRowInView++
		tm.lastRowInView++
	}
}

func (tm TableModel) View() string {
	if debugMode {
		tm.debug = style.BlockStyle.Render(fmt.Sprintf("%v Height: %s Row: %s First row: %s Last row: %s",
			style.DebugHeaderStyle.Render("[DEBUG]"),
			style.PropValueStyle.Render(strconv.Itoa(tm.table.Height())),
			style.PropValueStyle.Render(strconv.Itoa(tm.table.Cursor()+1)),
			style.PropValueStyle.Render(strconv.Itoa(tm.firstRowInView)),
			style.PropValueStyle.Render(strconv.Itoa(tm.lastRowInView)),
		))
	}
	return baseStyle.
		Render(tm.table.View()) +
		"\n" +
		tm.debug
}

func (tm *TableModel) updateTableDimensions() {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error getting terminal size: %v", err)
	}

	updatedColumns := []table.Column{}

	for i, c := range tm.table.Columns() {
		if i == 0 {
			updatedColumns = append(updatedColumns, table.Column{
				Title: c.Title,
				Width: 5,
			})
		} else {
			updatedColumns = append(updatedColumns, table.Column{
				Title: c.Title,
				Width: (width-4)/(len(tm.table.Columns())-1) - 3,
			})
		}
	}

	tm.table.SetColumns(updatedColumns)
	tm.table.SetHeight(height - 3)
	tm.table.SetWidth(width - 4)
}

func openLink(link string) {
	if err := exec.Command("open", link).Start(); err != nil {
		fmt.Println("Could not open link: ", link)
	}
}
