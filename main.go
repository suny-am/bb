package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func main() {
	envFile, _ := godotenv.Read(".env")
	client := resty.New()

	username := envFile["BITBUCKET_USERNAME"]
	appPassword := envFile["BITBUCKET_APP_PASSWORD"]
	credentials := fmt.Sprintf("%s:%s", username, appPassword)
	credentialsBase64 := base64.StdEncoding.EncodeToString([]byte(credentials))
	credentialsHeaderData := fmt.Sprintf("Basic %s", credentialsBase64)

	resp, err := client.R().
		SetHeader("Authorization", credentialsHeaderData).
		SetHeader("Accept", "application/json").
		Post("https://api.bitbucket.org/2.0/repositories/avupublicapis/integration-ecommerce-authorizer/commits")

	if resp.IsSuccess() {
		var dat map[string]interface{}
		json.Unmarshal([]byte(resp.String()), &dat)

		values := dat["values"].([]interface{})
		for i := range values {
			v := values[i].(map[string]interface{})
			fmt.Println(v["type"])
			fmt.Println(v["hash"])
			fmt.Println(v["author"])
			fmt.Println(v["message"])
			fmt.Println(v["summary"])
		}

	}

	if resp.IsError() {
		fmt.Println(err)
	}

}

// cURL example
// curl \
// -H "Authorization: Basic $my_credentials_after_base64_encoding" \
// -H "Accept: application/json" \
// 	https://api.bitbucket.org/2.0/repositories/avupublicapis/integration-ecommerce-authorizer/commits
