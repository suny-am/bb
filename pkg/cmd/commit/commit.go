/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package commit

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/pkg/cmd"
)

var Credentials = cmd.Credentials

type (
	CommitResponse struct {
		Author  Author
		Date    string
		Hash    string
		Message string
		Parents []ParentCommit
	}
	Author struct {
		Raw  string
		Type string
		User User
	}
	User struct {
		Account_Id   string
		Display_Name string
		Nickname     string
		Type         string
		Uuid         string
	}
	ParentCommit struct {
		Hash string
		Type string
	}
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Bitbucket commit information",
	Long: `Use this command to get commit activity information
	from either public or workspace repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		workspace, _ := cmd.Flags().GetString("workspace")
		repository, _ := cmd.Flags().GetString("repository")
		commit, _ := cmd.Flags().GetString("commit")

		client := resty.New()

		endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commit/%s", workspace, repository, commit)

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
			var response CommitResponse

			if err := json.Unmarshal([]byte(resp.String()), &response); err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Author: %s\n",
				response.Author.User.Display_Name)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("workspace", "w", "", "Target workspace")
	commitCmd.Flags().StringP("repository", "r", "", "Repository for the commit")
	commitCmd.Flags().StringP("commit", "c", "", "Target commit")

	commitCmd.MarkFlagRequired("workspace")
	commitCmd.MarkFlagRequired("repository")
	commitCmd.MarkFlagRequired("commit")
}
