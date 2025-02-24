package util

import (
	"os"
	"strings"
)

func GetCurrentDir() string {
	path, _ := os.Getwd()
	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
}
