package core;
import (
  "fmt"
)

type Todo struct {
  Task string
  Deps []string
  Done bool
  DoneDeps bool
}

func (todo Todo) IsDone(todomap TodoMap) bool {
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

func (todo Todo) PrintMarkdown(todomap TodoMap, depth int) {
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
