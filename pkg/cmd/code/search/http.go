package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/suny-am/bitbucket-cli/api"
)

func searchCode(opts *SearchOptions) (*api.CodeSearchResponse, error) {
	client := http.Client{}
	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)
	endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/workspaces/%s/search/code?search_query=%s", opts.workspace, strings.ReplaceAll(opts.searchParam, " ", "%20"))
	req, err := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("Authorization", authHeaderValue)
	req.Header.Add("Accept", "application/json")

	if opts.limit < 0 {
		opts.limit = 0
	}

	if opts.limit > 0 {
		q := req.URL.Query()
		q.Add("pagelen", strconv.Itoa(opts.limit))
		req.URL.RawQuery = q.Encode()
	}

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var code api.CodeSearchResponse

	if err = json.Unmarshal([]byte(body), &code); err != nil {
		return nil, err
	}

	for i := range code.Values {
		c := &code.Values[i]
		repo := strings.Split(c.File.Links.Self.Href, "/")[6]
		repoEndpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", opts.workspace, repo)
		defaultBranch := fetchDefaultBranch(&client, repoEndpoint, authHeaderValue)
		if defaultBranch != "" {
			c.File.Links.Html.Href = fmt.Sprintf("https://bitbucket.org/%s/%s/src/%s/%s", opts.workspace, repo, defaultBranch, c.File.Path)
		}
		if opts.includeSource {
			src := fetchSource(&client, c.File.Links.Self.Href, c.File.Path, authHeaderValue)
			if src == "" {
				continue
			}
			c.File.Source = src
		}
	}

	return &code, nil
}

func fetchSource(client *http.Client, srcPath string, extension string, authHeaderValue string) string {
	ext := strings.Replace(filepath.Ext(extension), ".", "", -1)
	srcReq, err := http.NewRequest("GET", srcPath, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	srcReq.Header.Add("Authorization", authHeaderValue)
	resp, err := client.Do(srcReq)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	src := fmt.Sprintf("```%s\n%s\n```", ext, body)
	return src
}

func fetchDefaultBranch(client *http.Client, endpoint string, authHeaderValue string) string {
	defaultBranchReq, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defaultBranchReq.Header.Add("Authorization", authHeaderValue)
	resp, err := client.Do(defaultBranchReq)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var repository api.Repository

	if err := json.Unmarshal([]byte(body), &repository); err != nil {
		return ""
	}

	return repository.Mainbranch.Name
}
