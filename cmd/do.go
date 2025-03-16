package cmd

import (
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(doCmd)
  doCmd.ValidArgsFunction = TaskKeyCompletionNotDone
}

var doCmd = &cobra.Command{
  Use: "do [tasks]",
  Short: "set tasks as completed",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    if len(keys) > 0 {
      for _,key := range keys {
        taskmap.Do(key)
      }
    } else {
      if !utils.ReadYesNo("This will set all tasks as completed. ARE YOU REALLY SURE? (yes/no) ")  {
        return
      }
      for key := range taskmap {
        taskmap.Do(key)
      }
    }
  },
}
