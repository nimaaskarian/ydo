package cmd

import (
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(rmCmd)
  rmCmd.ValidArgsFunction = TaskKeyCompletionNotDone
}

var rmCmd = &cobra.Command{
  Aliases: []string{"del", "delete"},
  Use: "rm [tasks]",
  Short: "remove tasks completely from the data (also removes it from all the dependency lists)",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    if len(keys) > 0 {
      for _,key := range keys {
        delete(taskmap, key)
        taskmap.WipeDependenciesToKey(key)
      }
    } else {
      if !utils.ReadYesNo("WARN This will DELETE ALL THE TASKS. ARE YOU REALLY SURE? (yes/no) ")  {
        return
      }
      for key := range taskmap {
        delete(taskmap, key)
        taskmap.WipeDependenciesToKey(key)
      }
    }
    taskmap.Write(tasks_path)
  },
}
