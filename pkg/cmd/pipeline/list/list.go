package list

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/internal/spinner"
	"github.com/suny-am/bb/internal/table"
)

type ListOptions struct {
	credentials string
	workspace   string
	repository  string
	limit       int
}

var opts ListOptions

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pipelines",
	Long:  `List one or more personal and/or workspace repository pipelines`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.limit < 0 {
			return errors.New("limit cannot be negative or 0")
		}

		var pipelines *api.Pipelines
		var err error

		go func() {
			opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
			pipelines, err = getPipelines(&opts)
			spinner.Stop()
		}()

		spinner.Start("Fetching pipelines")

		if err != nil {
			return err
		}

		if len(pipelines.Values) == 0 {
			fmt.Println("No results")
			return nil
		}

		if err := drawPipelineTable(pipelines); err != nil {
			return err
		}

		return nil
	},
}

func drawPipelineTable(pipelines *api.Pipelines) error {
	headerData := []table.HeaderModel{
		{Key: "Repository"},
		{Key: "Creator"},
		{Key: "Created"},
		{Key: "Completed"},
		{Key: "Error"},
		{Key: "State"},
	}
	rowData := []table.RowModel{}

	for i, p := range pipelines.Values {
		var state string

		switch p.State.Result.Name {

		case "FAILED":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(p.State.Result.Name)
		case "SUCCESSFUL":
			state = lipgloss.NewStyle().Foreground(lipgloss.Color("#119911")).Render(p.State.Result.Name)
		default:
			state = p.State.Result.Name
		}

		rowData = append(rowData, table.RowModel{
			Id: fmt.Sprintf("%d", i),
			Data: []string{
				p.Repository.Name,
				p.Creator.Display_Name,
				p.Created_On,
				p.Completed_On,
				p.State.Result.Error.Message,
				state,
			},
			Link: &p.Links.Html.Href,
		})
	}

	table.Draw(headerData, rowData)

	return nil
}

func init() {
	var workspaceDefaultValue string
	defaultWorkspace, err := config.GetWorkspace()
	if err != nil {
		ListCmd.MarkFlagRequired("workspace")
		workspaceDefaultValue = ""
	} else {
		workspaceDefaultValue = defaultWorkspace
	}
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
	ListCmd.Flags().StringVarP(&opts.repository, "repository", "r", "", "Target repository")
	ListCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Item limit")
	ListCmd.MarkFlagRequired("repository")
}
