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
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/http2"
	"github.com/suny-am/bb/internal/spinner"
	"github.com/suny-am/bb/internal/textinput"
)

func getRepos(opts *api.RepositoryListOptions, cmd *cobra.Command) (*api.Repositories, error) {
	var repositories api.Repositories
	var err error

	go func() {
		err = get(&repositories, cmd, opts)
		debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
		if debug {
			textinput.ConfirmKey()
		}
		spinner.Stop()
	}()

	spinner.Start("Searching repositories")

	return &repositories, err
}

func get(repositories *api.Repositories, cmd *cobra.Command, opts *api.RepositoryListOptions) error {
	client := http2.Init(cmd)
	req, err := generateRequest(opts)
	if err != nil {
		return err
	}
	fetchReposRecurse(repositories, client, req)
	return nil
}

func fetchReposRecurse(repositories *api.Repositories, client *http2.Client, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var partialRepositories api.Repositories

	if err := json.Unmarshal([]byte(body), &partialRepositories); err != nil {
		panic(err)
	}

	if partialRepositories.Values != nil {
		repositories.Values = append(repositories.Values, partialRepositories.Values...)
		if partialRepositories.Next != "" {
			req.URL, err = req.URL.Parse(partialRepositories.Next)
			if err != nil {
				fmt.Printf("Error parsing URL: %s", err)
				return
			}
			if len(repositories.Values) >= opts.PageLen {
				return
			}
			fetchReposRecurse(repositories, client, req)
		}
	}
}

func generateRequest(opts *api.RepositoryListOptions) (*http.Request, error) {
	authHeaderValue := fmt.Sprintf("Basic %s", opts.Credentials)
	endpoint := "https://api.bitbucket.org/2.0/repositories"

	if opts.Workspace != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, opts.Workspace)
	}

	if opts.Name != "" {
		endpoint = fmt.Sprintf("%s?q=name~\"%s\"", endpoint, opts.Name)
		endpoint = http2.DetermineQueryParametersDirect(endpoint, []string{
			fmt.Sprintf("q=name~\"%s\"", opts.Name),
			fmt.Sprintf("sort=%s", opts.Sort),
			fmt.Sprintf("pagelen=%d", opts.PageLen),
		})
	}

	req, err := http.NewRequest("GET", endpoint, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeaderValue)

	return req, err
}
