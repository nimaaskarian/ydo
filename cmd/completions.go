package cmd

import (
	"slices"
	"strings"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"

	"github.com/spf13/cobra"
)

func TaskKeyCompletionFilter(filter core.TaskFilter) cobra.CompletionFunc {
  return func (cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    taskmap = core.LoadTaskMap(tasks_path, time.Now())
    now = time.Now()
    if now_str != "" {
      var err error
      now, err = utils.ParseDue(now_str, now)
      if err != nil {
        return []string{}, cobra.ShellCompDirectiveError
      }
    }

    keys := make([]string, 0, len(taskmap))
    i := 0
    for key := range taskmap {
      if (filter == nil || filter(taskmap[key], taskmap, now)) && !slices.Contains(args, key) {
        keys = append(keys, key)
        i++
      }
    }
    return keys, cobra.ShellCompDirectiveNoFileComp
  }
}

func TaskKeyCompletionOnFirst(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  if len(args) > 0 {
    return []string{}, cobra.ShellCompDirectiveDefault
  }
  return TaskKeyCompletionFilter(nil)(cmd, args, toComplete)
}

func DurationCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  durations := [...]string{
    "months",
    "yrs",
    "wks",
    "days",
    "hrs",
    "mins",
    "scnds",
  }
  out := make([]string, 0, len(durations))
  index := strings.IndexFunc(toComplete, utils.IsDigit)
  if index != 0 && len(toComplete) != 0 {
    if index == -1 {
      index = len(toComplete)
    }
    for _,duration := range durations {
      out = append(out, toComplete[:index]+duration)
    }
  }
  return out, cobra.ShellCompDirectiveNoFileComp
}

func DueCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  date := [...]string{
  "sat",
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
  time := [...]string {
    "now",
  }
  durations, _ := DurationCompletion(cmd, args, toComplete)
  out := make([]string, 0, len(time)*len(date)+len(durations))
  out = append(out, durations...)
  if strings.Contains(toComplete, "/") {
    for _, date := range date {
      for _, time := range time {
        out = append(out, date+"/"+time)
      }
    }
  } else {
    for _, date := range date {
      out = append(out, date)
    }
  }
  return out, cobra.ShellCompDirectiveNoFileComp
}

func KeyCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  taskmap = core.LoadTaskMap(tasks_path, now)
  var words []string
  if cmd == editCmd {
    edit_key := args[0]
    if len(args)-1 > 0 {
      words = args[1:]
    }
    task := taskmap[edit_key]
    words = strings.Fields(task.Task.Value())
  }
  if cmd == addCmd {
    words = args
  }
  if words == nil {
    return nil, cobra.ShellCompDirectiveError
  }
  if depth := strings.Count(toComplete, "-"); depth != 0 {
    init_words := words
    for range depth {
      for _, word_i := range words {
        included_words := strings.Split(word_i, "-")
        for _, word_j := range init_words {
          if !slices.Contains(included_words, word_j) {
            words = append(words, word_i+"-"+word_j)
          }
        }
      }
    }
  }
  for _, word := range words {
    words = append(words, word+"-")
  }
  return words, cobra.ShellCompDirectiveDefault
}

func TagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	taskmap = core.LoadTaskMap(tasks_path, now)
	return taskmap.Tags(), cobra.ShellCompDirectiveDefault
}
