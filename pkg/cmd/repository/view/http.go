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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/http2"
	"github.com/suny-am/bb/internal/spinner"
	"github.com/suny-am/bb/internal/textinput"
)

func getRepo(opts *ViewOptions, cmd *cobra.Command) (*api.Repository, error) {
	var repository api.Repository
	var err error

	go func() {
		err = get(&repository, cmd, opts)
		debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
		if debug {
			textinput.ConfirmKey()
		}
		spinner.Stop()
	}()

	spinner.Start("Searching repositories")

	return &repository, err
}

func get(repository *api.Repository, cmd *cobra.Command, opts *ViewOptions) error {
	client := http2.Init(cmd)
	req, err := generateRequest(opts)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(body), &repository); err != nil {
		return err
	}

	req.URL, err = req.URL.Parse(fmt.Sprintf("%s/src/master/README", req.URL))
	if err != nil {
		return err
	}

	readmeResp, err := client.Do(req)

	var readmeBody []byte
	if err != nil || readmeResp.StatusCode != 200 {
		req.URL, err = url.Parse(strings.Replace(req.URL.String(), "README", "README.md", 1))
		if err != nil {
			return err
		}

		readmeResp, err = client.Do(req)
		if err != nil || readmeResp.StatusCode != 200 {
			readmeBody = nil
		} else {
			readmeBody, _ = io.ReadAll(readmeResp.Body)
		}
	}

	if readmeBody != nil {
		if readmeText := string(readmeBody); readmeText != "" {
			repository.Readme = readmeText
		}
	}

	return nil
}

func generateRequest(opts *ViewOptions) (*http.Request, error) {
	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)
	endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", opts.workspace, opts.repository)
	req, err := http.NewRequest("GET", endpoint, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeaderValue)

	return req, err
}
