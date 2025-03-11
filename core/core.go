package core

import (
  "gopkg.in/yaml.v3"
  "log"
)

type Todo struct {
  Message string
  Priority int
  Deps []int
  Done bool
}

type TodoMap map[int]Todo;

func (todo *Todo) YamlParse(input string) {
  err := yaml.Unmarshal([]byte(input), &todo)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

func (todomap TodoMap) YamlParseMap(input string) {
  err := yaml.Unmarshal([]byte(input), &todomap)
  if err != nil {
    log.Fatalf("%v\n", err)
  }
}

func (todomap TodoMap) NextId() int {
  i := 1
  for {
    if _, ok := todomap[i]; ok {
      i++
    } else {
      return i
    }
  }
}
