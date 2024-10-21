package table

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

type tableModel struct {
	table      *table.Table
	headerData []HeaderModel
	rowData    []RowModel

	coordinates  coordinates
	focusedRow   int
	windowWidth  int
	windowHeight int

	leftOffset int
	topOffset  int

	debug bool
}

type coordinates struct {
	x int
	y int
}

const (
	topOffset = 4 // offset required to register correct Y coordinate for row event ( 1 border + 1 headeing 1 border + 1 position)
)

func (tm tableModel) Init() tea.Cmd {
	return nil
}

func (tm tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		tm.windowWidth = msg.Width
		tm.windowHeight = msg.Height

	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return tm, tea.Quit
		case "down":
			if tm.focusedRow < len(tm.rowData)-1 {
				tm.focusedRow++
			}
		case "up":
			if tm.focusedRow > 0 {
				tm.focusedRow--
			}
		case "right":
			if tm.leftOffset < 300 {
				tm.leftOffset++
			}
		case "left":
			if tm.leftOffset > 0 {
				tm.leftOffset--
			}
		case "enter":
			if tm.rowData[tm.focusedRow].Link != nil {
				exec.Command("open", *tm.rowData[tm.focusedRow].Link).Start()
			}
		case "i":
			tm.debug = !tm.debug
		}

	case tea.MouseMsg:

		if tea.MouseEvent(msg).IsWheel() {
			switch tea.MouseButton(msg.Button) {
			case tea.MouseButtonWheelUp:
				if tm.focusedRow > 0 {
					tm.focusedRow--
				}
			case tea.MouseButtonWheelDown:
				if tm.focusedRow < len(tm.rowData)-1 {
					tm.focusedRow++
				}
			}
		}

		switch msg.Action {
		case tea.MouseAction(tea.MouseButtonLeft):
			tm.coordinates = coordinates{msg.X, msg.Y}
			rY := msg.Y - topOffset
			if rY >= 0 && rY < len(tm.rowData) {
				tm.focusedRow = rY
				if tm.rowData[tm.focusedRow].Link != nil && msg.Shift {
					exec.Command("open", *tm.rowData[tm.focusedRow].Link).Start()
				}

			}

		case tea.MouseAction(tea.MouseButtonWheelRight):
			tm.coordinates = coordinates{msg.X, msg.Y}
			if tm.focusedRow < len(tm.rowData)-1 {
				tm.focusedRow++
			}

		case tea.MouseAction(tea.MouseButtonWheelLeft):
			tm.coordinates = coordinates{msg.X, msg.Y}
			if tm.focusedRow > 0 {
				tm.focusedRow--
			}
		}

	}
	for i := 0; i < len(tm.rowData); i++ {
		if i == tm.focusedRow {
			tm.rowData[i].Focused = true
		} else {
			tm.rowData[i].Focused = false
		}
	}

	tm.table = genTable(tm.headerData, tm.rowData, tm.windowWidth, tm.windowHeight)
	return tm, nil
}

func (tm tableModel) View() string {
	var s string
	if tm.debug {
		s = fmt.Sprintf("width: %d, height: %d leftOffset: %d, topOffset: %d, focusedRow: %d, X: %d, Y: %d",

			tm.windowWidth,
			tm.windowHeight,
			tm.leftOffset,
			tm.topOffset,
			tm.focusedRow,
			tm.coordinates.x,
			tm.coordinates.y)
	}
	s = fmt.Sprintf("%s\n%s", s, tm.table.String())

	return s
}

func Draw(headerData []HeaderModel, rowData []RowModel) {
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))

	tm := tableModel{
		genTable(headerData, rowData, width, height),
		headerData,
		rowData,
		coordinates{0, 0},
		0,
		width,
		height,
		0,
		0,
		true,
	}

	p := tea.NewProgram(tm,
		tea.WithMouseAllMotion(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func genTable(headerData []HeaderModel, rowData []RowModel, width int, height int) *table.Table {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("#1188cc")).Bold(true)
	selectedStyle := baseStyle.Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	headers := []string{}

	for _, cm := range headerData {
		headers = append(headers, cm.Key)
	}

	tableRows := [][]string{}

	for i, rm := range rowData {
		if i <= height-topOffset {
			tableRows = append(tableRows, rm.Data)
		}
	}

	t := table.New().
		Headers(headers...).
		Rows(tableRows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			} else if rowData[row-1].Focused {
				return selectedStyle
			} else {
				return baseStyle
			}
		}).
		Border(lipgloss.ThickBorder()).
		Width(width).
		Height(height)

	return t
}
