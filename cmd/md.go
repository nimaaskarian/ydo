package cmd

import (
	"time"

	"github.com/nimaaskarian/ydo/core"
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
    due_time, err := utils.ParseDate(due, time.Now())
    if err != nil {
      return err
    }
    md_config := config.Markdown
    md_config.Limit = 0
    md_config.Filter = func(task core.Task, taskmap core.TaskMap) bool {
      return due_time.IsZero() || (task.Due.Sub(due_time).Abs() < time.Hour*24)
    }
    if len(keys) == 0 {
      taskmap.PrintMarkdown(&md_config)
    } else {
        seen_keys := make(map[string]bool, len(keys))
        for _, key := range keys {
          task, err := taskmap.GetTask(key)
          if err != nil {
            return err
          }
          task.PrintMarkdown(taskmap, 0, seen_keys, key, nil, &md_config)
        }
    }
    return nil
  },
}
