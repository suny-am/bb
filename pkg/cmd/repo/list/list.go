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

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
	"github.com/suny-am/bitbucket-cli/internal/table"
)

type ListOptions struct {
	credentials string
	workspace   string
	nameFilter  string
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

		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		repos, err := listRepos(&opts)
		if err != nil {
			return err
		}

		headers := []string{"NAME", "DESCRIPTION", "VISIBILITY", "UPDATED"}
		table.Draw(*repos, headers)

		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ListCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Item limit")
	ListCmd.Flags().StringVarP(&opts.nameFilter, "name", "n", "", "Name match filter")
}
