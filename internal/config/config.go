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

type ConfigOption struct {
	Key   string
	Value string
}

func LoadConfig() error {
	configFilePath, err := loadConfigFilePath()
	if err != nil {
		return err
	}

	if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
		return err
	}

	return nil
}

func loadConfigFilePath() (string, error) {
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

	return configFilePath, nil
}

func GetUsername() (string, error) {
	if err := LoadConfig(); err != nil {
		return "", err
	}
	username := k.String("user")
	if username == "" {
		return "", errors.New("could not get username from config")
	}
	return username, nil
}

func GetWorkspace() (string, error) {
	if err := LoadConfig(); err != nil {
		return "", err
	}
	workspace := k.String("workspace")
	if workspace == "" {
		return "", errors.New("could not get workspace from config")
	}
	return workspace, nil
}

func GetSpinnerStyle() (*int, error) {
	if err := LoadConfig(); err != nil {
		return nil, err
	}
	spinnerStyle := k.Int("spinner")
	if spinnerStyle == 0 {
		return nil, errors.New("could not get spinner style from config")
	}
	return &spinnerStyle, nil
}

func GetConfiguredItems() ([]string, error) {
	if err := LoadConfig(); err != nil {
		return nil, err
	}

	items := k.All()

	var itemKeys []string
	for k := range items {
		itemKeys = append(itemKeys, k)
	}
	return itemKeys, nil
}

func SetConfigOption(option ConfigOption) (string, error) {
	if err := LoadConfig(); err != nil {
		return "", err
	}

	if option.Key != "spinner" {
		option.Value = strings.Replace(option.Value, ">", "", 1)
		option.Value = strings.ReplaceAll(option.Value, " ", "")
	}

	err := writeToOptionFile(option.Key, option.Value)
	if err != nil {
		return "could not set value", err
	}

	return fmt.Sprintf("Option '%s' configured", option.Key), nil
}

func writeToOptionFile(key string, value string) error {
	configFilePath, err := loadConfigFilePath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	v, err := yaml.Parser().Unmarshal(data)
	if err != nil {
		return err
	}

	v[key] = value

	newData, err := yaml.Parser().Marshal(v)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, newData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigOption(option string) string {
	selectedOption := k.Get(option)
	return selectedOption.(string)
}
