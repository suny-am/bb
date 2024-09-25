package forks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/suny-am/bitbucket-cli/api"
)

func viewforks(opts *ForksOptions) (*api.Repositories, error) {

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)
	endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/forks", opts.workspace, opts.repository)

	client := &http.Client{}

	forksReq, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	forksReq.Header.Add("Accept", "application/json")
	forksReq.Header.Add("Authorization", authHeaderValue)

	resp, err := client.Do(forksReq)

	if err != nil {
		return nil, err
	}

	var forks api.Repositories

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(body), &forks); err != nil {
		return nil, err
	}

	return &forks, nil
}
