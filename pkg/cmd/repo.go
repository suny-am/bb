/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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

// Flag variables
var Workspace string

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Bitbucket repository information",
	Long: `Use this command to get general information about public or
	workspace repositories.`,
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

		endpoint := "https://api.bitbucket.org/2.0/repositories/avupublicapis" // TBD add workspace as argument

		resp, err := client.R().
			SetHeader("Authorization", authHeaderData).
			SetHeader("Accept", "application/json").
			EnableTrace().
			Get(endpoint)

		if resp.IsError() {
			fmt.Println(err.Error())
		}

		if resp.IsSuccess() {
			var dat map[string]interface{}

			if err := json.Unmarshal([]byte(resp.String()), &dat); err != nil {
				fmt.Println(err)
			}

			repositories := dat["values"].([]interface{})

			for repoIdx := range repositories {
				repo := repositories[repoIdx].(map[string]interface{})

				// links := repo["links"].(map[string]interface{})
				// for linkIdx := range links {
				// 	l, ok := links[linkIdx].(map[string]interface{})
				// 	if ok {
				// 		fmt.Printf("Link: %s\n", l["href"])
				// 	}
				// 	if !ok {
				// 		l, ok := links[linkIdx].([]interface{})
				// 		if ok {
				// 			for i := range l {
				// 				ref := l[i]
				// 				ll, ok := ref.(map[string]interface{})
				// 				if ok {
				// 					fmt.Printf("\tLink: %s\n", ll["href"])
				// 				}
				// 			}
				// 		}
				// 	}
				// }

				fmt.Println(repo["uuid"])
				fmt.Println(repo["name"])
				owner := repo["owner"].(map[string]interface{})["username"]
				fmt.Printf("Owner: %s\n", owner)
				fmt.Printf("Size: %f Bytes\n\n", repo["size"].(float64))
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(repoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	repoCmd.Flags().StringVarP(&Workspace, "workspace", "w", "", "workspace name")

}
