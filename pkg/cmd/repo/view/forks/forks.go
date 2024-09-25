package forks

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/internal/iostreams"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
	tablePrinter "github.com/suny-am/bitbucket-cli/internal/tableprinter"
)

type ForksOptions struct {
	repository  string
	workspace   string
	credentials string
}

var opts ForksOptions

var ForksCmd = &cobra.Command{
	Use:   "forks",
	Short: "View forks for a repository",
	Long:  `View one ore more forks for a given repository`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("<repository> argument is required")
		}

		if len(args) > 1 {
			return errors.New("only one <repository> argument is allowed")
		}

		opts.repository = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		forks, err := viewforks(&opts)

		if err != nil {
			return err
		}

		tp := tablePrinter.New(os.Stdout, true, 500)
		cs := *iostreams.NewColorScheme(true, true, true)

		headers := []string{"NAME", "INFO", "UPDATED"}
		tp.Header(headers, tablePrinter.WithColor(cs.LightGrayUnderline))
		for i := range forks.Values {
			repo := forks.Values[i]
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
	ForksCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ForksCmd.MarkFlagRequired("workspace")
}
