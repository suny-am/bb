/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	credentials "github.com/suny-am/bitbucket-cli/pkg/utils"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "bitbucket-cli",
	Short: "CLI solution for interacting with Bitbucket Cloud tenants",
	Long: `This CLI enables shell interaction with various
Bitbucket Cloud resources.

Fetch personal commit history, workspace statistics, branch activity,
Pull Request information and much more, all from your terminal.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

var Credentials string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
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

	CredProvider := credentials.NewCredentialsProvider()
	credentials, err := CredProvider.GetCredentials()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Credentials = credentials
}
