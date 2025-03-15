package core

import (
	"fmt"
	"slices"
	"time"

	"github.com/nimaaskarian/ydo/utils"
)

type Task struct {
  Task string           `yaml:",omitempty"`
  Deps []string         `yaml:",omitempty,flow"`
  Done bool             `yaml:",omitempty"`
  AutoComplete bool     `yaml:"auto-complete,omitempty"`
  CreatedAt time.Time   `yaml:"created-at,omitempty"`
  Due time.Time         `yaml:",omitempty"`
}

func (task Task) IsDone(taskmap TaskMap) bool {
  if task.AutoComplete {
    for _,key := range task.Deps {
      if !taskmap[key].IsDone(taskmap) {
        return false
      }
    }
    return true
  }
  return task.Done
}

func (task Task) IsNotDone(taskmap TaskMap) bool {
  return !task.IsDone(taskmap)
}

func (task Task) PrintMarkdown(taskmap TaskMap, depth int, seen_keys *[]string, key string, filter func(task Task, taskmap TaskMap) bool) {
  if filter != nil && !filter(taskmap[key], taskmap) {
    return
  }
  for range depth {
    fmt.Print("   ")
  }
  var print_key string
  if key != "" {
    print_key = key+": "
  }
  if task.IsDone(taskmap) {
    fmt.Printf("- [x] %s%s\n", print_key,task.Task)
  } else {
    due_print := ""
    if !task.Due.IsZero() {
      diff := task.Due.Sub(time.Now())
      due_print = " ("
      if diff < 0 {
        due_print = " (-"
        diff = - diff
      }
      due_print += utils.FormatDuration(diff)
      due_print += ")"
    }
    fmt.Printf("- [ ] %s%s%s\n", print_key,task.Task, due_print)
  }
  if seen_keys != nil && slices.Contains(*seen_keys, key) {
    return
  }
  if seen_keys != nil {
    *seen_keys = append(*seen_keys, key)
  }
  for _, key := range task.Deps {
    taskmap[key].PrintMarkdown(taskmap, depth+1, seen_keys, key, filter)
  }
}
