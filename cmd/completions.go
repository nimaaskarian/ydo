package cmd

import (
  "strings"
	"github.com/nimaaskarian/ydo/core"

	"github.com/spf13/cobra"
)

func TaskKeyCompletionFilter(filter func(core.Task, core.TaskMap) bool) func (*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
  return func (cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    taskmap = core.LoadTaskMap(tasks_path)
    keys := make([]string, len(taskmap))
    i := 0
    for key := range taskmap {
      if (filter == nil || filter(taskmap[key], taskmap)) && strings.HasPrefix(key, toComplete) {
        keys[i] = key
        i++
      }
    }
    return keys, cobra.ShellCompDirectiveDefault
  }
}

func TaskKeyCompletionOnFirst(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  if len(args) > 0 {
    return []string{}, cobra.ShellCompDirectiveDefault
  }
  return TaskKeyCompletionFilter(nil)(cmd, args, toComplete)
}

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

