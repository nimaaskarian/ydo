package cmd

import (
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(regenKeyCmd)
  regenKeyCmd.ValidArgsFunction = TaskKeyCompletion
}

var regenKeyCmd = &cobra.Command{
  Use: "regen-key [tasks]",
  Short: "regen key with the automatic key generator (respects config file)",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    if len(keys) > 0 {
      for _,key := range keys {
        task := taskmap[key]
        new_key = taskmap.TfidfNextKey(task.Task, config.Tfidf, key)
        taskmap.ReplaceKeyInDeps(key, new_key)
        taskmap[new_key] = task
      }
    } else {
      if !utils.ReadYesNo("Regen key for all the tasks? (yes/no) ")  {
        return
      }
      for key := range taskmap {
        task := taskmap[key]
        new_key = taskmap.TfidfNextKey(task.Task, config.Tfidf, key)
        taskmap.ReplaceKeyInDeps(key, new_key)
        taskmap[new_key] = task
      }
    }
  },
}
