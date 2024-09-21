package text

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

const (
	ellipsis            = "..."
	minWidthForEllipsis = len(ellipsis) + 2
)

func DisplayWidth(s string) int {
	return lipgloss.Width(s)
}

func Truncate(maxWidth int, s string) string {
	w := DisplayWidth(s)
	if w <= maxWidth {
		return s
	}
	tail := ""
	if maxWidth >= minWidthForEllipsis {
		tail = ellipsis
	}
	r := truncate.StringWithTail(s, uint(maxWidth), tail)
	if DisplayWidth(r) < maxWidth {
		r += " "
	}
	return r
}

func PadRight(maxWidth int, s string) string {
	if padWidth := maxWidth - DisplayWidth(s); padWidth > 0 {
		s += strings.Repeat(" ", padWidth)
	}
	return s
}
