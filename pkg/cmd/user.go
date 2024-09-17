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

var UserID string
var Email string

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Bitbucket user information",
	Long: `Use this command to get general information about one or more
Bitbucket users.`,
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

		endpoint := "https://api.bitbucket.org/2.0/user"

		if UserID != "" && Email == "" {
			endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/users/%{s}", UserID)
		}

		if UserID == "" && Email != "" {
			endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/user/emails/%s", Email)
		}

		// Email overrides User for now

		if UserID != "" && Email != "" {
			endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/user/emails/%s", Email)
		}

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
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	userCmd.Flags().StringVarP(&UserID, "user", "u", "", "User Bitbucket Account ID/UUID")
	userCmd.Flags().StringVarP(&Email, "email", "e", "", "Email registered on Bitbucket")
}
