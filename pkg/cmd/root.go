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
	"os"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/pkg/cmd/code"
	"github.com/suny-am/bb/pkg/cmd/configure"
	"github.com/suny-am/bb/pkg/cmd/pipeline"
	"github.com/suny-am/bb/pkg/cmd/pullrequest"
	"github.com/suny-am/bb/pkg/cmd/repository"
)

var rootCmd = &cobra.Command{
	Use:   "bb",
	Short: "CLI solution for interacting with Bitbucket Cloud tenants",
	Long: `This CLI enables shell interaction with various
Bitbucket Cloud resources.

Fetch personal commit history, workspace statistics, branch activity,
Pull Request information and much more, all from your terminal.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		CredProvider := keyring.NewCredentialsProvider()
		credentials, err := CredProvider.GetCredentials()

		ctx := context.WithValue(cmd.Context(), keyring.CredentialsKey{}, credentials)
		cmd.SetContext(ctx)

		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//	rootCmd.AddCommand(fork.ForkCmd)
	rootCmd.AddCommand(code.CodeCmd)
	rootCmd.AddCommand(configure.ConfigureCmd)
	rootCmd.AddCommand(pipeline.PipelineCmd)
	rootCmd.AddCommand(pullrequest.PullrequestCmd)
	rootCmd.AddCommand(repository.RepositoryCmd)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug mode")
}
