package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// edit flags
var (
new_key string
remove_deps bool
remove_dep_to bool
no_auto_complete bool
)

func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "append dependencies for the task")
  editCmd.Flags().StringArrayVarP(&dep_to, "dep-to", "D", []string{}, "append task keys for this task to be dependent to")
  editCmd.Flags().BoolVarP(&remove_deps, "remove-deps", "r", false, "remove previous dependencies for the task. using this with --deps causes to replace dependencies")
  editCmd.Flags().BoolVarP(&remove_dep_to, "remove-dep-to", "R", false, "remove previous 'dependent to' for the task. using this with --dep-to causes to replace 'dependent to's")
  editCmd.Flags().BoolVarP(&auto_complete, "auto-complete", "a", false, "enable auto complete for the task (done when deps are done)")
  editCmd.Flags().BoolVarP(&no_auto_complete, "no-auto-complete", "A", false, "disable auto complete for the task (done when deps are done)")
  editCmd.Flags().StringVarP(&new_key, "key", "k", "", "new key to the task")
  editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletion)
  editCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletion)
  editCmd.ValidArgsFunction = TaskKeyCompletionOnFirst
}

var editCmd = &cobra.Command{
  Aliases: []string{"e"},
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
      taskmap.WipeDependenciesToKey(key)
    }
    dep_to_changed := false
    for _, dep_key := range dep_to {
      if task, ok := taskmap[dep_key]; ok {
        if !slices.Contains(task.Deps, key) {
          dep_to_changed = true
          task.Deps = append(task.Deps, key)
          taskmap[dep_key] = task
        }
      } else {
        slog.Error("No such task", "key", dep_key)
      }
    }
    if new_key != "" {
      for dep_key, task := range taskmap {
        index := slices.Index(task.Deps, key)
        if index != -1 {
          task.Deps = slices.Replace(task.Deps, index, index+1, new_key)
          taskmap[dep_key] = task
        }
      }
      delete(taskmap, key)
      key = new_key
    }
    if auto_complete && no_auto_complete {
      slog.Error("Can't use both auto-complete and no-auto-complete flags at the same time")
      os.Exit(1)
    }
    if auto_complete {
      task.AutoComplete = true
    }
    if no_auto_complete {
      task.AutoComplete = false
    }
    task.Deps = append(task.Deps, deps...)
    if !dep_to_changed && reflect.DeepEqual(taskmap[key], task) {
      fmt.Println("Task not edited")
    } else {
      taskmap[key] = task
      slog.Info("Task edited", "task", task)
      rootCmd.Run(cmd, args)
      taskmap.Write(tasks_path)
    }
  },
}
