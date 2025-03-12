package cmd

import (
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(undoCmd)
}

var undoCmd = &cobra.Command{
  Use: "undo",
  Short: "set a task as not completed",
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) > 0 {
      key = args[0];
    }
    if key == ""  {
      if !utils.ReadYesNo("This will set all tasks as not completed. ARE YOU REALLY SURE? (yes/no) ")  {
        return
      }
      for key := range taskmap {
        taskmap.Undo(key)
      }
    } else {
      taskmap.Undo(key)
      utils.MustHaveTask(taskmap, key)
    }
    utils.WriteTaskmap(taskmap, path)
  },
}
