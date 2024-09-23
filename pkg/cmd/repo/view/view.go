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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
)

type ViewOptions struct {
	repository  string
	workspace   string
	credentials string
}

var opts ViewOptions

var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a repository",
	Long:  `View a repository in a given workspace`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("<repository> argument is required")
		}

		if len(args) > 1 {
			return errors.New("only one <repository> argument is allowed")
		}

		opts.repository = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		repo, err := viewRepo(&opts)

		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\n", repo.Full_Name)
		fmt.Printf("Description: %s\n", repo.Description)
		fmt.Printf("Created: %s\n", repo.Created_On)
		fmt.Printf("Updated: %s\n", repo.Updated_On)
		fmt.Printf("Private: %v\n", repo.Is_Private)
		fmt.Printf("Fork policy: %v\n", repo.Fork_Policy)
		fmt.Printf("Wiki: %v\n", repo.Has_Wiki)
		fmt.Printf("Language: %v\n", repo.Language)
		fmt.Printf("Owner: %s [%s]\n", repo.Owner.Display_Name, repo.Owner.Type)
		fmt.Printf("Size (MB): %v\n", repo.Size/1000/1000)
		fmt.Printf("Project: %v\n", repo.Project.Name)
		fmt.Printf("Main branch: %v\n", repo.Mainbranch.Name)
		fmt.Printf("Readme: %v\n", repo.Readme)

		return nil
	},
}

func init() {
	ViewCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ViewCmd.MarkFlagRequired("workspace")
}
