package keyring

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	config "github.com/suny-am/bb/internal/config"
	"github.com/zalando/go-keyring"
)

type CredentialsProvider interface {
	GetCredentials() (string, error)
}

type (
	OXSKeyChainProvider struct{}
	EnvVarProvider      struct{}
)

func (p *OXSKeyChainProvider) GetCredentials() (string, error) {
	username, err := config.GetUsername()
	if err != nil {
		return "", err
	}
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
