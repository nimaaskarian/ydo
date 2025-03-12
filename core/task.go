package core;
import (
  "fmt"
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

func (task Task) PrintMarkdown(taskmap TaskMap, depth int) {
  var inner string;
  if task.IsDone(taskmap) {
    inner = "x"
  } else {
    inner = " "
  }
  fmt.Printf("- [%s] %s\n", inner,task.Task)
  for _, id := range task.Deps {
    for range depth {
      fmt.Print("   ")
    }
    taskmap[id].PrintMarkdown(taskmap, depth+1)
  }
}
