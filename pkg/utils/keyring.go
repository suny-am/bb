package keyring

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	execPathKeychain = "/usr/bin/security"

	// encodingPrefix is a well-known prefix added to strings encoded by Set.
	encodingPrefix       = "go-keyring-encoded:"
	base64EncodingPrefix = "go-keyring-base64:"
)

// type macOSXKeyChain struct{}

// func (k macOSXKeyChain) Get(service, username string) (string, error) {
func GetPassword(service, username string) (string, error) {
	out, err := exec.Command(
		execPathKeychain,
		"find-generic-password",
		"-s", "bitbucket-cli",
		"-wa", username).CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "could not be found") {
			err = exec.ErrNotFound
		}
		return "", err
	}

	trimStr := strings.TrimSpace(string(out[:]))
	// if the string has the well-known prefix, assume it's encoded
	if strings.HasPrefix(trimStr, encodingPrefix) {
		dec, err := hex.DecodeString(trimStr[len(encodingPrefix):])
		return string(dec), err
	} else if strings.HasPrefix(trimStr, base64EncodingPrefix) {
		dec, err := base64.StdEncoding.DecodeString(trimStr[len(base64EncodingPrefix):])
		return string(dec), err
	}

	return trimStr, nil
}

func GetUsername() (string, error) {

	if os.Getenv("APP_ENV") == "development" {
		return os.Getenv("BITBUCKET_USERNAME"), nil
	}

	home := os.Getenv("HOME")

	configFilePath := fmt.Sprintf("%s/.config/bitbucket-cli/hosts.yml", home)
	buffer, err := os.ReadFile(configFilePath)

	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = yaml.Unmarshal(buffer, &data)

	if err != nil {
		return "", err
	}

	return data["user"].(string), nil
}

func Credentials() (string, error) {

	username, err := GetUsername()

	if err != nil {
		return "", err
	}

	var password string

	if os.Getenv("APP_ENV") == "development" {
		password = os.Getenv("BITBUCKET_PASSWORD")
	}

	if os.Getenv("APP_ENV") != "development" {
		password, err = GetPassword("bitbucket-cli", username)
		if err != nil {
			return "", err
		}
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))

	return credentials, nil
}
