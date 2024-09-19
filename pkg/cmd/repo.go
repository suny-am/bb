/*
// Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Bitbucket repository information",
	Long: `Use this command to get general information about public 
or workspace repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TBD: get system env instead?
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading .env file")
		}

		workspace, _ := cmd.Flags().GetString("workspace")
		repository, _ := cmd.Flags().GetString("repository")

		limit, _ := cmd.Flags().GetString("limit")

		if limit == "" {
			limit = "10"
		}

		credentials, err := CredProvider.GetCredentials()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		authHeaderData := fmt.Sprintf("Basic %s", credentials)

		client := resty.New()

		// TBD add workspace as argument

		endpoint := "https://api.bitbucket.org/2.0/repositories"

		if workspace != "" {
			endpoint = fmt.Sprintf("%s/%s", endpoint, workspace)

			// --repository requires --workspace
			if repository != "" {
				endpoint = fmt.Sprintf("%s/%s", endpoint, repository)
			}
		}

		resp, err := client.R().
			SetHeader("Authorization", authHeaderData).
			SetHeader("Accept", "application/json").
			SetQueryParam("pagelen", limit).
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

			output, err := json.MarshalIndent(data["values"], "", "  ")

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(output))

		}
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)

	repoCmd.Flags().StringP("workspace", "w", "", "Target workspace")
	repoCmd.Flags().StringP("repository", "r", "", "Target repository")
	repoCmd.Flags().StringP("limit", "l", "", "Item limit")
}
