package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nimaaskarian/ydo/core"
)

func init() {
  rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
  Use: "list",
  Short: "list all todos of selected files",
  Run: func(cmd *cobra.Command, args []string) {
    if todo, ok := todomap[key]; ok {
      core.PrintYaml(todo)
    } else {
      core.PrintYaml(todomap)
    }
  },
}
