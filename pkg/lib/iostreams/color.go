package iostreams

import (
	"fmt"
	"strings"

	"github.com/mgutz/ansi"
)

var (
	magenta            = ansi.ColorFunc("magenta")
	cyan               = ansi.ColorFunc("cyan")
	red                = ansi.ColorFunc("red")
	yellow             = ansi.ColorFunc("yellow")
	blue               = ansi.ColorFunc("blue")
	green              = ansi.ColorFunc("green")
	gray               = ansi.ColorFunc("black+h")
	lightGrayUnderline = ansi.ColorFunc("white+du")
	bold               = ansi.ColorFunc("default+b")
	magentaBold        = ansi.ColorFunc("magenta+b")
	cyanBold           = ansi.ColorFunc("cyan+b")
	redBold            = ansi.ColorFunc("red+b")
	yellowBold         = ansi.ColorFunc("yellow+b")
	blueBold           = ansi.ColorFunc("blue+b")
	greenBold          = ansi.ColorFunc("green+b")
	grayBold           = ansi.ColorFunc("black+hb")

	gray256 = func(t string) string {
		return fmt.Sprintf("\x1b[%d;5;%dm%s\x1b[m", 38, 242, t)
	}
)

func NewColorScheme(enabled, is256enabled bool, hasTrueColor bool) *ColorScheme {
	return &ColorScheme{
		enabled:      enabled,
		is256enabled: is256enabled,
		hasTrueColor: hasTrueColor,
	}
}

type ColorScheme struct {
	enabled      bool
	is256enabled bool
	hasTrueColor bool
}

func (c *ColorScheme) Enabled() bool {
	return c.enabled
}

func (c *ColorScheme) Bold(t string) string {
	if !c.enabled {
		return t
	}
	return bold(t)
}

func (c *ColorScheme) BoldF(t string, args ...interface{}) string {
	return c.Bold(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) Magenta(t string) string {
	if !c.enabled {
		return t
	}
	return magenta(t)
}

func (c *ColorScheme) MagentaF(t string, args ...interface{}) string {
	return c.Magenta(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) MagentaBold(t string) string {
	if !c.enabled {
		return t
	}
	return magentaBold(t)
}

func (c *ColorScheme) Cyan(t string) string {
	if !c.enabled {
		return t
	}
	return cyan(t)
}

func (c *ColorScheme) CyanF(t string, args ...interface{}) string {
	return c.Cyan(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) CyanBold(t string) string {
	if !c.enabled {
		return t
	}
	return cyanBold(t)
}

func (c *ColorScheme) Red(t string) string {
	if !c.enabled {
		return t
	}
	return red(t)
}

func (c *ColorScheme) RedF(t string, args ...interface{}) string {
	return c.Red(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) RedBold(t string) string {
	if !c.enabled {
		return t
	}
	return redBold(t)
}

func (c *ColorScheme) Yellow(t string) string {
	if !c.enabled {
		return t
	}
	return yellow(t)
}

func (c *ColorScheme) YellowF(t string, args ...interface{}) string {
	return c.Yellow(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) YellowBold(t string) string {
	if !c.enabled {
		return t
	}
	return yellowBold(t)
}

func (c *ColorScheme) Blue(t string) string {
	if !c.enabled {
		return t
	}
	return blue(t)
}

func (c *ColorScheme) BlueF(t string, args ...interface{}) string {
	return c.Magenta(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) BlueBold(t string) string {
	if !c.enabled {
		return t
	}
	return blueBold(t)
}

func (c *ColorScheme) Green(t string) string {
	if !c.enabled {
		return t
	}
	return green(t)
}

func (c *ColorScheme) GreenF(t string, args ...interface{}) string {
	return c.Green(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) GreenBold(t string) string {
	if !c.enabled {
		return t
	}
	return greenBold(t)
}

func (c *ColorScheme) Gray(t string) string {
	if !c.enabled {
		return t
	}
	return gray(t)
}

func (c *ColorScheme) GrayF(t string, args ...interface{}) string {
	return c.Gray(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) GrayBold(t string) string {
	if !c.enabled {
		return t
	}
	return grayBold(t)
}

func (c *ColorScheme) LightGrayUnderline(t string) string {
	if !c.enabled {
		return t
	}
	return lightGrayUnderline(t)
}

func (c *ColorScheme) LightGrayUnderlineF(t string, args ...interface{}) string {
	return c.LightGrayUnderline(fmt.Sprintf(t, args...))
}

func (c *ColorScheme) SuccessIcon() string {
	return c.SuccessIconWithColor(c.Green)
}

func (c *ColorScheme) SuccessIconWithColor(colo func(string) string) string {
	return colo("✓")
}

func (c *ColorScheme) WarningIcon() string {
	return c.WarningIconWithColor(c.Green)
}

func (c *ColorScheme) WarningIconWithColor(colo func(string) string) string {
	return colo("!")
}

func (c *ColorScheme) FailureIcon() string {
	return c.FailureIconWithColor(c.Green)
}

func (c *ColorScheme) FailureIconWithColor(colo func(string) string) string {
	return colo("X")
}

func (c *ColorScheme) ColorFromString(s string) func(string) string {
	s = strings.ToLower(s)
	var fn func(string) string
	switch s {
	case "red":
		fn = c.Red
	case "green":
		fn = c.Green
	case "blue":
		fn = c.Blue
	case "yellow":
		fn = c.Yellow
	case "magenta":
		fn = c.Magenta
	case "cyan":
		fn = c.Cyan
	case "gray":
		fn = c.Gray
	case "bold":
		fn = c.Bold
	default:
		fn = func(s string) string {
			return s
		}
	}

	return fn
}
