package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/suny-am/bitbucket-cli/api"
)

func listRepos(opts *ListOptions) (*api.RepoListResponse, error) {

	// authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)

	endpoint := "https://api.bitbucket.org/2.0/repositories"

	if opts.workspace != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, opts.workspace)

		// --repository requires --workspace
		if opts.repository != "" {
			endpoint = fmt.Sprintf("%s/%s", endpoint, opts.repository)
		}
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	// req.Header.Add("Authorize", authHeaderValue)
	req.Header.Add("Accept", "application/json")

	query := req.URL.Query()
	if opts.limit != nil {
		query.Add("size", strconv.Itoa(*opts.limit))
		req.URL.RawQuery = query.Encode()
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var response api.RepoListResponse

	fmt.Println(resp.Body)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return nil, err
	}

	return &response, nil
}
