package cmd

import (
	"log/slog"
	"os"
	"time"

	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(mdCmd)
  mdCmd.Flags().StringVarP(&due, "due", "u", "", "specify due for the tasks to print")
  mdCmd.RegisterFlagCompletionFunc("due", DueCompletion)
  mdCmd.ValidArgsFunction = TaskKeyCompletion
}

var mdCmd = &cobra.Command{
  Use: "md [tasks (optional)]",
  Short: "output tasks as markdown (run with no args so it'd output all tasks like `ydo` does)",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    var due_time time.Time
    if due != "" {
      due_time = utils.ParseDate(due)
    }
    if len(keys) == 0 {
      if due_time.IsZero() {
        taskmap.PrintMarkdown(nil, config.Markdown)
      } else {
        for key, task := range taskmap {
          if task.Due == due_time {
            task.PrintMarkdown(taskmap, 0, map[string]bool{}, key, nil, config.Markdown.Indent)
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
            task.PrintMarkdown(taskmap, 0, seen_keys, key, nil, config.Markdown.Indent)
          }
        }
    }
  },
}
