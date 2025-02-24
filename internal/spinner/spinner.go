package spinner

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/suny-am/bb/internal/config"
)

var (
	// Available spinners
	spinners = []spinner.Spinner{
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}

	textStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type spinnerModel struct {
	message       string
	spinner       spinner.Model
	index         int
	quitting      bool
	configureMode bool
}

type (
	stopMessage       struct{}
	viewUpdateMessage struct {
		content string
	}
)

const helpLegend = "←/↑: prev   |   →/↓: next"

var (
	spinnerInstance spinnerModel
	spinnerProgram  *tea.Program
	helpStyle       = lipgloss.NewStyle().PaddingLeft(4).PaddingBottom(1).Faint(true)
)

func Start(message string) {
	spinnerStyle := 1
	configSpinner, err := config.GetSpinnerStyle()
	if err == nil {
		spinnerStyle = *configSpinner
	}
	spinnerInstance = spinnerModel{message: message}

	spinnerInstance.initSpinner(spinnerStyle)

	spinnerProgram = tea.NewProgram(spinnerInstance)
	if _, err := spinnerProgram.Run(); err != nil {
		fmt.Println("Could not run Spinner program")
	}
}

func Configure() {
	spinnerStyle := 1
	configSpinner, err := config.GetSpinnerStyle()
	if err == nil {
		spinnerStyle = *configSpinner
	}
	spinnerInstance = spinnerModel{configureMode: true}

	spinnerInstance.initSpinner(spinnerStyle)

	spinnerInstance.message = fmt.Sprintf("Spinner style: %d", spinnerInstance.index)

	spinnerProgram = tea.NewProgram(spinnerInstance, tea.WithAltScreen())
	if _, err := spinnerProgram.Run(); err != nil {
		fmt.Println("Could not run Spinner program")
	}
}

func Stop() {
	if spinnerProgram != nil {
		spinnerProgram.Send(stopMessage{})
	} else {
		Stop()
	}
}

func AddToView(str string) {
	if spinnerProgram != nil {
		spinnerProgram.Send(viewUpdateMessage{str})
	} else {
		AddToView(str)
	}
}

func (m *spinnerModel) initSpinner(style int) {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[style]
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "right", "down":
			if m.index < 8 && m.configureMode {
				m.index++
				m.spinner.Spinner = spinners[m.index]
				m.message = fmt.Sprintf("Spinner style: %d", m.index)
			}

		case "left", "up":
			if m.index > 0 && m.configureMode {
				m.index--
				m.spinner.Spinner = spinners[m.index]
				m.message = fmt.Sprintf("Spinner style: %d", m.index)
			}
		case "enter":
			m.quitting = true
			if m.configureMode {
				stringVal := strconv.Itoa(m.index)
				_, err := config.SetConfigOption(config.ConfigOption{Key: "spinner", Value: stringVal})
				if err == nil {
					return m, tea.Quit
				}
			}
		}
		return m, nil

	case viewUpdateMessage:
		m.message = m.message + msg.content
		return m, nil

	case stopMessage:
		m.quitting = true
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m spinnerModel) View() (s string) {
	if m.quitting {
		return s
	}

	var gap string

	switch m.index {
	case 1:
		gap = ""
	default:
		gap = " "
	}

	s += fmt.Sprintf("\n %s%s%s\n\n",
		m.spinner.View(),
		gap,
		textStyle(fmt.Sprintf("%s...", m.message)),
	)

	if m.configureMode {
		s += fmt.Sprintf("\n%s", helpStyle.Render(helpLegend))
	}

	return s
}
