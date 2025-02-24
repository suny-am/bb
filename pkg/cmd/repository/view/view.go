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
	"github.com/suny-am/bb/internal/keyring"
	"github.com/suny-am/bb/internal/textview"
	"github.com/suny-am/bb/internal/util"
)

type ViewOptions struct {
	repository  string
	workspace   string
	credentials string
	current     bool
}

var opts ViewOptions

var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a repository",
	Long:  `View a repository in a given workspace`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if opts.current {
			opts.repository = util.GetCurrentDir()
		} else {

			if len(args) < 1 {
				return errors.New("<repository> argument is required")
			}

			if len(args) > 1 {
				return errors.New("only one <repository> argument is allowed")
			}
		}

		opts.repository = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)
		repo, err := getRepo(&opts, cmd)
		if err != nil {
			return err
		}

		if repo.Name == "" {
			fmt.Println(api.NoResults)
			return nil
		}

		viewRepo(repo)

		return nil
	},
}

func colorAttribute(key string, value string) string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3d8280")).PaddingLeft(2).Render(key)
	valueStyle := lipgloss.NewStyle().PaddingBottom(1).Render(value)
	return strings.Join([]string{keyStyle, valueStyle}, ": ")
}

func viewRepo(repo *api.Repository) {
	var status string
	if repo.Is_Private {
		status = "Private"
	} else {
		status = "Public"
	}
	content := []string{
		colorAttribute("Owner", fmt.Sprintf("%s <%s>", repo.Owner.Display_Name, repo.Owner.Nickname)),
		colorAttribute("Size", fmt.Sprintf("%d Kb", repo.Size/1000)),
		colorAttribute("Language", repo.Language),
		colorAttribute("Project", fmt.Sprintf("%s [%s]", repo.Project.Name, repo.Project.Type)),
		colorAttribute("Created", repo.Created_On),
		colorAttribute("Updated", repo.Updated_On),
		colorAttribute("Status", status),
		colorAttribute("Main branch", fmt.Sprintf("%s [%s]", repo.Mainbranch.Name, repo.Mainbranch.Type)),
		colorAttribute("Links", repo.Links.Html.Href),
	}

	if repo.Description != "" {
		description, err := glamour.Render(repo.Description, "light")
		if err != nil {
			panic(err)
		}
		content = append(content, description)
	}

	if repo.Readme != "" {
		readme, err := glamour.Render(repo.Readme, "light")
		if err != nil {
			panic(err)
		}
		content = append(content, readme)
	}
	textview.DrawView(repo.Name, strings.Join(content, "\n"))
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
	ViewCmd.Flags().BoolVarP(&opts.current, "current", "c", true, "Use current directory name as repository name")
}
