package core

import (
	"fmt"
	"slices"
	"time"
)

type Task struct {
  Task string           `yaml:",omitempty"`
  Deps []string         `yaml:",omitempty,flow"`
  Done bool             `yaml:",omitempty"`
  AutoComplete bool     `yaml:"auto-complete,omitempty"`
  CreatedAt time.Time   `yaml:"created-at,omitempty"`
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
  var inner string;
  if task.IsDone(taskmap) {
    inner = "x"
  } else {
    inner = " "
  }
  if key != "" {
    fmt.Printf("- [%s] %s: %s\n", inner, key, task.Task)
  } else {
    fmt.Printf("- [%s] %s\n", inner, task.Task)
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
