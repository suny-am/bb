/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package list

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/internal/style"
	"github.com/suny-am/bb/internal/table"
	"github.com/suny-am/bb/internal/util"
)

type ListOptions struct {
	credentials string
	workspace   string
	repository  string
	limit       int
	current     bool
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

		if opts.current {
			opts.repository = util.GetCurrentDir()
		}

		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
		pipelines, err := getPipelines(&opts, cmd)
		if err != nil {
			return err
		}

		if len(pipelines.Values) == 0 {
			fmt.Println(api.NoResults)
			return nil
		}

		if err := drawPipelineTable(pipelines); err != nil {
			return err
		}

		return nil
	},
}

func drawPipelineTable(pipelines *api.Pipelines) error {
	headerData := []table.ColumnData{
		{Key: "Repository"},
		{Key: "Creator"},
		{Key: "Created"},
		{Key: "Completed"},
		{Key: "State"},
	}

	rowData := []table.RowData{}

	for _, p := range pipelines.Values {

		link := p.Repository.Links.Html.Href + "/pipelines/results/" + strconv.Itoa(p.Build_Number)

		rowData = append(rowData, table.RowData{
			Content: []string{
				p.Repository.Name,
				p.Creator.Display_Name,
				p.Created_On,
				p.Completed_On,
				setState(p.State),
			},
			Link: &link,
		})
	}

	table.Draw(headerData, rowData)

	return nil
}

func setState(s api.PipelineState) string {
	var state string

	// TODO: read icons from config
	switch true {
	case s.Name == "PENDING":
		state = "🕗"
	case s.Result.Name == "FAILED":
		state = "❌"
	case s.Result.Name == "SUCCESSFUL":
		state = "✅"
	case s.Result.Name == "STOPPED":
		state = "⛔️"
	case s.Stage.Name == "PAUSED":
		state = "😴"
	case s.Stage.Name == "RUNNING":
		state = "⚡️"
	default:
		state = "👽"
	}

	return style.CenterAlignStyle.Render(state)
}

func init() {
	var workspaceDefaultValue string
	defaultWorkspace, err := config.GetWorkspace()
	if err != nil {
		_ = ListCmd.MarkFlagRequired("workspace")
		workspaceDefaultValue = ""
	} else {
		workspaceDefaultValue = defaultWorkspace
	}
	ListCmd.Flags().BoolVarP(&opts.current, "current", "c", false, "Reference repository from current directory")
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
	ListCmd.Flags().StringVarP(&opts.repository, "repository", "r", "", "Target repository")
	ListCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Item limit")
	ListCmd.MarkFlagsMutuallyExclusive("current", "repository")
	ListCmd.MarkFlagsOneRequired("current", "repository")
}
