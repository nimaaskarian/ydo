package cmd

import (
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

var no_deps bool
func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "add-deps", "d", []string{}, "add dependencies for the task")
  editCmd.Flags().BoolVarP(&no_deps, "no-deps", "D", false, "delete dependencies for a task (if used with dependencies, replaces the deps with new deps)")
  editCmd.RegisterFlagCompletionFunc("add-deps", TaskKeyCompletion)
  editCmd.ValidArgsFunction = TaskKeyCompletionOnFirst
}

var editCmd = &cobra.Command{
  Use: "edit [key] [new task message (optional)]",
  Short: "edit a task",
  Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
      slog.Error("No key to edit")
      return
    }
    key := args[0]
    task := taskmap[key]
    if new_task := strings.Join(args[1:], " "); new_task != "" {
      task.Task = new_task
    }
    taskmap.MustHaveTask(key)
    for _,key := range deps {
      taskmap.MustHaveTask(key)
    }
    if no_deps {
      task.Deps = make([]string, 0, len(deps))
    }
    task.Deps = append(task.Deps, deps...)
    taskmap[key] = task
    slog.Info("Task edited", "task", taskmap[key])
    taskmap.Write(tasks_path)
  },
}
