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
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/pkg/cmd/commit"
	"github.com/suny-am/bitbucket-cli/pkg/cmd/permission"
	"github.com/suny-am/bitbucket-cli/pkg/cmd/pr"
	"github.com/suny-am/bitbucket-cli/pkg/cmd/repo"
	"github.com/suny-am/bitbucket-cli/pkg/cmd/user"
	credentials "github.com/suny-am/bitbucket-cli/pkg/lib/credentials"
	"github.com/suny-am/bitbucket-cli/pkg/types"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bitbucket-cli",
	Short: "CLI solution for interacting with Bitbucket Cloud tenants",
	Long: `This CLI enables shell interaction with various
Bitbucket Cloud resources.

Fetch personal commit history, workspace statistics, branch activity,
Pull Request information and much more, all from your terminal.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		CredProvider := credentials.NewCredentialsProvider()
		credentials, err := CredProvider.GetCredentials()

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		ctx := context.WithValue(cmd.Context(), types.CredentialsKey{}, credentials)
		cmd.SetContext(ctx)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bitbucket-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.AddCommand(commit.CommitCmd)
	rootCmd.AddCommand(permission.PermissionCmd)
	rootCmd.AddCommand(pr.PrCmd)
	rootCmd.AddCommand(repo.RepoCmd)
	rootCmd.AddCommand(user.UserCmd)
}
