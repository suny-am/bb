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

var Commit string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Bitbucket commit information",
	Long: `Use this command to get commit activity information
	from either public or workspace repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading .env file")
		}

		username := os.Getenv("BITBUCKET_USERNAME")
		appPassword := os.Getenv("BITBUCKET_APP_PASSWORD")
		credentials := fmt.Sprintf("%s:%s", username, appPassword)
		b64 := base64.StdEncoding.EncodeToString([]byte(credentials))
		authHeaderData := fmt.Sprintf("Basic %s", b64)

		client := resty.New()

		// TBD add workspace as argument

		if Workspace == "" {
			fmt.Println("--workspace required")
			return
		}

		if Repository == "" {
			fmt.Println("--repository required")
			return
		}

		if Commit == "" {
			fmt.Println("--commit required")
			return
		}

		endpoint := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/commit/%s", Workspace, Repository, Commit)

		fmt.Println(endpoint)

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitCmd.Flags().StringVarP(&Workspace, "workspace", "w", "", "Workspace for the repository")
	commitCmd.Flags().StringVarP(&Repository, "repository", "r", "", "Repository for the commit")
	commitCmd.Flags().StringVarP(&Commit, "commit", "c", "", "Target commit")
}
