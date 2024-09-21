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
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/pkg/lib/iostreams"
	tablePrinter "github.com/suny-am/bitbucket-cli/pkg/lib/tableprinter"
	"github.com/suny-am/bitbucket-cli/pkg/types"
)

type (
	RepoListResponse struct {
		Values []Repository
	}
	Repository struct {
		Created_On  string
		Updated_On  string
		Description string
		Full_Name   string
		Is_Private  bool
	}
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	Long:  `List one or more personal and/or workspace repositories`,

	Run: func(cmd *cobra.Command, args []string) {

		workspace, _ := cmd.Flags().GetString("workspace")
		repository, _ := cmd.Flags().GetString("repository")

		limit, _ := cmd.Flags().GetString("limit")

		if limit == "" {
			limit = "10"
		}

		cmd.Root().PreRun(cmd, nil)
		credentials := cmd.Context().Value(types.CredentialsKey{})

		authHeaderData := fmt.Sprintf("Basic %s", credentials)

		client := resty.New()

		// TBD add workspace as argument

		endpoint := "https://api.bitbucket.org/2.0/repositories"

		if workspace != "" {
			endpoint = fmt.Sprintf("%s/%s", endpoint, workspace)

			// --repository requires --workspace
			if repository != "" {
				endpoint = fmt.Sprintf("%s/%s", endpoint, repository)
			}
		}

		resp, err := client.R().
			SetHeader("Authorization", authHeaderData).
			SetHeader("Accept", "application/json").
			SetQueryParam("pagelen", limit).
			EnableTrace().
			Get(endpoint)

		if resp.IsError() {
			fmt.Println(err.Error())
		}

		if resp.IsSuccess() {
			var response RepoListResponse

			if err := json.Unmarshal([]byte(resp.String()), &response); err != nil {
				fmt.Println(err)
			}

			tp := tablePrinter.New(os.Stdout, true, 200) // TBD pass dynamic opts

			cs := *iostreams.NewColorScheme(true, true, true)

			headers := []string{"NAME", "INFO", "UPDATED"}
			tp.Header(headers, tablePrinter.WithColor(cs.LightGrayUnderline))
			for i := range response.Values {
				repo := response.Values[i]
				tp.Field(repo.Full_Name, tablePrinter.WithColor(cs.Bold))
				tp.Field("public", tablePrinter.WithColor(cs.Gray))
				tp.Field(repo.Updated_On, tablePrinter.WithColor(cs.Gray))
				tp.EndRow()
			}

			tp.Render()
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	ListCmd.Flags().StringP("workspace", "w", "", "Target workspace")
	ListCmd.Flags().StringP("repo", "r", "", "Target repository")
	ListCmd.Flags().StringP("limit", "l", "", "Item limit")
}
