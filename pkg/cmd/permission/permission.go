/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package permission

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/suny-am/bitbucket-cli/pkg/cmd"
)

type (
	PermissionResponse struct {
		Values []Permission
	}

	Permission struct {
		Permission string
		Type       string
		User       User
		Repository Repository
	}

	User struct {
		AccountId   string
		DisplayNamy string
		Nickname    string
		Type        string
		Uuid        string
	}

	Repository struct {
		Full_Name string
		Name      string
		DataType  string
		Uuid      string
		Type      string
	}
)

var Credentials = cmd.Credentials

// permissionCmd represents the permission command
var permissionCmd = &cobra.Command{
	Use:   "permission",
	Short: "User permissions",
	Long:  `Retrieve permissions for a specific user`,
	Run: func(cmd *cobra.Command, args []string) {
		authHeaderData := fmt.Sprintf("Basic %s", Credentials)

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
			var data PermissionResponse

			if err := json.Unmarshal([]byte(resp.String()), &data); err != nil {
				fmt.Println(err)
			}

			for i := range data.Values {
				fmt.Printf("Repository: %s\n", data.Values[i].Repository.Full_Name)
				fmt.Printf("Permission: %s\n\n", data.Values[i].Permission)
			}
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(permissionCmd)

	permissionCmd.Flags().StringP("user", "u", "", "Target user")
}
