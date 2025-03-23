package cmd

import (
	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(undoCmd)
}

var undoCmd = &cobra.Command{
  Aliases: []string{"u"},
  Use: "undo [keys]",
  Short: "set tasks as not completed",
  ValidArgsFunction: TaskKeyCompletionFilter(func(t core.Task, tm core.TaskMap) bool {return t.Done && !t.AutoComplete }),
  RunE: func(cmd *cobra.Command, keys []string) error {
    if len(keys) > 0 {
      for _,key := range keys {
        if err := taskmap.Undo(key); err != nil {
          return err
        }
      }
    } else {
      if always_yes || utils.ReadYesNo("This will set all tasks as not completed. ARE YOU REALLY SURE? (yes/no) ")  {
        for key := range taskmap {
          taskmap.Undo(key)
        }
      }
    }
    return nil
  },
  PostRunE: SaveChanges,
  PreRun: UpdateOldTaskMap,
}
