package cmd

import (
	"fmt"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(doCmd)
}

var doCmd = &cobra.Command{
  Use: "do",
  Short: "set a task as done",
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) > 0 {
      key = args[0];
    }
    utils.MustHaveTask(taskmap, key)
    taskmap.Do(key)
    utils.WriteTaskmap(taskmap, path)
    fmt.Printf("Completed task %q (%q)\n", key, taskmap[key].Task)
  },
}
