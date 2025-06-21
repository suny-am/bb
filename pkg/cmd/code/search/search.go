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
package search

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/internal/pager"
)

var (
	opts    api.CodeSearchOptions
	hlFg    = lipgloss.AdaptiveColor{Light: "#000", Dark: "#fff"}
	hlBg    = lipgloss.AdaptiveColor{Light: "#ffbb99", Dark: "#aa11dd"}
	mdStyle = lipgloss.NewStyle().Foreground(hlFg).Background(hlBg)
)

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

		opts.Search_Query = args[0]
		opts.Credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
		code, err := searchCode(&opts, cmd)
		if err != nil {
			fmt.Println(err)
		}

		if len(code.Values) == 0 {
			fmt.Println(api.NoResults)
			return nil
		}

		displaySearchResults(code)

		return nil
	},
}

func init() {
	var workspaceDefaultValue string
	defaultWorkspace, err := config.GetWorkspace()
	if err != nil {
		_ = SearchCmd.MarkFlagRequired("workspace")
		workspaceDefaultValue = ""
	} else {
		workspaceDefaultValue = defaultWorkspace
	}

	SearchCmd.Flags().StringVarP(&opts.Workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
	SearchCmd.Flags().StringVarP(&opts.Repository, "repo", "r", "", "Target repository")
	SearchCmd.Flags().IntVarP(&opts.PageLen, "limit", "l", 0, "Result limit")
	SearchCmd.Flags().BoolVarP(&opts.IncludeSource, "source", "s", false, "Include source")
}

func displaySearchResults(code *api.CodeSearchResponse) {
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

		if opts.IncludeSource {
			r, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
			source, err := r.Render(c.File.Source)
			if err == nil {
				ssb.WriteString(fmt.Sprintf("- Source: \n%s", source))
			}
		}

		pageData = append(pageData, ssb.String())
	}

	pager.Draw(pageData)
}
