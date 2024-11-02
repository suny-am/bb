package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/config"
)

func listPullrequests(opts *PrListOptions) (*api.Pullrequests, error) {
	client := &http.Client{}
	var pullrequests api.Pullrequests

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)

	var endpoint string

	if opts.repository == "" {

		user, err := config.GetUsername()
		if err != nil {
			return nil, err
		}
		endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/pullrequests/%s", user)
	} else {

		workspace, err := config.GetWorkspace()
		if err != nil {
			if opts.workspace == "" {
				return nil, err
			}
			workspace = opts.workspace
		}

		endpoint = "https://api.bitbucket.org/2.0/repositories"
		endpoint = fmt.Sprintf("%s/%s/%s/pullrequests", endpoint, workspace, opts.repository)

		if opts.titleFilter != "" {
			endpoint = fmt.Sprintf("%s?q=title~\"%s\"", endpoint, opts.titleFilter)
		} else if opts.authorFilter != "" {
			endpoint = fmt.Sprintf("%s?q=author.nickname=\"%s\"", endpoint, opts.authorFilter)
			endpoint = strings.ReplaceAll(endpoint, " ", "%20")
		}
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

	fetchPullrequestsRecurse(client, req, &pullrequests)

	if err != nil {
		return nil, err
	}

	return &pullrequests, nil
}

func fetchPullrequestsRecurse(client *http.Client, req *http.Request, pullrequests *api.Pullrequests) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var particalPullrequests api.Pullrequests

	if err := json.Unmarshal([]byte(body), &particalPullrequests); err != nil {
		fmt.Println(err)
		return
	}

	if particalPullrequests.Values != nil {
		pullrequests.Values = append(pullrequests.Values, particalPullrequests.Values...)
		if particalPullrequests.Next != "" {
			newReq, err := http.NewRequest("GET", particalPullrequests.Next, nil)
			newReq.Header.Add("Authorization", req.Header["Authorization"][0])
			newReq.Header.Add("Accept", req.Header["Accept"][0])
			if err != nil || len(pullrequests.Values) >= opts.limit {
				fmt.Println(err)
				return
			}
			fetchPullrequestsRecurse(client, newReq, pullrequests)
		}
	}
}
