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
  mdCmd.ValidArgsFunction = TaskKeyCompletionFilter(nil)
}

var mdCmd = &cobra.Command{
  Use: "md [tasks (optional)]",
  Short: "output tasks as markdown (run with no args so it'd output all tasks like `ydo` does)",
  RunE: func(cmd *cobra.Command, keys []string) error {
    var due_time time.Time
    due_time, err := utils.ParseDate(due)
    if err != nil {
      return err
    }
    if len(keys) == 0 {
      if due_time.IsZero() {
        taskmap.PrintMarkdown(nil, &config.Markdown)
      } else {
        for key, task := range taskmap {
          if task.Due == due_time {
            task.PrintMarkdown(taskmap, 0, map[string]bool{}, key, nil, &config.Markdown)
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
            task.PrintMarkdown(taskmap, 0, seen_keys, key, nil, &config.Markdown)
          }
        }
    }
    return nil
  },
}
