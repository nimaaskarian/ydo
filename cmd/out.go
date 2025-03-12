package cmd

import (
	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(outCmd)
}

var outCmd = &cobra.Command{
  Use: "out",
  Short: "output tasks as yaml",
  Run: func(cmd *cobra.Command, args []string) {
    if task, ok := taskmap[key]; ok {
      core.PrintYaml(task)
    } else {
      core.PrintYaml(taskmap)
    }
  },
}
