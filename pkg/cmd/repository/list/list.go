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
	"strings"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/internal/table"
)

type ListOptions struct {
	credentials string
	workspace   string
	nameFilter  string
	limit       int
	sort        string
}

var opts ListOptions

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	Long:  `List one or more personal and/or workspace repositories`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.limit < 0 {
			return errors.New("limit cannot be negative or 0")
		}

		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
		repos, err := getRepos(&opts, cmd)
		if err != nil {
			return err
		}

		if len(repos.Values) == 0 {
			fmt.Println(api.NoResults)
			return nil
		}

		if err := drawRepoTable(repos); err != nil {
			return err
		}

		return nil
	},
}

func drawRepoTable(repos *api.Repositories) error {
	headerData := []table.ColumnData{
		{Key: "Name"},
		{Key: "Description"},
		{Key: "Access"},
		{Key: "Updated"},
	}
	rowData := []table.RowData{}

	for _, r := range repos.Values {
		var access string
		if r.Is_Private {
			access = "Private"
		} else {
			access = "Public"
		}

		desc := strings.ReplaceAll(r.Description, "\r\n", " ")

		rowData = append(rowData, table.RowData{
			Content: []string{
				r.Name, desc, access, r.Updated_On,
			},
			Link: &r.Links.Html.Href,
		})
	}

	table.Draw(headerData, rowData)

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
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
	ListCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Item limit")
	ListCmd.Flags().StringVarP(&opts.nameFilter, "name", "n", "", "Name match filter")
	ListCmd.Flags().StringVarP(&opts.sort, "sort", "s", "", "Sorting mode")
}
