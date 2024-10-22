package code

import (
	"github.com/spf13/cobra"
	"github.com/suny-am/bb/pkg/cmd/code/search"
)

var CodeCmd = &cobra.Command{
	Use:   "code",
	Short: "Bitbucket code command",
	Long:  "Search for or edit (TBD) code in a bitbucket repository",
}

func init() {
	CodeCmd.AddCommand(search.SearchCmd)
}
