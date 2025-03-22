package cmd

import (
	"time"

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
    due_time, err := utils.ParseDate(due, time.Now())
    if err != nil {
      return err
    }
    todo_config := config.Markdown
    todo_config.Limit = 0
    todo_config.Filter = func(task core.Task, taskmap core.TaskMap) bool {
      return (due_time.IsZero() || (task.Due.Sub(due_time).Abs() < time.Hour*24)) && task.IsNotDone(taskmap)
    }
    if len(keys) == 0 {
      taskmap.PrintMarkdown(&todo_config)
    } else {
      seen_keys := make(map[string]bool, len(keys))
      for _, key := range keys {
        task, err := taskmap.GetTask(key)
        if err != nil {
          return err
        }
        task.PrintMarkdown(taskmap, 0, seen_keys, key, &todo_config)
      }
    }
    return nil
  },
}
