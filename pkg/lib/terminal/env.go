package terminal

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/term"
)

type Terminal struct {
	in           *os.File
	out          *os.File
	errOut       *os.File
	isTTY        bool
	colorEnabled bool
	is256enabled bool
	hasTrueColor bool
	width        int
	widthPercent int
}

func FromEnv() Terminal {
	var stdoutIsTTY bool
	var isColorEnabled bool
	var termWidthOverride int
	var termWidthPercentage int

	spec := os.Getenv("BIT_CLI_FORCE_TTY")
	if spec != "" {
		stdoutIsTTY = true
		isColorEnabled = !IsColorDisabled()

		if w, err := strconv.Atoi(spec); err == nil {
			termWidthOverride = w
		} else if strings.HasSuffix(spec, "%") {
			if p, err := strconv.Atoi(spec[:len(spec)-1]); err == nil {
				termWidthPercentage = p
			}
		}
	} else {
		stdoutIsTTY = IsTerminal(os.Stdout)
		isColorEnabled = IsColorForced() || (!IsColorDisabled() && stdoutIsTTY)
	}

	isVirtualTerminal := false
	if stdoutIsTTY {
		if err := enableVirtualTerminalProcessing(os.Stdout); err == nil {
			isVirtualTerminal = true
		}
	}

	return Terminal{
		in:           os.Stdin,
		out:          os.Stdout,
		errOut:       os.Stderr,
		isTTY:        stdoutIsTTY,
		colorEnabled: isColorEnabled,
		is256enabled: isVirtualTerminal || is256ColorSupported(),
		hasTrueColor: isVirtualTerminal || isTrueColorSupported(),
		width:        termWidthOverride,
		widthPercent: termWidthPercentage,
	}
}

func (t Terminal) In() io.Reader {
	return t.in
}

func (t Terminal) Out() io.Writer {
	return t.out
}

func (t Terminal) ErrOut() io.Writer {
	return t.errOut
}

func (t Terminal) IsTerminalOutput() bool {
	return t.isTTY
}

func (t Terminal) IsColorEnabled() bool {
	return t.colorEnabled
}

func (t Terminal) Is256ColorSupported() bool {
	return t.is256enabled
}

func (t Terminal) IsTrueColorSupported() bool {
	return t.hasTrueColor
}

func (t Terminal) Size() (int, int, error) {
	if t.width > 0 {
		return t.width, -1, nil
	}

	ttyOut := t.out
	if ttyOut == nil || !IsTerminal(ttyOut) {
		if f, err := openTTY(); err == nil {
			defer f.Close()
			ttyOut = f
		} else {
			return -1, -1, err
		}
	}

	width, height, err := terminalSize(ttyOut)
	if err == nil && t.widthPercent > 0 {
		return int(float64(width) * float64(t.widthPercent) / 100), height, nil
	}

	return width, height, err
}

func (t Terminal) Theme() string {
	if !t.IsColorEnabled() {
		return "none"
	}
	if termenv.HasDarkBackground() {
		return "dark"
	}
	return "light"
}

// IsTerminal reports whether a file descriptor is connected to a terminal.
func IsTerminal(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

func terminalSize(f *os.File) (int, int, error) {
	return term.GetSize(int(f.Fd()))
}

func IsColorDisabled() bool {
	return os.Getenv("NO_COLOR") != "" || os.Getenv("CLICOLOR") == "0"
}

func IsColorForced() bool {
	return os.Getenv("CLICOLOR_FORCE") != "" && os.Getenv("CLICOLOR_FORCE") != "0"
}

func is256ColorSupported() bool {
	return isTrueColorSupported() ||
		strings.Contains(os.Getenv("TERM"), "256") ||
		strings.Contains(os.Getenv("COLORTERM"), "256")
}

func isTrueColorSupported() bool {
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")

	return strings.Contains(term, "24bit") ||
		strings.Contains(term, "truecolor") ||
		strings.Contains(colorterm, "24bit") ||
		strings.Contains(colorterm, "truecolor")
}
