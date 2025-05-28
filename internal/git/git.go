package git

import (
	"log"
	"os/exec"
)

func GetGitRepo() string {
	cmd, err := exec.Command("zsh", "-c", "git rev-parse --show-toplevel").Output()
	if err != nil {
		log.Fatal("Not a git repository")
	}

	return string(cmd)
}
