package core

import (
	"fmt"
	"slices"
)

type Task struct {
  Task string     `yaml:",omitempty"`
  Deps []string   `yaml:",omitempty,flow"`
  Done bool       `yaml:",omitempty"`
  DoneDeps bool   `yaml:",omitempty"`
}

func (task Task) IsDone(taskmap TaskMap) bool {
  if task.DoneDeps {
    for _,key := range task.Deps {
      if !taskmap[key].IsDone(taskmap) {
        return false
      }
    }
    return true
  }
  return task.Done
}

func (task Task) PrintMarkdown(taskmap TaskMap, depth int, seen_keys *[]string) {
  var inner string;
  if task.IsDone(taskmap) {
    inner = "x"
  } else {
    inner = " "
  }
  fmt.Printf("- [%s] %s\n", inner,task.Task)
  for _, key := range task.Deps {
    if seen_keys != nil && slices.Contains(*seen_keys, key) {
      continue
    }
    for range depth {
      fmt.Print("   ")
    }
    if seen_keys != nil {
      *seen_keys = append(*seen_keys, key)
    }
    taskmap[key].PrintMarkdown(taskmap, depth+1, seen_keys)
  }
}
