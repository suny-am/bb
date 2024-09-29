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
	"os"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/api"
	"github.com/suny-am/bitbucket-cli/internal/iostreams"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
	"github.com/suny-am/bitbucket-cli/internal/markdown"
	tablePrinter "github.com/suny-am/bitbucket-cli/internal/tableprinter"
	"github.com/suny-am/bitbucket-cli/pkg/cmd/repo/view/forks"
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

		if len(args) < 1 {
			return errors.New("<pullrequest> argument is required")
		}

		if len(args) > 1 {
			return errors.New("only one <pullrequest> argument is allowed")
		}

		opts.pullrequest = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		repo, err := viewPullrequest(&opts)

		if err != nil {
			return err
		}

		renderFields(*repo)

		return nil
	},
}

func init() {
	ViewCmd.AddCommand(forks.ForksCmd)

	ViewCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ViewCmd.Flags().StringVarP(&opts.repository, "repo", "r", "", "Target repository")
	ViewCmd.MarkFlagRequired("workspace")
	ViewCmd.MarkFlagRequired("repo")
}

func renderFields(pr api.Pullrequest) {

	tp := tablePrinter.New(os.Stdout, true, 200)
	cs := *iostreams.NewColorScheme(true, true, true)

	prVal := reflect.ValueOf(pr)
	vType := prVal.Type()

	for i := 0; i < prVal.NumField(); i++ {
		if vType.Field(i).Name == "Readme" ||
			vType.Field(i).Name == "Links" ||
			vType.Field(i).Name == "Rendered" ||
			vType.Field(i).Name == "Comments" ||
			vType.Field(i).Name == "Summary" {
			continue
		}
		v := prVal.Field(i)
		tp.Field(vType.Field(i).Name, tablePrinter.WithColor(cs.Bold))
		switch v.Kind() {
		case reflect.Bool:
			tp.Field(fmt.Sprintf("%v", v))
		case reflect.String:
			if v.String() == "" {
				tp.Field("NA", tablePrinter.WithColor(cs.Red))
			} else {
				if vType.Field(i).Name == "State" {
					tp.Field(v.String(), tablePrinter.WithColor(cs.GreenBold))
				} else {
					tp.Field(v.String())
				}
			}
		case reflect.Int:
			tp.Field(fmt.Sprintf("%d", v.Int()))
		case reflect.Struct:
			switch vType.Field(i).Name {
			case "Author":
				tp.Field(pr.Author.Display_Name)
			case "Closed_By":
				tp.Field(pr.Closed_By.Name)
			case "Merge_Commit":
				tp.Field(pr.Merge_Commit.Hash)
			}
		case reflect.Slice:
			switch vType.Field(i).Name {
			case "Reviewers":
				var entryString string
				for _, r := range pr.Reviewers {
					if entryString == "" {
						entryString = r.Nickname
					} else {
						entryString = fmt.Sprintf("%s, %s", entryString, r.Display_Name)
					}
				}
				tp.Field(entryString)
			case "Participants":
				var entryString string
				for _, p := range pr.Participants {
					if entryString == "" {
						entryString = p.User.Display_Name
					}
					entryString = fmt.Sprintf("%s, %s", entryString, p.User.Display_Name)
				}

				tp.Field(entryString)
			}
		}
		tp.EndRow()
	}

	tp.Render()
	// list comments

	fmt.Println("")

	for i, c := range pr.Comments.Values {

		md := c.User.Display_Name
		inline := "NA"
		if c.Inline.Path != "" {
			inline = c.Inline.Path
		}
		md = fmt.Sprintf("- [Comment #%d] **%s**\n\nFile: `%s`\n\nLine: %d\n\n%s\n\n",
			i+1,
			md,
			inline,
			c.Inline.To,
			c.Content.Raw)
		markdown.Render(md)
	}

}
