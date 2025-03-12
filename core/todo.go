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

func (todo Task) IsDone(todomap TaskMap) bool {
  if todo.DoneDeps {
    for _,key := range todo.Deps {
      if !todomap[key].Done {
        return false
      }
    }
    return true
  }
  return todo.Done
}

func (todo Task) PrintMarkdown(todomap TaskMap, depth int) {
  var inner string;
  if todo.IsDone(todomap) {
    inner = "x"
  } else {
    inner = " "
  }
  fmt.Printf("- [%s] %s\n", inner,todo.Task)
  for _, id := range todo.Deps {
    for range depth {
      fmt.Print("   ")
    }
    todomap[id].PrintMarkdown(todomap, depth+1)
  }
}
