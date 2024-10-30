package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/suny-am/bb/api"
)

func listRepos(opts *ListOptions) (*api.Repositories, error) {
	client := &http.Client{}
	var repositories api.Repositories

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)
	endpoint := "https://api.bitbucket.org/2.0/repositories"

	if opts.workspace != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, opts.workspace)
	}

	if opts.nameFilter != "" {
		endpoint = fmt.Sprintf("%s?q=name~\"%s\"", endpoint, opts.nameFilter)
	}

	var pageLength int

	if opts.limit > 100 {
		pageLength = 100
	} else {
		pageLength = opts.limit
	}

	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if opts.limit > 0 {
		query := endpointUrl.Query()
		query.Add("pagelen", strconv.Itoa(pageLength))
		endpointUrl.RawQuery = query.Encode()
	}

	req, err := http.NewRequest("GET", endpointUrl.String(), nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeaderValue)

	fetchReposRecurse(client, req, &repositories)

	if err != nil {
		return nil, err
	}

	return &repositories, nil
}

func fetchReposRecurse(client *http.Client, req *http.Request, repositories *api.Repositories) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request Error: %s", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Body Parsing Error: %s", err)
		return
	}

	var partialRepositories api.Repositories

	if err := json.Unmarshal([]byte(body), &partialRepositories); err != nil {
		fmt.Printf("Unmarshalling error Error: %s", err)
		return
	}

	if partialRepositories.Values != nil {
		repositories.Values = append(repositories.Values, partialRepositories.Values...)
		if partialRepositories.Next != "" {
			newReq, err := http.NewRequest("GET", partialRepositories.Next, nil)
			newReq.Header.Add("Authorization", req.Header["Authorization"][0])
			newReq.Header.Add("Accept", req.Header["Accept"][0])
			if err != nil {
				fmt.Printf("Request Error: %s", err)
				return
			}
			if len(repositories.Values) >= opts.limit {
				return
			}
			fetchReposRecurse(client, newReq, repositories)
		}
	}
}
