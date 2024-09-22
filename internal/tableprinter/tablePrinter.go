package tablePrinter

import (
	"fmt"
	"io"

	"github.com/suny-am/bitbucket-cli/internal/text"
)

type fieldOption func(*tableField)

type TablePrinter interface {
	Header([]string, ...fieldOption)
	Field(string, ...fieldOption)
	EndRow()
	Render() error
}

func WithTruncate(fn func(int, string) string) fieldOption {
	return func(f *tableField) {
		f.truncateFunc = fn
	}
}

func WithPadding(fn func(int, string) string) fieldOption {
	return func(f *tableField) {
		f.paddingFunc = fn
	}
}

func WithColor(fn func(string) string) fieldOption {
	return func(f *tableField) {
		f.colorFunc = fn
	}
}

func New(w io.Writer, inTerminal bool, maxWidth int) TablePrinter {
	if inTerminal {
		return &terminalTablePrinter{
			out:      w,
			maxWidth: maxWidth,
		}
	}
	return nil
}

type tableField struct {
	text         string
	truncateFunc func(int, string) string
	paddingFunc  func(int, string) string
	colorFunc    func(string) string
}

type terminalTablePrinter struct {
	out        io.Writer
	maxWidth   int
	withHeader bool
	rows       [][]tableField
}

func (t *terminalTablePrinter) Header(columns []string, opts ...fieldOption) {
	if t.withHeader {
		return
	}

	t.withHeader = true
	for _, column := range columns {
		t.Field(column, opts...)
	}
	t.EndRow()
}

func (t *terminalTablePrinter) Field(s string, opts ...fieldOption) {
	if t.rows == nil {
		t.rows = make([][]tableField, 1)
	}
	rowIdx := len(t.rows) - 1
	field := tableField{
		text:         s,
		truncateFunc: text.Truncate,
	}
	for _, opt := range opts {
		opt(&field)
	}
	t.rows[rowIdx] = append(t.rows[rowIdx], field)
}

func (t *terminalTablePrinter) EndRow() {
	t.rows = append(t.rows, []tableField{})
}
func (t *terminalTablePrinter) Render() error {

	if len(t.rows) == 0 {
		return nil
	}

	delim := " "
	colCount := len(t.rows[0])
	colWidths := t.calculateColumnWidths(len(delim))

	for _, row := range t.rows {
		for col, field := range row {
			if col > 0 {
				_, err := fmt.Fprint(t.out, delim)
				if err != nil {
					return err
				}
			}
			truncVal := field.text
			if field.truncateFunc != nil {
				truncVal = field.truncateFunc(colWidths[col], field.text)
			}
			if field.paddingFunc != nil {
				truncVal = field.paddingFunc(colWidths[col], truncVal)
			} else if col < colCount-1 {
				truncVal = text.PadRight(colWidths[col], truncVal)
			}
			if field.colorFunc != nil {
				truncVal = field.colorFunc(truncVal)
			}
			_, err := fmt.Fprint(t.out, truncVal)
			if err != nil {
				return err
			}
		}
		if len(row) > 0 {
			_, err := fmt.Fprint(t.out, "\n")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *terminalTablePrinter) calculateColumnWidths(delimSize int) []int {
	colCount := len(t.rows[0])
	maxColWidths := make([]int, colCount)
	colWidths := make([]int, colCount)

	for _, row := range t.rows {
		for col, field := range row {
			w := text.DisplayWidth(field.text)
			if w > maxColWidths[col] {
				maxColWidths[col] = w
			}
			if field.truncateFunc == nil && w > colWidths[col] {
				colWidths[col] = w
			}
		}
	}

	availWidth := func() int {
		setWidths := 0
		for col := 0; col < colCount; col++ {
			setWidths += colWidths[col]
		}
		return t.maxWidth - delimSize*(colCount-1) - setWidths
	}
	numFixedCols := func() int {
		fixedCols := 0
		for col := 0; col < colCount; col++ {
			if colWidths[col] > 0 {
				fixedCols++
			}
		}
		return fixedCols
	}

	if w := availWidth(); w > 0 {
		if numFlexCols := colCount - numFixedCols(); numFlexCols > 0 {
			perCol := w / numFlexCols
			for col := 0; col < colCount; col++ {
				if max := maxColWidths[col]; max < perCol {
					colWidths[col] = max
				}
			}
		}
	}

	if w := availWidth(); w > 0 {
		for col := 0; col < colCount; col++ {
			d := maxColWidths[col] - colWidths[col]
			toAdd := w
			if d < toAdd {
				toAdd = d
			}
			colWidths[col] += toAdd
			w -= toAdd
			if w <= 0 {
				break
			}
		}
	}

	return colWidths
}
