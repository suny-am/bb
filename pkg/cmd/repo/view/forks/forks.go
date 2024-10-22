package forks

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/internal/keyring"
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

		/*
			forks, err := viewforks(&opts)
			if err != nil {
				return err
			}
		*/
		return nil
	},
}

func init() {
	ForksCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ForksCmd.MarkFlagRequired("workspace")
}
