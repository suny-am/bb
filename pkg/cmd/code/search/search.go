package search

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
	"github.com/suny-am/bitbucket-cli/internal/markdown"
)

type SearchOptions struct {
	repository    string
	workspace     string
	credentials   string
	searchParam   string
	includeSource bool
	limit         int
}

var opts SearchOptions

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for code",
	Long:  `Search for code in a workspace or repository`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("<searchParam> required")
		}

		if len(args) > 1 {
			return errors.New("only 1 <searchParam> allowed")
		}

		opts.searchParam = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		code, err := searchCode(&opts)

		if err != nil {
			fmt.Println(err)
		}

		for _, c := range code.Values {
			markdown.Render(fmt.Sprintf("- **%s**\n- ```%s````", c.File.Path, c.File.Links.Html.Href))
			if opts.includeSource {
				markdown.Render(c.File.Source)
			}
		}

		return nil
	},
}

func init() {
	SearchCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	SearchCmd.Flags().StringVarP(&opts.repository, "repo", "r", "", "Target repository")
	SearchCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Result limit")
	SearchCmd.Flags().BoolVarP(&opts.includeSource, "source", "s", false, "Include source")
	SearchCmd.MarkFlagRequired("workspace")
}
