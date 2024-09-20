/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/pkg/cmd"
)

type (
	PullRequestResponse struct {
		Values []PullRequest
	}
	PullRequest struct {
		Author              Author
		Close_Source_Branch bool
		Comment_Count       int
		Description         string
		Destination         Destination
		Id                  int
		Merge_Commit        Commit
		Links               Links
		Reason              string
		Source              Source
		State               string
		Task_Count          int
		Title               string
		Type                string
		Updated_On          string
	}
	Links struct {
		Self     map[string]string
		Html     map[string]string
		Commits  map[string]string
		Approve  map[string]string
		Diff     map[string]string
		DiffStat map[string]string
		Comments map[string]string
		Activity map[string]string
		Merge    map[string]string
		Decline  map[string]string
	}
	Source struct {
		Branch     Branch
		Commit     Commit
		Repository Repository
	}
	Destination struct {
		Branch     Branch
		Commit     Commit
		Repository Repository
	}
	Commit struct {
		Hash string
		Type string
	}
	Repository struct {
		Full_Name string
		Name      string
		Type      string
		Uuid      string
	}
	Branch struct {
		Name string
	}
	Author struct {
		Account_Id   string
		Display_Name string
		Nickname     string
		Type         string
		Uuid         string
	}
)

var Credentials = cmd.Credentials

var prCmd = &cobra.Command{
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

		authHeaderData := fmt.Sprintf("Basic %s", Credentials)

		resp, err := client.R().
			SetHeader("Authorization", authHeaderData).
			SetHeader("Accept", "application/json").
			EnableTrace().
			Get(endpoint)

		if resp.IsError() {
			fmt.Println(err.Error())
		}

		if resp.IsSuccess() {
			var response PullRequestResponse

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
	cmd.RootCmd.AddCommand(prCmd)

	prCmd.Flags().StringP("workspace", "w", "", "Target workspace")
	prCmd.Flags().StringP("repository", "r", "", "Target repository")
	prCmd.Flags().StringP("commit", "c", "", "commit for the target PR(s)")

	prCmd.MarkFlagRequired("workspace")
	prCmd.MarkFlagRequired("repository")
}
