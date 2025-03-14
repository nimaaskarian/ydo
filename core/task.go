package core

import (
	"fmt"
	"slices"
)

type Task struct {
  Task string     `yaml:",omitempty"`
  Deps []string   `yaml:",omitempty,flow"`
  Done bool       `yaml:",omitempty"`
  AutoComplete bool   `yaml:"auto-complete,omitempty"`
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

func (task Task) PrintMarkdown(taskmap TaskMap, depth int, seen_keys *[]string, key string) {
  var inner string;
  if task.IsDone(taskmap) {
    inner = "x"
  } else {
    inner = " "
  }
  if key != "" {
    key = key + ": "
  }
  fmt.Printf("- [%s] %s%s\n", inner, key, task.Task)
  if seen_keys != nil && slices.Contains(*seen_keys, key) {
    return
  }
  for _, key := range task.Deps {
    for range depth {
      fmt.Print("   ")
    }
    if seen_keys != nil {
      *seen_keys = append(*seen_keys, key)
    }
    taskmap[key].PrintMarkdown(taskmap, depth+1, seen_keys, key)
  }
}
