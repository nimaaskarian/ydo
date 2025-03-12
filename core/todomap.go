package core

import (
  "gopkg.in/yaml.v3"
  "log"
  "fmt"
  "strconv"
)

func ParseYaml(obj any, input []byte) {
  err := yaml.Unmarshal([]byte(input), obj)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

type TaskMap map[string]Todo;

func (todomap TaskMap) Do(key string) {
  if entry, ok := todomap[key]; ok {
    entry.Done = true
    todomap[key] = entry
  }
}

func PrintYaml(obj any) {
  s,err:=yaml.Marshal(obj)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
  fmt.Printf("%s", s)
}

func (todomap TaskMap) PrintMarkdown() {
  for _,todo := range todomap {
    todo.PrintMarkdown(todomap,1)
  }
}

func (todomap TaskMap) NextKey() string {
  i := 1
  for {
    s_i := "t"+strconv.Itoa(i);
    if _, ok := todomap[s_i]; ok {
      i++
    } else {
      return s_i
    }
  }
}

func (todomap TaskMap) PrintKeys() {
  for key := range todomap {
    fmt.Printf("%s\n", key)
  }
}
