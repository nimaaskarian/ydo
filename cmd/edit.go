package cmd

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var new_key string
var remove_deps, remove_dep_to bool
func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "append dependencies for the task")
  editCmd.Flags().StringArrayVarP(&dep_to, "dep-to", "D", []string{}, "append task keys for this task to be dependent to")
  editCmd.Flags().BoolVarP(&remove_deps, "remove-deps", "r", false, "remove previous dependencies for the task. using this with --deps causes to replace dependencies")
  editCmd.Flags().BoolVarP(&remove_dep_to, "remove-dep-to", "R", false, "remove previous 'dependent to' for the task. using this with --dep-to causes to replace 'dependent to's")
  editCmd.Flags().StringVarP(&new_key, "new-key", "k", "", "new key to the task")
  editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletion)
  editCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletion)
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
    if remove_deps {
      task.Deps = make([]string, 0, len(deps))
    }
    if remove_dep_to {
      for dep_key, task := range taskmap {
        index := slices.IndexFunc(task.Deps, func(dep string)bool {
          return dep == key
        })
        if index == -1 {
          continue
        }
        task.Deps = slices.Delete(task.Deps, index, index+1)
        taskmap[dep_key] = task
      }
    }
    for _, dep_key := range dep_to {
      if task, ok := taskmap[dep_key]; ok {
        task.Deps = append(task.Deps, key)
        taskmap[dep_key] = task
      } else {
        slog.Error("No such task", "key", dep_key)
      }
    }
    if new_key != "" {
      for dep_key, task := range taskmap {
        index := slices.IndexFunc(task.Deps, func(dep string)bool {
          return dep == key
        })
        if index == -1 {
          continue
        }
        task.Deps = slices.Replace(task.Deps, index, index+1, new_key)
        taskmap[dep_key] = task
      }
      delete(taskmap, key)
      key = new_key
    }
    task.Deps = append(task.Deps, deps...)
    taskmap[key] = task
    slog.Info("Task edited", "task", taskmap[key])
    taskmap.Write(tasks_path)
  },
}
