package core

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"

	"github.com/nimaaskarian/ydo/utils"
	"gopkg.in/yaml.v3"
)

func ParseYaml(obj any, input []byte) {
  err := yaml.Unmarshal([]byte(input), obj)
  if err != nil {
    slog.Error("Failed unmarshaling yaml", "err", err)
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
    slog.Info("Completed task","key" ,key, "task", entry.Task)
  } else {
    slog.Error("No such task", "key",key)
    os.Exit(1)
  }
}

func (taskmap TaskMap) Undo(key string) {
  if entry, ok := taskmap[key]; ok {
    entry.Done = false
    taskmap[key] = entry
    slog.Info("Un-completed task","key" ,key, "task", entry.Task)
  } else {
    slog.Error("No such task", "key",key)
    os.Exit(1)
  }
}

func PrintYaml(obj any) {
  s,err:=yaml.Marshal(obj)
  if err != nil {
    slog.Error("Failed marshaling yaml", "err", err)
  }
  fmt.Printf("%s", s)
}

func (taskmap TaskMap) Depth(key string, visited []string) int {
  depth := 0
  depth += len(taskmap[key].Deps)
  for _,key := range taskmap[key].Deps {
    if slices.Contains(visited, key) {
      continue
    }
    visited = append(visited, key)
    depth += taskmap.Depth(key, visited)
  }
  return depth
}
func (taskmap TaskMap) PrintMarkdown() {
  depths := make(map[string]int, len(taskmap))
  keys := make([]string, 0 ,len(taskmap))
  for key := range taskmap {
    depths[key] = taskmap.Depth(key, []string{})
    keys = append(keys, key)
    if depths[key] == len(taskmap) {
      keys = []string{key}
      break
    }
  }
  sort.SliceStable(keys, func(i, j int) bool {
    return depths[keys[i]] > depths[keys[j]]
  })

  var seen_keys []string
  var seen_in_deps []string
  for _,key := range keys {
    if !slices.Contains(seen_keys, key) && !slices.Contains(seen_in_deps, key) {
      seen_keys = append(seen_keys, key)
      taskmap[key].PrintMarkdown(taskmap, 1, seen_keys, &seen_in_deps)
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

func (taskmap TaskMap) MustHaveTask(key string) {
  if !taskmap.HasTask(key) {
    slog.Error("No such task", "key",key)
    os.Exit(1)
  }
}

func (taskmap TaskMap) MustNotHaveTask(key string) {
  if taskmap.HasTask(key) {
    slog.Error("Task already exists", "key",key)
    os.Exit(1)
  }
}

func (taskmap TaskMap) Write(path string) {
  content, err := yaml.Marshal(taskmap)
  utils.Check(err)
  err = os.WriteFile(path, content, 0644)
  utils.Check(err)
  slog.Info("Wrote to file", "path", path)
}

func LoadTaskMap(path string) TaskMap {
  slog.Info("Task file loaded.", "path", path)
  taskmap := make(TaskMap)
  content, _ := os.ReadFile(path)
  ParseYaml(taskmap, content)
  return taskmap
}
