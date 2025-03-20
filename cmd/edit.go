package cmd

import (
	"errors"
	"log/slog"
	"reflect"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// edit flags
var (
  new_key string
  remove_deps bool
  key_regen bool
  remove_dep_to bool
  no_auto_complete bool
)

func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "append dependencies for the task")
  editCmd.Flags().StringArrayVarP(&dep_tos, "dep-to", "D", []string{}, "append task keys for this task to be dependent to")
  editCmd.Flags().BoolVarP(&remove_deps, "remove-deps", "r", false, "remove previous dependencies for the task. using this with --deps causes to replace dependencies")
  editCmd.Flags().BoolVarP(&key_regen, "key-regen", "K", false, "regen key using the automatic next key generator (respects the config file)")
  editCmd.Flags().BoolVarP(&remove_dep_to, "remove-dep-to", "R", false, "remove previous 'dependent to' for the task. using this with --dep-to causes to replace 'dependent to's")
  editCmd.Flags().BoolVarP(&auto_complete, "auto-complete", "a", false, "enable auto complete for the task")
  editCmd.Flags().BoolVarP(&no_auto_complete, "no-auto-complete", "A", false, "disable auto complete for the task")
  editCmd.MarkFlagsMutuallyExclusive("no-auto-complete", "auto-complete")
  editCmd.Flags().StringVarP(&new_key, "key", "k", "", "new key to the task")
  editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletionFilter(nil))
  editCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletionFilter(nil))
  editCmd.ValidArgsFunction = TaskKeyCompletionOnFirst
}

var editCmd = &cobra.Command{
  Aliases: []string{"e"},
  Use: "edit [key] [new task message (optional)]",
  Short: "edit a task",
  Args: cobra.MinimumNArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
    if _, err := taskmap.GetTask(args[0]); err != nil {
      return err
    }
    edit_key := args[0]
    task := taskmap[edit_key]
    if new_task := strings.Join(args[1:], " "); new_task != "" {
      task.Task = new_task
    }
    if key_regen {
      new_key = taskmap.TfidfNextKey(task.Task, config.Tfidf, edit_key)
    }
    if _, err := taskmap.GetTask(edit_key); err != nil {
      return err
    }
    for _,dep := range deps {
      if _, err := taskmap.GetTask(dep); err != nil {
        return err
      }
    }
    if remove_deps {
      task.Deps = make([]string, 0, len(deps))
    }
    if remove_dep_to {
      taskmap.WipeDependenciesToKey(edit_key)
    }
    for _, dep_key := range dep_tos {
      task, err := taskmap.GetTask(dep_key)
      if err != nil {
        return err
      }
      if !slices.Contains(task.Deps, edit_key) {
        task.AddDep(taskmap, edit_key)
      }
      taskmap[dep_key] = task
    }
    edit_key = taskmap.ReplaceKeyInDeps(edit_key, new_key)
    if auto_complete {
      task.AutoComplete = true
    }
    if no_auto_complete {
      task.AutoComplete = false
    }
    task.Deps = append(task.Deps, deps...)
    taskmap[edit_key] = task
    if reflect.DeepEqual(taskmap, old_taskmap) {
      return errors.New("Not edited")
    }
    slog.Info("Task edited", "task", task)
    return nil
  },
}
