package cmd

import (
	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(outCmd)
}

var outCmd = &cobra.Command{
  Use: "out",
  Short: "output todos as yaml",
  Run: func(cmd *cobra.Command, args []string) {
    taskmap = utils.LoadTasks(path)
    if todo, ok := taskmap[key]; ok {
      core.PrintYaml(todo)
    } else {
      core.PrintYaml(taskmap)
    }
  },
}
