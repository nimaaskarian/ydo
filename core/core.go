package core

import (
  "gopkg.in/yaml.v3"
  "log"
  "fmt"
  "strconv"
)

type Todo struct {
  Task string
  Deps []string
  Done bool
}

func (todo Todo) PrintMarkdown(todomap TodoMap, depth int) {
  var inner string;
  if todo.Done {
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
type TodoMap map[string]Todo;

func (todomap TodoMap) PrintYaml() {
  s,err:=yaml.Marshal(todomap)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
  fmt.Printf("%s", s)
}

func (todomap TodoMap) PrintMarkdown() {
  for _,todo := range todomap {
    todo.PrintMarkdown(todomap,1)
  }
}

func ParseYaml(obj any, input []byte) {
  err := yaml.Unmarshal([]byte(input), obj)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

func (todomap TodoMap) NextId() string {
  i := 0
  for {
    s_i := strconv.Itoa(i);
    if _, ok := todomap[s_i]; ok {
      i++
    } else {
      return s_i
    }
  }
}
