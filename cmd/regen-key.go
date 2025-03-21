package cmd

import (
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(regenKeyCmd)
  regenKeyCmd.ValidArgsFunction = TaskKeyCompletionFilter(nil)
}

var regenKeyCmd = &cobra.Command{
  Use: "regen-key [tasks]",
  Short: "regen key with the automatic key generator (respects config file)",
  Run: func(cmd *cobra.Command, keys []string) {
    if len(keys) > 0 {
      for _,key := range keys {
        task := taskmap[key]
        new_key = taskmap.TfidfNextKey(task.Task, config.Tfidf, key)
        taskmap.ReplaceKeyInDeps(key, new_key)
        taskmap[new_key] = task
      }
    } else {
      if always_yes || utils.ReadYesNo("Regen key for all the tasks? (yes/no) ")  {
        for key := range taskmap {
          task := taskmap[key]
          new_key = taskmap.TfidfNextKey(task.Task, config.Tfidf, key)
          taskmap.ReplaceKeyInDeps(key, new_key)
          taskmap[new_key] = task
        }
      }
    }
  },
}
