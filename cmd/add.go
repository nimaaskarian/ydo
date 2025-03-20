package cmd

import (
	"fmt"
	"log/slog"
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
tfidf bool
)

func DueCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  date := []string{
  "saturday",
  "sunday",
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sun",
  "mon",
  "tue",
  "wed",
  "thu",
  "fri",
  "today",
  "tomorrow",
  "yesterday",
  "later",
  }
  time := []string {
    "now",
    "8:00",
    "20:00",
  }
  var out []string
  if strings.Contains(toComplete, "/") && !strings.HasPrefix(toComplete, "later/"){
    out = make([]string, 0, len(time)*len(date))
    for _, date := range date {
      for _, time := range time {
        out = append(out, date+"/"+time)
      }
    }
  } else {
    out = date
  }
  return out, cobra.ShellCompDirectiveDefault
}

func init() {
  rootCmd.AddCommand(addCmd)
  addCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the task to add")
  addCmd.Flags().StringVarP(&due, "due", "u", "", "specify due for the tasks to print")
  addCmd.Flags().StringArrayVarP(&dep_tos, "dep-to", "D", []string{}, "task keys for this task to be dependent to")
  addCmd.Flags().BoolVarP(&auto_complete, "auto-complete", "a", false, "enable auto complete for the task (done when deps are done)")
  addCmd.Flags().BoolVarP(&tfidf, "tfidf", "t", false, "use tfidf for automatic key generation (overrides config file and --key flag)")
  addCmd.Flags().StringVarP(&key, "key", "k", "", "key of the new task")
  addCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletionFilter(nil))
  addCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletionFilter(nil))
  addCmd.RegisterFlagCompletionFunc("due", DueCompletion)
}

var addCmd = &cobra.Command{
  Aliases: []string{"a"},
  Use: "add [your task here yay]",
  Short: "add a task",
  Args: cobra.MinimumNArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
    task := strings.Join(args, " ")
    if tfidf {
      key = taskmap.TfidfNextKey(task, core.TfidfConfig {Enabled: true}, "")
    } else {
      if key == "" {
        key = taskmap.TfidfNextKey(task, config.Tfidf, "")
      }
    }
    for _,key := range deps {
      if _, err := taskmap.GetTask(key); err != nil {
        return err
      }
    }
    due_time, err := utils.ParseDate(due)
    if err != nil {
      return err
    }
    taskmap[key] = core.Task {Task: task, Deps: deps, AutoComplete: auto_complete, CreatedAt: time.Now(), Due: due_time }
    for _, dep_to := range dep_tos {
      task, err := taskmap.GetTask(dep_to)
      if err != nil {
        return err
      }
      if err := task.AddDep(taskmap, key); err != nil {
        return err
      }
      taskmap[dep_to] = task
    }
    fmt.Printf("Task %q added\n", key)
    slog.Debug("Added a task", "task", taskmap[key])
    return nil
  },
}
