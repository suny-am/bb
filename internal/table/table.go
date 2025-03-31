package table

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func Draw(columnData []ColumnData, rowData []RowData) {
	columns := []table.Column{}
	rows := []table.Row{}

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Error getting terminal size: %v", err)
	}

	width = width - 10

	for _, cd := range columnData {
		columns = append(columns, table.Column{
			Title: cd.Key,
			Width: width / len(columnData),
		})
	}

	for _, r := range rowData {
		row := table.Row{}
		row = append(row, r.Content...)
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(width),
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

	tm := TableModel{
		t,
		rowData,
		1,
		49,
		"",
		"",
	}

	if _, err := tea.NewProgram(
		tm,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
