package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/http2"
	"github.com/suny-am/bb/internal/spinner"
	"github.com/suny-am/bb/internal/textinput"
)

func searchCode(opts *SearchOptions, cmd *cobra.Command) (*api.CodeSearchResponse, error) {
	var code api.CodeSearchResponse
	var err error

	go func() {
		err = search(&code, cmd, opts)
		debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
		if debug {
			textinput.ConfirmKey()
		}
		spinner.Stop()
	}()

	spinner.Start("Searching code")

	return &code, err
}

func search(code *api.CodeSearchResponse, cmd *cobra.Command, opts *SearchOptions) error {
	client := http2.Init(cmd)
	var req *http.Request
	var resp *http.Response
	var body []byte
	var err error

	req, err = generateRequest(opts)
	if err != nil {
		return err
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &code); err != nil {
		return err
	}

	for i := range code.Values {
		includeDefaultBranch(&code.Values[i], client, req)
		if opts.includeSource {
			includeSource(&code.Values[i], client, req)
		}
	}
	return nil
}

func generateRequest(opts *SearchOptions) (*http.Request, error) {
	searchQuery := strings.ReplaceAll(opts.searchParam, " ", "%20")
	if opts.repository != "" {
		searchQuery = fmt.Sprintf("%s+repo:%s",
			searchQuery,
			opts.repository)
	}
	endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/workspaces/%s/search/code?search_query=%s",
		opts.workspace,
		searchQuery)

	authHeaderValue := fmt.Sprintf("Basic %s", opts.credentials)

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

	return req, nil
}

func includeDefaultBranch(code *api.CodeItem, client *http2.Client, req *http.Request) {
	repo := strings.Split(code.File.Links.Self.Href, "/")[6]
	req.URL, _ = req.URL.Parse(fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s", opts.workspace, repo))
	defaultBranch := fetchDefaultBranch(client, req)
	if defaultBranch != "" {
		code.File.Links.Html.Href = fmt.Sprintf("https://bitbucket.org/%s/%s/src/%s/%s", opts.workspace, repo, defaultBranch, code.File.Path)
	}
}

func includeSource(code *api.CodeItem, client *http2.Client, req *http.Request) {
	req.URL, _ = req.URL.Parse(code.File.Links.Self.Href)
	src := fetchSource(client, req, code.File.Path)
	if src == "" {
		return
	}
	code.File.Source = src
}

func fetchSource(client *http2.Client, req *http.Request, extension string) string {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	ext := strings.ReplaceAll(filepath.Ext(extension), ".", "")
	src := fmt.Sprintf("```%s\n%s\n```", ext, body)
	return src
}

func fetchDefaultBranch(client *http2.Client, req *http.Request) string {
	resp, err := client.Do(req)
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
