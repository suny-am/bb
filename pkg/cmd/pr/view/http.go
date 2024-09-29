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
	"strings"

	"github.com/suny-am/bitbucket-cli/api"
)

func viewPullrequest(opts *ViewOptions) (*api.Pullrequest, error) {

	client := &http.Client{}

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)
	endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests", opts.workspace, opts.repository)
	idEndpoint := fmt.Sprintf("%s?q=title~\"%s\"", endpoint, opts.pullrequest)
	idEndpoint = strings.ReplaceAll(idEndpoint, " ", "%20")

	prIdReq, err := http.NewRequest("GET", idEndpoint, nil)
	if err != nil {
		return nil, err
	}

	prIdReq.Header.Add("Accept", "application/json")
	prIdReq.Header.Add("Authorization", authHeaderValue)
	resp, err := client.Do(prIdReq)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pullrequests api.Pullrequests
	var pullrequest api.Pullrequest
	if err := json.Unmarshal([]byte(body), &pullrequests); err != nil {
		return nil, err
	}

	prEndpoint := fmt.Sprintf("%s/%d", endpoint, pullrequests.Values[0].Id)
	prReq, err := http.NewRequest("GET", prEndpoint, nil)
	prReq.Header.Add("Accept", "application/json")
	prReq.Header.Add("Authorization", authHeaderValue)
	if err != nil {
		return nil, err
	}

	resp, err = client.Do(prReq)
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(body), &pullrequest); err != nil {
		return nil, err
	}

	commentsReq, err := http.NewRequest("GET", pullrequest.Links.Comments.Href, nil)
	if err != nil {
		return nil, err
	}

	commentsReq.Header = prReq.Header
	resp, _ = client.Do(commentsReq)
	body, _ = io.ReadAll(resp.Body)

	var comments api.Comments
	json.Unmarshal([]byte(body), &comments)
	pullrequest.Comments = comments

	return &pullrequest, nil
}
