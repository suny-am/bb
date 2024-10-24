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

func getPipelines(opts *ListOptions) (*api.Pipelines, error) {
	client := &http.Client{}
	var pipelines api.Pipelines

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)
	endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pipelines", opts.workspace, opts.repository)

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

	endpoint = fmt.Sprintf("%s?sort=-created_on", endpointUrl.String())

	req, err := http.NewRequest("GET", endpoint, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", authHeaderValue)

	fetchPipelinesRecurse(client, req, &pipelines)

	if err != nil {
		return nil, err
	}

	return &pipelines, nil
}

func fetchPipelinesRecurse(client *http.Client, req *http.Request, pipelines *api.Pipelines) {
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

	var partialpipelines api.Pipelines

	if err := json.Unmarshal([]byte(body), &partialpipelines); err != nil {
		fmt.Printf("Unmarshalling error Error: %s", err)
		return
	}

	if partialpipelines.Values != nil {
		pipelines.Values = append(pipelines.Values, partialpipelines.Values...)
		if partialpipelines.Next != "" {
			newReq, err := http.NewRequest("GET", partialpipelines.Next, nil)
			newReq.Header.Add("Authorization", req.Header["Authorization"][0])
			newReq.Header.Add("Accept", req.Header["Accept"][0])
			if err != nil {
				fmt.Printf("Request Error: %s", err)
				return
			}
			if len(pipelines.Values) >= opts.limit {
				return
			}
			fetchPipelinesRecurse(client, newReq, pipelines)
		}
	}
}
