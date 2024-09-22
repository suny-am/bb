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
	"os"

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/internal/iostreams"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
	tablePrinter "github.com/suny-am/bitbucket-cli/internal/tableprinter"
)

type ListOptions struct {
	repository  string
	workspace   string
	credentials string
	limit       int
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

		cmd.Root().PreRun(cmd, nil)
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		resp, err := listRepos(&opts)

		if err != nil {
			return err
		}

		tp := tablePrinter.New(os.Stdout, true, 500)

		cs := *iostreams.NewColorScheme(true, true, true)

		headers := []string{"NAME", "INFO", "UPDATED"}
		tp.Header(headers, tablePrinter.WithColor(cs.LightGrayUnderline))
		for i := range resp.Values {
			repo := resp.Values[i]
			tp.Field(repo.Full_Name, tablePrinter.WithColor(cs.Bold))
			if repo.Is_Private {
				tp.Field("private", tablePrinter.WithColor(cs.Gray))
			} else {
				tp.Field("public", tablePrinter.WithColor(cs.Yellow))
			}
			tp.Field(repo.Updated_On, tablePrinter.WithColor(cs.Gray))
			tp.EndRow()
		}

		tp.Render()

		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ListCmd.Flags().StringVarP(&opts.repository, "repo", "r", "", "Target repository")
	ListCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Item limit")
}
