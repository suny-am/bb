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

	"github.com/suny-am/bitbucket-cli/api"
)

func viewRepo(opts *ViewOptions) (*api.Repository, error) {

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)
	endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", opts.workspace, opts.repository)

	client := &http.Client{}

	repoReq, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	repoReq.Header.Add("Accept", "application/json")
	repoReq.Header.Add("Authorization", authHeaderValue)

	resp, err := client.Do(repoReq)

	if err != nil {
		return nil, err
	}

	var repository api.Repository

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(body), &repository); err != nil {
		return nil, err
	}

	endpoint = fmt.Sprintf("%s/src/master/README.md", endpoint)

	readmeReq, _ := http.NewRequest("GET", endpoint, nil)
	readmeReq.Header.Add("Authorization", authHeaderValue)
	readmeResp, _ := client.Do(readmeReq)
	readmeBody, _ := io.ReadAll(readmeResp.Body)

	if readmeText := string(readmeBody); readmeText != "" {
		repository.Readme = readmeText
	}

	return &repository, nil
}
