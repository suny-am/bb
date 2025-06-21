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

func getForks(opts *api.ForkListptions, cmd *cobra.Command) (*api.Repositories, error) {
	var forks api.Repositories
	var err error

	go func() {
		err = get(&forks, cmd, opts)
		debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
		if debug {
			textinput.ConfirmKey()
		}
		spinner.Stop()
	}()

	spinner.Start("Searching forks")

	return &forks, err
}

func get(forks *api.Repositories, cmd *cobra.Command, opts *api.ForkListptions) error {
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

	if err := json.Unmarshal([]byte(body), &forks); err != nil {
		return err
	}
	return nil
}

func generateRequest(opts *api.ForkListptions) (*http.Request, error) {
	authHeaderValue := fmt.Sprintf("Basic %s", opts.Credentials)
	endpoint := fmt.Sprintf("%s/forks", http2.DetermineRepositoryEndpoint(opts))
	forksReq, err := http.NewRequest("GET", endpoint, nil)

	forksReq.Header.Add("Accept", "application/json")
	forksReq.Header.Add("Authorization", authHeaderValue)

	return forksReq, err
}
