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

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/internal/table2"
	"github.com/suny-am/bb/internal/util"
)

type PrListOptions struct {
	current       bool
	credentials   string
	workspace     string
	repository    string
	titleFilter   string
	creatorFilter string
	stateFilter   string
	approvals     int
	limit         int
}

var (
	approvalCountStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#22dd99"))
	commentCountStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcc00"))
	zeroCountStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	opts PrListOptions
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pullrequests",
	Long:  `List one or more public or workspace related pullrequests`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.limit < 0 {
			return errors.New("limit cannot be negative or 0")
		}

		if opts.current {
			opts.repository = util.GetCurrentDir()
		}

		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
		pullrequests, err := getPullrequests(&opts, cmd)
		if err != nil {
			return err
		}

		if len(pullrequests.Values) == 0 {
			fmt.Println(api.NoResults)
			return nil
		}

		if err := viewPullrequests(pullrequests); err != nil {
			return err
		}

		return nil
	},
}

func viewPullrequests(pullrequests *api.Pullrequests) error {
	headerData := []table2.ColumnData{
		{Key: "Branch"},
		{Key: "Repository"},
		{Key: "Creator"},
		{Key: "Comments"},
		{Key: "Approvals"},
		{Key: "State"},
		{Key: "Updated"},
	}
	rowData := []table2.RowData{}

	for _, p := range pullrequests.Values {

		approvalCount := 0

		for _, pcp := range p.Participants {
			if pcp.Approved {
				approvalCount++
			}
		}

		if opts.approvals >= 0 && approvalCount > opts.approvals {
			continue
		}

		var approvalCountText string
		if approvalCount > 0 {
			approvalCountText = approvalCountStyle.Render(fmt.Sprintf("%d", approvalCount))
		} else {
			approvalCountText = zeroCountStyle.Render(fmt.Sprintf("%d", approvalCount))
		}

		var commentCountText string
		if p.Comment_Count > 0 {
			commentCountText = commentCountStyle.Render(fmt.Sprintf("%d", p.Comment_Count))
		} else {
			commentCountText = zeroCountStyle.Render(fmt.Sprintf("%d", p.Comment_Count))
		}

		rowData = append(rowData, table2.RowData{
			Content: []string{
				p.Source.Branch.Name,
				p.Source.Repository.Name,
				p.Author.Nickname,
				commentCountText,
				approvalCountText,
				p.State,
				p.Updated_On,
			},
			Link: &p.Links.Html.Href,
		})
	}

	table2.Draw(headerData, rowData)
	return nil
}

func init() {
	var workspaceDefaultValue string
	defaultWorkspace, err := config.GetWorkspace()
	if err != nil {
		workspaceDefaultValue = ""
	} else {
		workspaceDefaultValue = defaultWorkspace
	}

	ListCmd.Flags().BoolVarP(&opts.current, "current", "c", false, "Reference repository from current directory")
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
	ListCmd.Flags().StringVarP(&opts.repository, "repo", "r", "", "Target repository")
	ListCmd.Flags().StringVarP(&opts.titleFilter, "title", "t", "", "Title match filter")
	ListCmd.Flags().StringVar(&opts.creatorFilter, "creator", "", "Creator match filter")
	ListCmd.Flags().StringVarP(&opts.stateFilter, "state", "s", "", "Pullrequest state filter")
	ListCmd.Flags().IntVarP(&opts.approvals, "approvals", "a", -1, "Approvals count filter")
	ListCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Item limit")
}
