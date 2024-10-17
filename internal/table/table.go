package table

import (
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/suny-am/bitbucket-cli/api"
)

func Draw(data api.Repositories, headers []string) {
	app := tview.NewApplication()
	app.EnableMouse(true)
	table := tview.NewTable().
		SetBorders(false)
	table.SetBackgroundColor(tcell.ColorDefault)
	// headers
	for i, v := range headers {
		table.SetCell(0, i,
			tview.NewTableCell(v).
				SetSelectable(false).
				SetBackgroundColor(tcell.ColorDarkBlue).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter),
		)
	}
	// data
	for r, v := range data.Values {
		for c := range headers {
			color := tcell.ColorWhite
			var text string
			var maxWidth *int
			switch headers[c] {
			case "NAME":
				text = v.Name
			case "DESCRIPTION":
				if v.Description != "" {
					text = v.Description
					maxWidthVal := 80
					maxWidth = &maxWidthVal
				} else {
					text = "NA"
					color = tcell.ColorRed
				}
			case "VISIBILITY":
				if v.Is_Private {
					text = "private"
					color = tcell.ColorRed
				} else {
					text = "public"
				}
			case "UPDATED":
				text = v.Updated_On
			}
			cell := tview.NewTableCell(text).
				SetTextColor(color).
				SetAlign(tview.AlignLeft).
				SetReference(v.Links.Html.Href).
				SetClickedFunc(func() bool {
					cellRef := table.GetCell(r+1, c).GetReference().(string)
					err := exec.Command("open", cellRef).Start()
					return err == nil
				})
			if maxWidth != nil {
				cell.SetMaxWidth(*maxWidth)
			}
			table.SetCell(r+1, c, cell)
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
	table.SetSelectable(true, false)
	table.SetSelectedStyle(
		tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorWhiteSmoke))
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		table.SetSelectedFunc(func(row, column int) {
			cellRef := table.GetCell(row, column).GetReference().(string)
			exec.Command("open", cellRef).Start()
		})
	})
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
