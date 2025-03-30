package style

import "github.com/charmbracelet/lipgloss"

var (
	greenFg  = lipgloss.AdaptiveColor{Light: "#00aa00", Dark: "#00ff00"}
	purpleFg = lipgloss.AdaptiveColor{Light: "#9900ff", Dark: "#ff99ff"}
)

var DebugHeaderStyle = lipgloss.NewStyle().Foreground(greenFg)

var BlockStyle = lipgloss.NewStyle().PaddingLeft(1)

var PropValueStyle = lipgloss.NewStyle().Foreground(purpleFg)

var CenterAlignStyle = lipgloss.NewStyle().Align(lipgloss.Center).PaddingLeft(2)
