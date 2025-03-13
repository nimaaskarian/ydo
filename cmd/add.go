package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/nimaaskarian/ydo/core"
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
    task := strings.Join(args, " ")
    if task == "" {
      slog.Error("Task can't be empty")
      os.Exit(1)
    }
    taskmap.MustNotHaveTask(key)
    for _,key := range deps {
      taskmap.MustHaveTask(key)
    }
    taskmap[key] = core.Task {Task: task, Deps: deps}
    slog.Info("Task added", "task", taskmap[key])
    taskmap.Write(tasks_path)
  },
}
