package http2

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/suny-am/bb/internal/spinner"
)

type Client struct {
	Instance *http.Client
	Debug    bool
}

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
