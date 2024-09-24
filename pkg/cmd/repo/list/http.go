package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/suny-am/bitbucket-cli/api"
)

func listRepos(opts *ListOptions) (*api.Repositories, error) {

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)

	endpoint := "https://api.bitbucket.org/2.0/repositories"

	if opts.workspace != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, opts.workspace)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeaderValue)

	if opts.limit > 0 {
		query := req.URL.Query()
		query.Add("pagelen", strconv.Itoa(opts.limit))
		req.URL.RawQuery = query.Encode()
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var repositories api.Repositories

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(body), &repositories); err != nil {
		return nil, err
	}

	return &repositories, nil
}
