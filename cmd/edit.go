package cmd

import (
	"log/slog"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the task to add")
  editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletion)
  editCmd.ValidArgsFunction = TaskKeyCompletionOnFirst
}

var editCmd = &cobra.Command{
  Use: "edit [key] [new task]",
  Short: "edit a task",
  Run: func(cmd *cobra.Command, args []string) {
    key := args[0]
    task := strings.Join(args[1:], " ")
    if task == "" {
      task = taskmap[key].Task
    }
    taskmap.MustHaveTask(key)
    for _,key := range deps {
      taskmap.MustHaveTask(key)
    }
    taskmap[key] = core.Task {Task: task, Deps: deps}
    slog.Info("Task edited", "task", taskmap[key])
    taskmap.Write(tasks_path)
  },
}
