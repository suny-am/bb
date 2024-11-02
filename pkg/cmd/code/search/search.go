package search

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/internal/pager"
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
		hlFg := lipgloss.AdaptiveColor{Light: "#000", Dark: "#fff"}
		hlBg := lipgloss.AdaptiveColor{Light: "#ffbb99", Dark: "#aa11dd"}
		mdStyle := lipgloss.NewStyle().Foreground(hlFg).Background(hlBg)

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

		if len(code.Values) == 0 {
			fmt.Println("No results")
			return nil
		}

		pageData := []string{}

		for _, c := range code.Values {

			var ssb strings.Builder

			fileType := strings.Split(c.File.Path, ".")

			ext := fileType[len(fileType)-1]

			ssb.WriteString(fmt.Sprintf("- Path: %s\n", c.File.Path))
			ssb.WriteString(fmt.Sprintf("- Type: %s\n", ext))
			ssb.WriteString(fmt.Sprintf("- Link: %s\n\n", c.File.Links.Html.Href))

			if c.Content_match_count > 0 {
				for i, cm := range c.Content_matches {

					ssb.WriteString(fmt.Sprintf("- Match #%d:\n", i))
					for _, l := range cm.Lines {
						var sssb strings.Builder
						for _, s := range l.Segments {
							var segment string
							if s.Match {
								segment = mdStyle.Render(s.Text)
							} else {
								segment = s.Text
							}
							sssb.WriteString(segment)
						}
						ssb.WriteString(fmt.Sprintf("%s\n", sssb.String()))
					}
					ssb.WriteString("\n")
				}
			}

			if opts.includeSource {
				r, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
				source, err := r.Render(c.File.Source)
				if err == nil {
					ssb.WriteString(fmt.Sprintf("- Source: \n%s", source))
				}
			}

			pageData = append(pageData, ssb.String())
		}

		pager.Draw(pageData)

		return nil
	},
}

func init() {
	var workspaceDefaultValue string
	defaultWorkspace, err := config.GetWorkspace()
	if err != nil {
		SearchCmd.MarkFlagRequired("workspace")
		workspaceDefaultValue = ""
	} else {
		workspaceDefaultValue = defaultWorkspace
	}

	SearchCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
	SearchCmd.Flags().StringVarP(&opts.repository, "repo", "r", "", "Target repository")
	SearchCmd.Flags().IntVarP(&opts.limit, "limit", "l", 0, "Result limit")
	SearchCmd.Flags().BoolVarP(&opts.includeSource, "source", "s", false, "Include source")
}
