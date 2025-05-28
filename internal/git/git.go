package git

import (
	"log"
	"os/exec"
	"strings"
)

func GetGitRepo() string {
	cmd, err := exec.Command("zsh", "-c", "git rev-parse --show-toplevel").Output()
	if err != nil {
		log.Fatal("Not a git repository")
	}

	slice := strings.Split(string(cmd), "/")

	return strings.TrimSuffix(slice[len(slice)-1], "\n")
}
