package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
  "time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

// add flags
var (
deps []string
dep_to []string
key string
auto_complete bool
)

func init() {
  rootCmd.AddCommand(addCmd)
  addCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the task to add")
  addCmd.Flags().StringArrayVarP(&dep_to, "dep-to", "D", []string{}, "task keys for this task to be dependent to")
  addCmd.Flags().BoolVarP(&auto_complete, "auto-complete", "a", false, "enable auto complete for the task (done when deps are done)")
  addCmd.Flags().StringVarP(&key, "key", "k", "", "key of the new task")
  addCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletion)
  addCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletion)
}

var addCmd = &cobra.Command{
  Aliases: []string{"a"},
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
    for _, dep_key := range dep_to {
      if task, ok := taskmap[dep_key]; ok {
        task.Deps = append(task.Deps, key)
        taskmap[dep_key] = task
      } else {
        slog.Error("No such task", "key", dep_key)
      }
    }
    taskmap[key] = core.Task {Task: task, Deps: deps, AutoComplete: auto_complete, CreatedAt: time.Now() }
    fmt.Printf("Task %q added\n", key)
    slog.Info("Task added", "task", taskmap[key])
    taskmap.Write(tasks_path)
  },
}
