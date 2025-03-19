package cmd

import (
	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(undoCmd)
  undoCmd.ValidArgsFunction = TaskKeyCompletionFilter(core.Task.IsDone)
}

var undoCmd = &cobra.Command{
  Aliases: []string{"u"},
  Use: "undo [keys]",
  Short: "set tasks as not completed",
  Run: func(cmd *cobra.Command, keys []string) {
    if len(keys) > 0 {
      for _,key := range keys {
        taskmap.Undo(key)
      }
    } else {
      if !utils.ReadYesNo("This will set all tasks as not completed. ARE YOU REALLY SURE? (yes/no) ")  {
        return
      }
      for key := range taskmap {
        taskmap.Undo(key)
      }
    }
  },
}
