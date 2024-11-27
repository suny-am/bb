package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func LoadConfig() error {
	var configFilePath string
	home, homeFound := os.LookupEnv("HOME")

	if !homeFound {
		pathString, pathFound := os.LookupEnv("PATH")
		if !pathFound {
			return errors.New("no system $PATH found")
		}
		paths := strings.Split(pathString, ":")

		configIndex := slices.IndexFunc(paths, func(p string) bool {
			pathArray := strings.Split(p, "/")
			return pathArray[len(pathArray)-2] == "bitbucket-cli"
		})

		if configIndex == -1 {
			return errors.New("config not found in $PATH")
		}

		configFilePath = fmt.Sprintf("%s/config.yml", paths[configIndex])
	} else {
		configFilePath = fmt.Sprintf("%s/.config/bitbucket-cli/config.yml", home)
	}

	if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		return err
	}

	return nil
}

func GetUsername() (string, error) {
	if err := LoadConfig(); err != nil {
		return "", err
	}
	username := k.String("user")
	if username == "" {
		return "", errors.New("Could not get username from config")
	}
	return username, nil
}

func GetWorkspace() (string, error) {
	if err := LoadConfig(); err != nil {
		return "", err
	}
	workspace := k.String("workspace")
	if workspace == "" {
		return "", errors.New("Could not get workspace from config")
	}
	return workspace, nil
}

func GetSpinnerStyle() (*int, error) {
	if err := LoadConfig(); err != nil {
		return nil, err
	}
	spinnerStyle := k.Int("spinner")
	if spinnerStyle == 0 {
		return nil, errors.New("Could not get spinner style from config")
	}
	return &spinnerStyle, nil
}
