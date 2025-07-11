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

func getPipelines(opts *api.PipelineListOptions, cmd *cobra.Command) (*api.Pipelines, error) {
	var pipelines api.Pipelines
	var err error

	go func() {
		err = get(&pipelines, cmd, opts)
		debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
		if debug {
			textinput.ConfirmKey()
		}
		spinner.Stop()
	}()

	spinner.Start("Getting pipelines")

	return &pipelines, err
}

func get(pipelines *api.Pipelines, cmd *cobra.Command, opts *api.PipelineListOptions) error {
	client := http2.Init(cmd)
	req, err := generateRequest(opts)
	if err != nil {
		return err
	}

	fetchPipelinesRecurse(client, req, pipelines)

	return nil
}

func fetchPipelinesRecurse(client *http2.Client, req *http.Request, pipelines *api.Pipelines) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var partialpipelines api.Pipelines

	if err := json.Unmarshal([]byte(body), &partialpipelines); err != nil {
		panic(err)
	}

	if partialpipelines.Values != nil {
		pipelines.Values = append(pipelines.Values, partialpipelines.Values...)
		if partialpipelines.Next != "" {
			newReq, err := http.NewRequest("GET", partialpipelines.Next, nil)
			newReq.Header.Add("Authorization", req.Header["Authorization"][0])
			newReq.Header.Add("Accept", req.Header["Accept"][0])
			if err != nil {
				panic(err)
			}
			if len(pipelines.Values) >= opts.PageLen {
				return
			}
			fetchPipelinesRecurse(client, newReq, pipelines)
		}
	}
}

func generateRequest(opts *api.PipelineListOptions) (*http.Request, error) {
	authHeaderValue := fmt.Sprintf("Basic %s", opts.Credentials)

	endpoint := fmt.Sprintf("%s/pipelines",
		http2.DetermineRepositoryEndpoint(opts))

	endpoint = http2.DetermineQueryParameters(opts, endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeaderValue)

	return req, err
}
