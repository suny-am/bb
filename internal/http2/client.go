package http2

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/api"
	"github.com/suny-am/bb/internal/spinner"
)

type EndpointOptions interface {
	*api.PullrequestViewOptions |
		*api.PullrequestListOptions |
		*api.PipelineListOptions |
		*api.CodeSearchOptions |
		*api.ForkListptions |
		*api.RepositoryViewOptions |
		*api.RepositoryListOptions
}

type Client struct {
	Instance *http.Client
	Debug    bool
}

const apiBaseUrl = "https://api.bitbucket.org/2.0"

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.Debug {
		debugMsg := fmt.Sprintf("\n[DEBUG]\nRequest URL:\n%s", req.URL)
		spinner.AddToView(debugMsg)
	}

	resp, err := c.Instance.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		debugMsg := fmt.Sprintf("\n[DEBUG]\nResponse:\n%v", resp)
		spinner.AddToView(debugMsg)
	}

	return resp, nil
}

func Init(cmd *cobra.Command) *Client {
	debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
	client := &Client{&http.Client{}, debug}
	return client
}

func DetermineQueryParametersDirect(endpoint string, queryParams []string) string {
	queryParamEntries := 0

	for _, qp := range queryParams {

		fmt.Println("QueryParameters", qp, "\r")

		if qp == "" {
			continue
		}

		var prefix string

		if queryParamEntries == 0 {
			prefix = "?"
		} else {
			prefix = "&"
		}

		queryParamEntries++

		endpoint = fmt.Sprintf("%s%s%s", endpoint, prefix, qp)
	}

	endpoint = strings.ReplaceAll(endpoint, " ", "%20")

	return endpoint
}

func DetermineQueryParameters[E EndpointOptions](opts E, endpoint string) string {
	optEntries := reflect.Indirect(reflect.ValueOf(opts))

	queryParams := optEntries.FieldByName("QueryParameters")

	queryParamsLength := queryParams.Type().NumField()

	queryParamEntries := 0

	for i := range queryParamsLength {

		queryParamKey := queryParams.Type().Field(i).Name
		queryParamValue := queryParams.Field(i)

		if !queryParamValue.IsZero() {

			if queryParamKey == "PageLen" && queryParamValue.Int() > 100 {
				queryParamValue.SetInt(100)
			}

			var prefix string

			if queryParamEntries == 0 {
				prefix = "?"
			} else {
				prefix = "&"
			}

			queryParamEntries++

			endpoint = fmt.Sprintf("%s%s%s=%v", endpoint, prefix, strings.ToLower(queryParamKey), queryParamValue)
		}
	}
	endpoint = strings.ReplaceAll(endpoint, " ", "%20")

	return endpoint
}

func DetermineRepositoryEndpoint[E EndpointOptions](opts E) string {
	optEntries := reflect.Indirect(reflect.ValueOf(opts))

	workspace := optEntries.FieldByName("Workspace")
	repository := optEntries.FieldByName("Repository")

	return fmt.Sprintf("%s/repositories/%s/%s",
		apiBaseUrl,
		workspace,
		repository)
}

func DetermineWorkspaceEndpoint[E EndpointOptions](opts E) string {
	optEntries := reflect.Indirect(reflect.ValueOf(opts))

	workspace := optEntries.FieldByName("Workspace")

	return fmt.Sprintf("%s/workspaces/%s",
		apiBaseUrl,
		workspace)
}
