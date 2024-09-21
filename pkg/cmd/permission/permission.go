/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package permission

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/api"
	"github.com/suny-am/bitbucket-cli/pkg/types"
)

// PermissionCmd represents the permission command
var PermissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "User permissions",
	Long:  `Retrieve permissions for a specific user`,
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Root().PreRun(cmd, nil)
		credentials := cmd.Context().Value(types.CredentialsKey{})

		authHeaderData := fmt.Sprintf("Basic %s", credentials)

		client := resty.New()

		endpoint := "https://api.bitbucket.org/2.0/user/permissions/repositories"

		resp, err := client.R().
			SetHeader("Authorization", authHeaderData).
			SetHeader("Accept", "application/json").
			EnableTrace().
			Get(endpoint)

		if resp.IsError() {
			fmt.Println(err.Error())
		}

		if resp.IsSuccess() {
			var response api.Permissions

			if err := json.Unmarshal([]byte(resp.String()), &response); err != nil {
				fmt.Println(err)
			}

			for i := range response.Values {
				fmt.Printf("Repository: %s\nPermission: %s\n\n",
					response.Values[i].Repository.Full_Name,
					response.Values[i].Permission)
			}
		}
	},
}

func init() {
	PermissionCmd.Flags().StringP("user", "u", "", "Target user")
}
