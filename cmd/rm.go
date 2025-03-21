package cmd

import (
	"log/slog"
	"os"
	"slices"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var cascade bool;

func init() {
  rootCmd.AddCommand(rmCmd)
  rmCmd.Flags().BoolVarP(&cascade, "cascade", "C", false, "cascade remove the dependencies of this task that aren't a dependency to other tasks")
  rmCmd.ValidArgsFunction = TaskKeyCompletion
}

func delete_task(taskmap core.TaskMap, key string) {
  if task, ok := taskmap[key]; ok{
    delete(taskmap, key)
    taskmap.WipeDependenciesToKey(key)
    if cascade {
      for _, key := range task.Deps {
        should_remove := true
        for _, task := range taskmap {
          if slices.Contains(task.Deps, key) {
            should_remove = false
            break
          }
        }
        if should_remove {
          delete(taskmap, key)
        }
      }
    }
  } else {
    slog.Error("No such task", "key", key)
    os.Exit(1)
  }
}
var rmCmd = &cobra.Command{
  Aliases: []string{"del", "delete"},
  Use: "rm [tasks]",
  Short: "remove tasks completely from the data (also removes it from all the dependency lists)",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    if len(keys) > 0 {
      for _,key := range keys {
        delete_task(taskmap, key)
      }
    } else {
      if !utils.ReadYesNo("WARN This will DELETE ALL THE TASKS. ARE YOU REALLY SURE? (yes/no) ")  {
        return
      }
      taskmap = make(core.TaskMap)
    }
  },
}
