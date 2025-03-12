package cmd

import (
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(doCmd)
  doCmd.ValidArgsFunction = TaskKeyCompletion
}

var doCmd = &cobra.Command{
  Use: "do",
  Short: "set a task as done",
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) > 0 {
      key = args[0];
    }
    if key == ""  {
      if !utils.ReadYesNo("This will set all tasks as completed. ARE YOU REALLY SURE? (yes/no) ")  {
        return
      }
      for key := range taskmap {
        taskmap.Do(key)
      }
    } else {
      utils.MustHaveTask(taskmap, key)
      taskmap.Do(key)
    }
    utils.WriteTaskmap(taskmap, path)
  },
}
