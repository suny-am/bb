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
package list

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/api"
	"github.com/suny-am/bitbucket-cli/internal/bb"
	"github.com/suny-am/bitbucket-cli/lib/iostreams"
	tablePrinter "github.com/suny-am/bitbucket-cli/lib/tableprinter"
	cmdUtil "github.com/suny-am/bitbucket-cli/pkg/cmdutil"
)

type ListOptions struct {
	HttpClient func() (*http.Client, error)
	Config     func() (bb.Config, error)
	IO         *iostreams.IOStreams

	Limit int
	Owner string

	Visibility string

	Now func() time.Time
}

func NewCmdList(f *cmdUtil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := ListOptions{
		IO:         f.IOStreams,
		Config:     f.Config,
		HttpClient: f.HttpClient,
		Now:        time.Now,
	}

	var (
		flagPublic  bool
		flagPrivate bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List repositories",
		Long:  `List one or more personal and/or workspace repositories`,

		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Limit < 1 {
				return cmdUtil.FlagErrorf("invalid limit: %v", opts.Limit)
			}

			if err := cmdUtil.MutuallyExclusive("specify only one of `--public`, `--private`, or `--visibility`", flagPublic, flagPrivate, opts.Visibility != ""); err != nil {
				return err
			}

			if flagPrivate {
				opts.Visibility = "private"
			} else if flagPublic {
				opts.Visibility = "public"
			}

			if len(args) > 0 {
				opts.Owner = args[0]
			}

			if runF != nil {
				return runF(&opts)
			}

			return listRun(&opts)

		},
	}

	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of repositories to list")
	cmdUtil.StringEnumFlag(cmd, &opts.Visibility, "visibility", "", "", []string{"public", "private", "internal"}, "Filter by repository visibility")

	cmd.Flags().BoolVar(&flagPrivate, "private", false, "Show only private repositories")
	cmd.Flags().BoolVar(&flagPublic, "public", false, "Show only public repositories")
	_ = cmd.Flags().MarkDeprecated("public", "use `--visibility=public` instead")
	_ = cmd.Flags().MarkDeprecated("private", "use `--visibility=private` instead")

	return cmd
}

func listRun(opts *ListOptions) error {

	limit := opts.Limit

	credentials := "test"

	authHeaderData := fmt.Sprintf("Basic %s", credentials)

	client := resty.New()

	// TBD add workspace as argument

	endpoint := "https://api.bitbucket.org/2.0/repositories"

	resp, err := client.R().
		SetHeader("Authorization", authHeaderData).
		SetHeader("Accept", "application/json").
		SetQueryParam("pagelen", strconv.Itoa(limit)).
		EnableTrace().
		Get(endpoint)

	if resp.IsError() {
		fmt.Println(err.Error())
	}

	if resp.IsSuccess() {
		var repositories api.Repositories

		if err := json.Unmarshal([]byte(resp.String()), &repositories); err != nil {
			fmt.Println(err)
		}

		tp := tablePrinter.New(os.Stdout, true, 200) // TBD pass dynamic opts

		cs := *iostreams.NewColorScheme(true, true, true)

		headers := []string{"NAME", "INFO", "UPDATED"}
		tp.Header(headers, tablePrinter.WithColor(cs.LightGrayUnderline))
		for i := range repositories.Values {
			repo := repositories.Values[i]
			tp.Field(repo.Full_Name, tablePrinter.WithColor(cs.Bold))
			tp.Field("public", tablePrinter.WithColor(cs.Gray))
			tp.Field(repo.Updated_On, tablePrinter.WithColor(cs.Gray))
			tp.EndRow()
		}

		tp.Render()
	}
	return nil
}
