package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(todoCmd)
  todoCmd.Flags().StringVarP(&due, "due", "u", "", "specify due for the tasks to print")
  todoCmd.RegisterFlagCompletionFunc("due", DueCompletion)
  todoCmd.ValidArgsFunction = TaskKeyCompletionFilter(core.Task.IsNotDone)
}

var todoCmd = &cobra.Command{
  Aliases: []string{"t"},
  Use: "todo [tasks (optional)]",
  Short: "output to-do (unfinished) tasks as markdown",
  Run: func(cmd *cobra.Command, keys []string) {
    var due_time time.Time
    if due != "" {
      due_time = utils.ParseDate(due)
    }
    if len(keys) == 0 {
      if due_time.IsZero() {
        taskmap.PrintMarkdown(core.Task.IsNotDone, config.Markdown)
      } else {
        for key, task := range taskmap {
          if task.Due == due_time {
            task.PrintMarkdown(taskmap, 0, map[string]bool{}, key, core.Task.IsNotDone, config.Markdown.Indent)
          }
        }
      }
    } else {
        seen_keys := make(map[string]bool, len(keys))
        for _, key := range keys {
          task, ok := taskmap[key];
          if !ok {
            slog.Error("No such task", "key", key)
            os.Exit(1)
          }
          if due_time.IsZero() != (task.Due == due_time) {
            task.PrintMarkdown(taskmap, 0, seen_keys, key, core.Task.IsNotDone, config.Markdown.Indent)
          }
        }
    }
  },
}
