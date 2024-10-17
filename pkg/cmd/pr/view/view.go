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
package view

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
	"github.com/suny-am/bitbucket-cli/pkg/cmd/repo/view/forks"
)

type ViewOptions struct {
	repository  string
	workspace   string
	pullrequest string
	credentials string
}

var opts ViewOptions

var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a pullrequest",
	Long:  `View a pullrequest in a given workspace`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("<pullrequest> argument is required")
		}

		if len(args) > 1 {
			return errors.New("only one <pullrequest> argument is allowed")
		}

		opts.pullrequest = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
		/*
			repo, err := viewPullrequest(&opts)

			if err != nil {
				return err
			}
		*/
		return nil
	},
}

func init() {
	ViewCmd.AddCommand(forks.ForksCmd)

	ViewCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ViewCmd.Flags().StringVarP(&opts.repository, "repo", "r", "", "Target repository")
	ViewCmd.MarkFlagRequired("workspace")
	ViewCmd.MarkFlagRequired("repo")
}
