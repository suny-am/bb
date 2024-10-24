package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

func GetUsername() (string, error) {
	var configFilePath string
	home, homeFound := os.LookupEnv("HOME")

	if !homeFound {
		pathString, pathFound := os.LookupEnv("PATH")
		if !pathFound {
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

		configFilePath = fmt.Sprintf("%s/config.yml", paths[configIndex])
	} else {
		configFilePath = fmt.Sprintf("%s/.config/bitbucket-cli/config.yml", home)
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

	return username, nil
}
