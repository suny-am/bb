package table

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Draw() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true)
	lorem := strings.Split("Lorem ipsum aasfafasfasfafq36123622y2y423dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	cols, rows := 5, 100
	word := 0
	headers := strings.Split("1 2 3 4 5", " ")
	for i, v := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(v).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter),
		)
	}
	for r := 1; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			table.SetCell(r, c,
				tview.NewTableCell(lorem[word]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(lorem)
		}
	}
	table.SetBorderPadding(1, 1, 1, 1)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
		}
		return event
	})
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyCtrlQ {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
