package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

// add flags
var (
deps []string
dep_tos []string
key string
due string
auto_complete bool
description string
tfidf bool
taskmsg string 
recur string
)

func init() {
  rootCmd.AddCommand(addCmd)
  addCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the task to add")
  addCmd.Flags().StringVarP(&description, "description", "e", "", "description of the task")
  addCmd.Flags().StringArrayVarP(&dep_tos, "dep-to", "D", []string{}, "task keys for this task to be dependent to")
  addCmd.Flags().BoolVarP(&auto_complete, "auto-complete", "a", false, "enable auto complete for the task (done when deps are done)")
  addCmd.Flags().BoolVarP(&tfidf, "tfidf", "t", false, "use tfidf for automatic key generation (overrides config file and --key flag)")
  addCmd.Flags().StringVarP(&key, "key", "k", "", "key of the new task")
  addCmd.RegisterFlagCompletionFunc("key", KeyCompletion)

  addCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletionFilter(nil))
  addCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletionFilter(nil))

  addCmd.Flags().StringVarP(&due, "due", "u", "", "specify due for the tasks to print")
  addCmd.RegisterFlagCompletionFunc("due", DueCompletion)

  addCmd.Flags().StringVarP(&recur, "recur", "r", "", "duration of in which the ask recurs")
  addCmd.RegisterFlagCompletionFunc("recur", DurationCompletion)
}

var addCmd = &cobra.Command{
  Aliases: []string{"a"},
  Use: "add [your task here yay]",
  Short: "add a task",
  Args: func(cmd *cobra.Command, args []string) (err error) {
    if err = cobra.MinimumNArgs(1)(cmd, args); err != nil {
      return err
    }
    taskmsg, err = TaskTitleFromArgs(args)
    return err
  },
  RunE: func(cmd *cobra.Command, args []string) error {
    if tfidf {
      key = taskmap.TfidfNextKey(taskmsg, core.TfidfConfig {Enabled: true}, "")
    } else {
      if key == "" {
        key = taskmap.TfidfNextKey(taskmsg, config.Tfidf, "")
      }
    }
    for _,dep := range deps {
      if _, err := taskmap.GetTask(dep); err != nil {
        return err
      }
    }
    due_time, err := utils.ParseDue(due, time.Now())
    if err != nil {
      return err
    }
    if _, err := utils.ParseDuration(recur, time.Now()); err != nil {
      return err
    }
    taskmap[key] = core.Task {Task: taskmsg, Deps: deps, AutoComplete: auto_complete, CreatedAt: time.Now(), Due: due_time, Description: description, Recur: recur }
    for _, dep_to := range dep_tos {
      task, err := taskmap.AddDep(dep_to, key)
      if err != nil {
        return err
      }
      if ; err != nil {
        return err
      }
      taskmap[dep_to] = task
    }
    fmt.Printf("Task %q added\n", key)
    slog.Debug("Added a task", "task", taskmap[key])
    return nil
  },
  PostRunE: SaveChanges,
  PreRun: UpdateOldTaskMap,
}

func TaskTitleFromArgs(args []string) (taskmsg string, err error) {
  has_non_empty := slices.ContainsFunc(args, func (arg string) bool {
    return arg != ""
  })
  if !has_non_empty {
    return "", errors.New("Task cannot be empty")
  }
  // in task title, replace all newlines with spaces
  return strings.Replace(strings.Join(args, " "), "\n", " ", 0), nil
}
