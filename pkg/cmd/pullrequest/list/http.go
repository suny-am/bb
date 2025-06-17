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
	"math"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/http2"
	"github.com/suny-am/bb/internal/spinner"
	"github.com/suny-am/bb/internal/textinput"
)

const pageLenMax = 50

func getPullrequests(opts *api.PullrequestListOptions, cmd *cobra.Command) (*api.Pullrequests, error) {
	var pullrequests api.Pullrequests
	var err error

	go func() {
		err = get(&pullrequests, cmd, opts)
		debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
		if debug {
			textinput.ConfirmKey()
		}
		spinner.Stop()
	}()

	spinner.Start("Searching pullrequests")

	return &pullrequests, err
}

func get(pullrequests *api.Pullrequests, cmd *cobra.Command, opts *api.PullrequestListOptions) error {
	client := http2.Init(cmd)

	req, err := generateRequest(opts)
	if err != nil {
		return err
	}
	fetchPullrequestsRecurse(pullrequests, client, req)
	return nil
}

func generateRequest(opts *api.PullrequestListOptions) (*http.Request, error) {
	var endpoint string
	authHeaderValue := fmt.Sprintf("Basic %s", opts.Credentials)

	if opts.Repository == "" {

		user, err := config.GetUsername()
		if err != nil {
			return nil, err
		}
		endpoint = fmt.Sprintf("%s/pullrequests/%s",
			http2.DetermineWorkspaceEndpoint(opts), user)
	} else {

		if opts.Workspace == "" {
			workspace, err := config.GetWorkspace()
			if err != nil {
				return nil, err
			}
			opts.Workspace = workspace
		}

		endpoint = fmt.Sprintf("%s/pullrequests", http2.DetermineRepositoryEndpoint(opts))
	}

	if opts.Title != "" {
		opts.Title = fmt.Sprintf("q=title~\"%s\"", opts.Title)
	}
	if opts.Creator != "" {
		opts.Creator = fmt.Sprintf("q=author.nickname~\"%s\"", opts.Creator)
	}
	if opts.State != "" {
		endpoint = fmt.Sprintf("q=state=\"%s\"", opts.State)
	}

	opts.PageLen = int(math.Min(float64(opts.PageLen), float64(pageLenMax)))

	endpoint = http2.DetermineQueryParametersDirect(endpoint, []string{
		opts.Title,
		opts.Creator,
		opts.State,
		fmt.Sprintf("pageLen=%d", opts.PageLen),
	})

	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", endpointUrl.String(), nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeaderValue)

	fmt.Println(endpointUrl.String())

	return req, err
}

func fetchPullrequestsRecurse(pullrequests *api.Pullrequests, client *http2.Client, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var particalPullrequests api.Pullrequests

	if err := json.Unmarshal([]byte(body), &particalPullrequests); err != nil {
		panic(err)
	}

	for i := range particalPullrequests.Values {
		p := &particalPullrequests.Values[i]
		link := p.Links.Self.Href

		singlePrReq, err := http.NewRequest("GET", link, nil)
		singlePrReq.Header.Add("Authorization", req.Header["Authorization"][0])
		if err == nil {
			singlePrResp, err := client.Do(singlePrReq)
			if err == nil {
				body, err := io.ReadAll(singlePrResp.Body)
				if err == nil {
					var singlePullrequest api.Pullrequest
					if err := json.Unmarshal([]byte(body), &singlePullrequest); err == nil {
						p.Participants = append(p.Participants, singlePullrequest.Participants...)
					}
				}
			}
		}
	}

	if particalPullrequests.Values != nil {
		pullrequests.Values = append(pullrequests.Values, particalPullrequests.Values...)
		if particalPullrequests.Next != "" {
			newReq, err := http.NewRequest("GET", particalPullrequests.Next, nil)
			newReq.Header.Add("Authorization", req.Header["Authorization"][0])
			newReq.Header.Add("Accept", req.Header["Accept"][0])
			if err != nil || len(pullrequests.Values) >= opts.PageLen {
				fmt.Println(err)
				return
			}
			fetchPullrequestsRecurse(pullrequests, client, newReq)
		}
	}
}
