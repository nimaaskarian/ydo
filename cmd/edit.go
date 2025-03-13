package cmd

import (
	"log/slog"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the task to add")
  editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletion)
}

var editCmd = &cobra.Command{
  Use: "edit",
  Short: "edit a task",
  Run: func(cmd *cobra.Command, args []string) {
    utils.MustHaveTask(taskmap, key)
    task := strings.Join(args, " ")
    if task == "" {
      task = taskmap[key].Task
    }
    for _,key := range deps {
      utils.MustHaveTask(taskmap, key)
    }
    taskmap[key] = core.Task {Task: task, Deps: deps}
    utils.WriteTaskmap(taskmap, tasks_path)
    slog.Info("Task created", "key", key)
  },
}
