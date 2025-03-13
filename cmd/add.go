package cmd

import (
	"log/slog"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var deps []string
func init() {
  rootCmd.AddCommand(addCmd)
  addCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the task to add")
  addCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletion)

}

var addCmd = &cobra.Command{
  Use: "add",
  Short: "add a task",
  Run: func(cmd *cobra.Command, args []string) {
    if key == "" {
      key = taskmap.NextKey()
    }
    utils.MustNotHaveTask(taskmap, key)
    task := strings.Join(args, " ")
    if task == "" {
      slog.Error("Task can't be empty")
      panic("task can not be empty")
    }
    for _,key := range deps {
      utils.MustHaveTask(taskmap, key)
    }
    taskmap[key] = core.Task {Task: task, Deps: deps}
    utils.WriteTaskmap(taskmap, tasks_path)
    slog.Info("Task created", "key", key)
  },
}
