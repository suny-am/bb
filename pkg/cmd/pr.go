/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Pull Request information",
	Long: `Get information for a pull request,
such as status, commit tree and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading .env file")
		}

		workspace, _ := cmd.Flags().GetString("workspace")
		repository, _ := cmd.Flags().GetString("repository")
		commit, _ := cmd.Flags().GetString("commit")

		endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests", workspace, repository)

		if commit != "" {
			endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commit/%s/pullrequests", workspace, repository, commit)
		}

		client := resty.New()

		username := os.Getenv("BITBUCKET_USERNAME")
		appPassword := os.Getenv("BITBUCKET_APP_PASSWORD")
		credentials := fmt.Sprintf("%s:%s", username, appPassword)
		b64 := base64.StdEncoding.EncodeToString([]byte(credentials))
		authHeaderData := fmt.Sprintf("Basic %s", b64)

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

			output, err := json.MarshalIndent(data["values"], "", "  ")

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(output))
		}
	},
}

func init() {
	rootCmd.AddCommand(prCmd)

	prCmd.Flags().StringP("workspace", "w", "", "Target workspace")
	prCmd.Flags().StringP("repository", "r", "", "Target repository")
	prCmd.Flags().StringP("commit", "c", "", "commit for the target PR(s)")

	prCmd.MarkFlagRequired("workspace")
	prCmd.MarkFlagRequired("repository")
}