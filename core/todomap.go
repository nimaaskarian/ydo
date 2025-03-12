package core

import (
  "gopkg.in/yaml.v3"
  "log"
  "fmt"
  "strconv"
	"slices"
)

func ParseYaml(obj any, input []byte) {
  err := yaml.Unmarshal([]byte(input), obj)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

type TaskMap map[string]Task;

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

func (taskmap TaskMap) PrintMarkdown() {
  var seen_keys []string
  for _,todo := range taskmap {
    seen_keys = append(seen_keys, todo.Deps...)
  }
  for i,todo := range taskmap {
    if !slices.Contains(seen_keys, i) {
      todo.PrintMarkdown(taskmap, 1)
    }
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
