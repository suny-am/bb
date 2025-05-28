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
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/git"
	"github.com/suny-am/bb/internal/keyring"
)

type ViewOptions struct {
	repository  string
	workspace   string
	pullrequest string
	credentials string
}

var opts ViewOptions

var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a pullrequest",
	Long:  `View a pullrequest in a given workspace`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.repository == "" {
			opts.repository = git.GetGitRepo()
		}

		if len(args) < 1 {
			return errors.New("<pullrequest> argument is required")
		}

		if len(args) > 1 {
			return errors.New("only one <pullrequest> argument is allowed")
		}

		opts.pullrequest = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
		pullrequest, err := getPullrequest(&opts, cmd)
		if err != nil {
			return err
		}

		if pullrequest.Title == "" {
			fmt.Println(api.NoResults)
			return nil
		}

		viewPullrequest(pullrequest)

		return nil
	},
}

type comment struct {
	timestamp string
	content   string
}

var (
	sb      strings.Builder
	tsStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#003300", Dark: "#11bb99"})
	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#003300", Dark: "#00ffff"})

	mdStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingTop(1)
)

func colorAttribute(key string, value string) string {
	return fmt.Sprintf("%s: %s\n", keyStyle.Render(key), value)
}

func viewPullrequest(pr *api.Pullrequest) {
	re, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(),
		glamour.WithEmoji())

	sb.WriteString(colorAttribute("Title", pr.Source.Branch.Name))
	sb.WriteString(colorAttribute("Author", fmt.Sprintf("%s [%s]", pr.Author.Display_Name, pr.Author.Nickname)))
	sb.WriteString(colorAttribute("Created", pr.Created_On))
	sb.WriteString(colorAttribute("State", pr.State))
	sb.WriteString(colorAttribute("Link", pr.Links.Html.Href))

	if pr.Comment_Count > 0 {
		count := lipgloss.NewStyle().
			Width(100).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			Render(colorAttribute("Comments", fmt.Sprintf("%d", pr.Comment_Count)))
		sb.WriteString(count)
		commentData := []comment{}
		for _, c := range pr.Comments.Values {
			comment := comment{}
			var ssb strings.Builder
			tsString := tsStyle.Render(c.Created_On)
			ssb.WriteString(fmt.Sprintf("\nAuthor: %s [%s] [%s]\n", c.User.Display_Name, c.User.Nickname, tsString))
			content, err := re.Render(c.Content.Raw)
			if err == nil {
				ssb.WriteString(fmt.Sprintf("%s\n\n", content))
			}
			comment.timestamp = c.Created_On
			comment.content = ssb.String()
			commentData = append(commentData, comment)
		}

		for _, c := range commentData {
			sb.WriteString(c.content)
		}
	}

	fmt.Println(mdStyle.Render(sb.String()))
}

func init() {
	var workspaceDefaultValue string
	defaultWorkspace, err := config.GetWorkspace()
	if err != nil {
		_ = ViewCmd.MarkFlagRequired("workspace")
		workspaceDefaultValue = ""
	} else {
		workspaceDefaultValue = defaultWorkspace
	}

	ViewCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", workspaceDefaultValue, "Target workspace")
	ViewCmd.Flags().StringVarP(&opts.repository, "repository", "r", "", "Target repository")
}
