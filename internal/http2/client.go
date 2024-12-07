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
		debugMsg := fmt.Sprintf("\n [DEBUG]\n Request URL:\n %s", req.URL)
		spinner.AddToView(debugMsg)
	}
	return c.Instance.Do(req)
}

func Init(cmd *cobra.Command) *Client {
	debug, _ := cmd.Root().PersistentFlags().GetBool("debug")
	client := &Client{&http.Client{}, debug}
	return client
}
