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
package browse

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/git"
)

var BrowseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Bitbucket browse command",
	Long:  "Open a given repository in the default browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		var repoName string

		if len(args) == 0 {
			repoName = git.GetGitRepo()
		} else {
			repoName = args[0]
		}

		defaultWorkspace, err := config.GetWorkspace()
		if err != nil {
			return err
		}

		repoUrl := fmt.Sprintf("https://bitbucket.org/%s/%s", defaultWorkspace, repoName)

		browseCmd := exec.Command("open", repoUrl)

		browseCmd.Run()

		return nil
	},
}

func init() {
	BrowseCmd.Execute()
}
