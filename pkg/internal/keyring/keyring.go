package keyring

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v3"
)

type CredentialsProvider interface {
	GetCredentials() (string, error)
}

type OXSKeyChainProvider struct{}
type EnvVarProvider struct{}

func (p *OXSKeyChainProvider) GetCredentials() (string, error) {

	var configFilePath string
	home, empty := os.LookupEnv("HOME")

	if empty {
		pathString, empty := os.LookupEnv("PATH")
		if empty {
			return "", errors.New("no system $PATH found")
		}
		paths := strings.Split(pathString, ":")

		configIndex := slices.IndexFunc(paths, func(p string) bool {
			pathArray := strings.Split(p, "/")
			return pathArray[len(pathArray)-2] == "bitbucket-cli"
		})

		if configIndex == -1 {
			return "", errors.New("config not found in $PATH")
		}

		configFilePath = fmt.Sprintf("%s/hosts.yml", paths[configIndex])
	} else {
		configFilePath = fmt.Sprintf("%s/.config/bitbucket-cli/hosts.yml", home)
	}

	buffer, err := os.ReadFile(configFilePath)

	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = yaml.Unmarshal(buffer, &data)

	if err != nil {
		return "", err
	}

	username := data["user"].(string)
	password, err := keyring.Get("bitbucket-cli", username)

	if err != nil {
		return "", err
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	return credentials, nil
}

func (p *EnvVarProvider) GetCredentials() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	username := os.Getenv("BITBUCKET_USERNAME")
	if username == "" {
		return "", errors.New("BITBUCKET_USERNAME not set")
	}

	password := os.Getenv("BITBUCKET_PASSWORD")
	if password == "" {
		return "", errors.New("BITBUCKET_PASSWORD not set")
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	return credentials, nil
}

func NewCredentialsProvider() CredentialsProvider {
	if os.Getenv("APP_ENV") == "development" {
		return &EnvVarProvider{}
	}
	return &OXSKeyChainProvider{}
}
