/*
Copyright © 2024 Calle Sandberg <visualarea.1@gmail.com>
*/
package main

import (
	"github.com/suny-am/bitbucket-cli/pkg/cmd"
	_ "github.com/suny-am/bitbucket-cli/pkg/cmd/commit"
	_ "github.com/suny-am/bitbucket-cli/pkg/cmd/permission"
	_ "github.com/suny-am/bitbucket-cli/pkg/cmd/pr"
	_ "github.com/suny-am/bitbucket-cli/pkg/cmd/repo"
	_ "github.com/suny-am/bitbucket-cli/pkg/cmd/user"
)

func main() {
	cmd.Execute()
}
