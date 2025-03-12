package cmd

import (
	"github.com/spf13/cobra"
	"github.com/nimaaskarian/ydo/utils"
)

func init() {
  rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
  Use: "list",
  Short: "list todos as markdown",
  Run: func(cmd *cobra.Command, args []string) {
    todomap = utils.LoadTodos(path)
    if todo, ok := todomap[key]; ok {
      todo.PrintMarkdown(todomap, 1)
    } else {
      todomap.PrintMarkdown()
    }
  },
}
