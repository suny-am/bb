package configure

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/list"
	"github.com/suny-am/bb/internal/spinner"
)

var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure settings",
	Long:  `Configure settings of the CLI, like output schema, theme and more`,

	RunE: func(cmd *cobra.Command, args []string) error {
		var items []list.Item

		itemKeys, err := config.GetConfiguredItems()
		if err != nil {
			return err
		}

		sort.Strings(itemKeys)

		for _, k := range itemKeys {
			items = append(items, list.Item(k))
		}

		updateFunc := func(m list.ListModel, keypress string) (tea.Model, tea.Cmd) {
			switch keypress {
			case "q", "ctrl+c":
				m.Quitting = true
				return m, tea.Quit

			case "enter":
				i, ok := m.List.SelectedItem().(list.Item)
				if ok {
					m.Choice = string(i)
					switch m.Choice {
					case "spinner":
						spinner.Configure()
					default:
						configureOption(m.Choice)
					}
				}
				return m, tea.Quit
			}
			return m, nil
		}

		list.DrawList(items, updateFunc)

		return nil
	},
}
