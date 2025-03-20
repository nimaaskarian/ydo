package cmd

import (
	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(todoCmd)
  todoCmd.Flags().StringVarP(&due, "due", "u", "", "specify due for the tasks to print")
  todoCmd.RegisterFlagCompletionFunc("due", DueCompletion)
}

var todoCmd = &cobra.Command{
  Aliases: []string{"t"},
  Use: "todo [tasks (optional)]",
  Short: "output to-do (unfinished) tasks as markdown",
  ValidArgsFunction: TaskKeyCompletionFilter(core.Task.IsNotDone),
  RunE: func(cmd *cobra.Command, keys []string) error {
    due_time, err := utils.ParseDate(due)
    if err != nil {
      return err
    }
    if len(keys) == 0 {
      if due_time.IsZero() {
        taskmap.PrintMarkdown(core.Task.IsNotDone, &config.Markdown)
      } else {
        for key, task := range taskmap {
          if task.Due == due_time {
            task.PrintMarkdown(taskmap, 0, map[string]bool{}, key, core.Task.IsNotDone, &config.Markdown)
          }
        }
      }
    } else {
      seen_keys := make(map[string]bool, len(keys))
      for _, key := range keys {
        task, err := taskmap.GetTask(key)
        if err != nil {
          return err
        }
        if due_time.IsZero() != (task.Due == due_time) {
          task.PrintMarkdown(taskmap, 0, seen_keys, key, core.Task.IsNotDone, &config.Markdown)
        }
      }
    }
    return nil
  },
}
