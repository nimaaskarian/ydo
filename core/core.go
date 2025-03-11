package core

import (
  "gopkg.in/yaml.v3"
  "log"
)

type Todo struct {
  Task string
  Priority int
  Deps []int
  Done bool
}
type TodoMap map[int]Todo;

func ParseYaml(obj any, input []byte) {
  err := yaml.Unmarshal([]byte(input), obj)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

func (todomap TodoMap) NextId() int {
  i := 0
  for {
    if _, ok := todomap[i]; ok {
      i++
    } else {
      return i
    }
  }
}
