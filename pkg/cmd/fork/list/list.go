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

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/keyring"
)

type ForksOptions struct {
	current     bool
	repository  string
	workspace   string
	credentials string
}

var opts ForksOptions

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forks for a repository (TBD)",
	Long:  `List one ore more forks for a given repository (TBD)`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("<repository> argument is required")
		}

		if len(args) > 1 {
			return errors.New("only one <repository> argument is allowed")
		}

		opts.repository = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		forks, err := getForks(&opts, cmd)
		if err != nil {
			return err
		}

		if len(forks.Values) == 0 {
			fmt.Println(api.NoResults)
			return nil
		}

		fmt.Println("COMMAND IS WORK IN PROGRESS")
		fmt.Println(forks)
		return nil
	},
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

	ListCmd.Flags().BoolVarP(&opts.current, "current", "cd", true, "Reference repository from current directory")
	ListCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
}
