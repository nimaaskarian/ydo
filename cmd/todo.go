package cmd

import (
	"log/slog"
	"os"

	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(todoCmd)
  todoCmd.ValidArgsFunction = TaskKeyCompletionNotDone
}

var todoCmd = &cobra.Command{
  Aliases: []string{"t"},
  Use: "todo [tasks (optional)]",
  Short: "output to-do (unfinished) tasks as markdown",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    if len(keys) == 0 {
      taskmap.PrintMarkdown(core.Task.IsNotDone)
    } else {
        for _, key := range keys {
          task, ok := taskmap[key];
          if !ok {
            slog.Error("No such task", "key", key)
            os.Exit(1)
          }
          task.PrintMarkdown(taskmap, 0, &[]string{}, key, core.Task.IsNotDone)
        }
    }
  },
}
