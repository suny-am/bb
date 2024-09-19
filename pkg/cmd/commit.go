/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
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

		credentials, err := CredProvider.GetCredentials()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

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
			var data map[string]interface{}

			if err := json.Unmarshal([]byte(resp.String()), &data); err != nil {
				fmt.Println(err)
			}

			output, err := json.MarshalIndent(data, "", "  ")

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(output))
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringP("workspace", "w", "", "Target workspace")
	commitCmd.Flags().StringP("repository", "r", "", "Repository for the commit")
	commitCmd.Flags().StringP("commit", "c", "", "Target commit")

	commitCmd.MarkFlagRequired("workspace")
	commitCmd.MarkFlagRequired("repository")
	commitCmd.MarkFlagRequired("commit")
}
