/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package pr

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/api"
	"github.com/suny-am/bitbucket-cli/pkg/types"
)

var PrCmd = &cobra.Command{
	Use:   "pr",
	Short: "Pull Request information",
	Long: `Get information for a pull request,
such as status, commit tree and more.`,
	Run: func(cmd *cobra.Command, args []string) {

		workspace, _ := cmd.Flags().GetString("workspace")
		repository, _ := cmd.Flags().GetString("repository")
		commit, _ := cmd.Flags().GetString("commit")

		endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests", workspace, repository)

		if commit != "" {
			endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commit/%s/pullrequests", workspace, repository, commit)
		}

		client := resty.New()

		cmd.Root().PreRun(cmd, nil)
		credentials := cmd.Context().Value(types.CredentialsKey{})

		authHeaderData := fmt.Sprintf("Basic %s", credentials)

		resp, err := client.R().
			SetHeader("Authorization", authHeaderData).
			SetHeader("Accept", "application/json").
			EnableTrace().
			Get(endpoint)

		if resp.IsError() {
			fmt.Println(err.Error())
		}

		if resp.IsSuccess() {
			var response api.PullRequests

			if err := json.Unmarshal([]byte(resp.String()), &response); err != nil {
				fmt.Println(err)
			}

			for i := range response.Values {
				fmt.Printf("Title: %s\n", response.Values[i].Title)
				fmt.Printf("Description: %s\n", response.Values[i].Description)
				fmt.Printf("State: %s\n", response.Values[i].State)
				fmt.Printf("Type: %s\n", response.Values[i].Type)
				fmt.Printf("Source: %s\n", response.Values[i].Source.Commit.Hash)
				fmt.Printf("Destination: %s\n", response.Values[i].Destination.Commit.Hash)
				fmt.Printf("Comment_Count: %d\n", response.Values[i].Comment_Count)
				fmt.Printf("Link: %s\n", response.Values[i].Links.Self["href"])
				fmt.Printf("Updated_On: %s\n\n", response.Values[i].Updated_On)
			}
		}
	},
}

func init() {
	PrCmd.Flags().StringP("workspace", "w", "", "Target workspace")
	PrCmd.Flags().StringP("repository", "r", "", "Target repository")
	PrCmd.Flags().StringP("commit", "c", "", "commit for the target PR(s)")

	PrCmd.MarkFlagRequired("workspace")
	PrCmd.MarkFlagRequired("repository")
}
