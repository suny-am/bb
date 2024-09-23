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

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/internal/iostreams"
	"github.com/suny-am/bitbucket-cli/internal/keyring"
	tablePrinter "github.com/suny-am/bitbucket-cli/internal/tableprinter"
)

type ViewOptions struct {
	repository  string
	workspace   string
	credentials string
}

var opts ViewOptions

var ViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View a repository",
	Long:  `View a repository in a given workspace`,

	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) < 1 {
			return errors.New("<repository> argument is required")
		}

		if len(args) > 1 {
			return errors.New("only one <repository> argument is allowed")
		}

		opts.repository = args[0]
		opts.credentials = cmd.Context().Value(keyring.CredentialsKey{}).(string)

		repo, err := viewRepo(&opts)

		if err != nil {
			return err
		}

		markdown := markdown.Render(repo.Readme, 80, 6)

		fmt.Println(string(markdown))

		tp := tablePrinter.New(os.Stdout, true, 200)
		cs := *iostreams.NewColorScheme(true, true, true)

		repoVal := reflect.ValueOf(*repo)
		vType := repoVal.Type()

		for i := 0; i < repoVal.NumField(); i++ {
			if vType.Field(i).Name == "Readme" {
				continue
			}
			tp.Field(vType.Field(i).Name, tablePrinter.WithColor(cs.Bold))
			v := repoVal.Field(i)
			switch v.Kind() {
			case reflect.Bool:
				tp.Field(fmt.Sprintf("%v", v))
			case reflect.String:
				if v.String() == "" {
					tp.Field("NA", tablePrinter.WithColor(cs.Red))
				} else {
					tp.Field(v.String())
				}
			case reflect.Int:
				tp.Field(fmt.Sprintf("%d", v.Int()))
			case reflect.Struct:
				switch vType.Field(i).Name {
				case "Owner":
					tp.Field(repo.Owner.Display_Name)
				case "Mainbranch":
					tp.Field(repo.Mainbranch.Name)
				case "Project":
					tp.Field(repo.Project.Name)
				}
			}

			tp.EndRow()
		}

		tp.Render()

		return nil
	},
}

func init() {
	ViewCmd.Flags().StringVarP(&opts.workspace, "workspace", "w", "", "Target workspace")
	ViewCmd.MarkFlagRequired("workspace")
}
