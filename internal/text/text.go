package text

import (
	"strings"
	"unicode"

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

func TruncateSimple(text string, max int) string {
	lastSpaceIdx := -1
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIdx = i
		}
		len++
		if len >= max {
			if lastSpaceIdx != -1 {
				return text[:lastSpaceIdx] + "..."
			}
		}
	}
	return text
}

func FormatCR(text string) string {
	split := strings.Split(text, "\r\n")
	return strings.Join(split, " ")
}
