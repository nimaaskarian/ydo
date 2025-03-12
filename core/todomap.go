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

type TodoMap map[string]Todo;

func (todomap TodoMap) Do(key string) {
  if entry, ok := todomap[key]; ok {
    entry.Done = true
    todomap[key] = entry
  }
}

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

func (todomap TodoMap) NextKey() string {
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

func (todomap TodoMap) PrintKeys() {
  for key := range todomap {
    fmt.Printf("%s\n", key)
  }
}
