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

func (taskmap TaskMap) HasTask(key string) bool {
  _, ok := taskmap[key]
  return ok
}

func (taskmap TaskMap) Do(key string) {
  if entry, ok := taskmap[key]; ok {
    entry.Done = true
    taskmap[key] = entry
    log.Printf("Completed task %q (%q)\n", key, entry.Task)
  }
}

func (taskmap TaskMap) Undo(key string) {
  if entry, ok := taskmap[key]; ok {
    entry.Done = false
    taskmap[key] = entry
    log.Printf("Un-completed task %q (%q)\n", key, taskmap[key].Task)
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
  for _,task := range taskmap {
    seen_keys = append(seen_keys, task.Deps...)
  }
  for i,task := range taskmap {
    if !slices.Contains(seen_keys, i) {
      task.PrintMarkdown(taskmap, 1)
    }
  }
}

func (taskmap TaskMap) NextKey() string {
  i := 1
  for {
    s_i := "t"+strconv.Itoa(i);
    if _, ok := taskmap[s_i]; ok {
      i++
    } else {
      return s_i
    }
  }
}

func (taskmap TaskMap) PrintKeys() {
  for key := range taskmap {
    fmt.Printf("%s\n", key)
  }
}
