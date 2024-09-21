/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package user

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/pkg/types"
)

type (
	UserResponse struct {
		Type         string
		Created_On   string
		Display_Name string
		Uuid         string
	}
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Bitbucket user information",
	Long: `Use this command to get general information about one or more
Bitbucket users.`,
	Run: func(cmd *cobra.Command, args []string) {

		userId, _ := cmd.Flags().GetString("userId")
		email, _ := cmd.Flags().GetString("email")

		cmd.Root().PreRun(cmd, nil)

		credentials := cmd.Context().Value(types.CredentialsKey{})

		authHeaderData := fmt.Sprintf("Basic %s", credentials)

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
			var response UserResponse

			if err := json.Unmarshal([]byte(resp.String()), &response); err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Display name: %s\nType: %s\nUUID: %s\nCreation time: %s",
				response.Display_Name,
				response.Type,
				response.Uuid,
				response.Created_On)
		}
	},
}

func init() {
	UserCmd.Flags().StringP("user", "u", "", "Target user Bitbucket account ID/UUID")
	UserCmd.Flags().StringP("email", "e", "", "Targer user email registered on Bitbucket")
}
