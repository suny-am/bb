package textinput

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/suny-am/bb/internal/config"
)

type errMsg error

type inputModel struct {
	err          error
	configOption string
	prompt       string
	input        textinput.Model
	configMode   bool
}

func Listen(prompt string, placeholder string) {
	p := tea.NewProgram(initialModel(prompt, placeholder,
		false,
		nil),
		tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func ConfigListen(prompt string, option string, currentOptionValue string) {
	p := tea.NewProgram(initialModel(option,
		currentOptionValue,
		true,
		&option),
		tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel(prompt string, placeholder string, configMode bool, configOption *string) inputModel {
	input := textinput.New()
	input.CharLimit = 156
	input.Focus()
	input.Placeholder = placeholder
	input.Placeholder = placeholder
	input.Width = 20

	return inputModel{
		err:          nil,
		prompt:       prompt,
		configOption: *configOption,
		input:        input,
		configMode:   configMode,
	}
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.configMode {
				_, err := config.SetConfigOption(config.ConfigOption{Key: m.configOption, Value: m.input.Value()})
				if err == nil {
					return m, tea.Quit
				}
				return m, nil
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf("%s: \n\n%s\n\n%s",
		m.prompt,
		m.input.View(),
		"press esc to quit") + "\n"
}
