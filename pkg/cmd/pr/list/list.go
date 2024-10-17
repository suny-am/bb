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
)

type PrListOptions struct {
	credentials  string
	workspace    string
	repository   string
	titleFilter  string
	authorFilter string
	limit        int
}

var opts PrListOptions

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pullrequests",
	Long:  `List one or more public or workspace related pullrequests`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.limit < 0 {
			return errors.New("limit cannot be negative or 0")
		}

		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		//pullrequests, err := listPullrequests(&opts)
		//if err != nil {
		//	return err
		//}

		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ListCmd.Flags().StringVarP(&opts.repository, "repo", "r", "", "Target repository")
	ListCmd.Flags().StringVarP(&opts.titleFilter, "title", "t", "", "Title match filter")
	ListCmd.Flags().StringVarP(&opts.authorFilter, "author", "a", "", "Author name match filter")
	ListCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Item limit")

	ListCmd.MarkFlagRequired("workspace")
	ListCmd.MarkFlagRequired("repo")
}
