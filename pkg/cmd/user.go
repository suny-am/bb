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

		userId, _ := cmd.Flags().GetString("userId")
		email, _ := cmd.Flags().GetString("email")

		username := os.Getenv("BITBUCKET_USERNAME")
		appPassword := os.Getenv("BITBUCKET_APP_PASSWORD")
		credentials := fmt.Sprintf("%s:%s", username, appPassword)
		b64 := base64.StdEncoding.EncodeToString([]byte(credentials))
		authHeaderData := fmt.Sprintf("Basic %s", b64)

		client := resty.New()

		// TBD add workspace as argument

		endpoint := "https://api.bitbucket.org/2.0/user"

		if userId != "" && email == "" {
			endpoint = fmt.Sprintf(`https://api.bitbucket.org/2.0/users/{%s}`, userId)
		}

		if userId == "" && email != "" {
			endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/user/emails/%s", email)
		}

		// email overrides User for now

		if userId != "" && email != "" {
			endpoint = fmt.Sprintf("https://api.bitbucket.org/2.0/user/emails/%s", email)
		}

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

	userCmd.Flags().StringP("user", "u", "", "Target user Bitbucket account ID/UUID")
	userCmd.Flags().StringP("email", "e", "", "Targer user email registered on Bitbucket")
}
