package core

import (
	"cmp"
	"fmt"
	"log/slog"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseYaml(obj any, input []byte) {
  err := yaml.Unmarshal([]byte(input), obj)
  if err != nil {
    slog.Error("Failed unmarshaling yaml", "err", err)
  }
}

type TaskMap map[string]Task;

func (taskmap TaskMap) WipeDependenciesToKey(key string) {
    for each_key, task := range taskmap {
    index := slices.Index(task.Deps, key)
    if index != -1 {
      task.Deps = slices.Delete(task.Deps, index, index+1)
      taskmap[each_key] = task
    }
  }
}

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
    slog.Error("No such task")
    os.Exit(1)
  }
}

func (taskmap TaskMap) Undo(key string) {
  if entry, ok := taskmap[key]; ok {
    entry.Done = false
    taskmap[key] = entry
    slog.Info("Un-completed task","key" ,key, "task", entry.Task)
  } else {
    slog.Error("No such task")
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

func (taskmap TaskMap) PrintMarkdown(filter func(task Task, taskmap TaskMap) bool) {
  if len(taskmap) == 0 {
    fmt.Println("No tasks found")
    return
  }
  keys := make([]string, 0 ,len(taskmap))
  for key := range taskmap {
    keys = append(keys, key)
  }
  slices.SortFunc(keys, func(k1, k2 string) int {
    return taskmap[k1].CreatedAt.Compare(taskmap[k2].CreatedAt)
  })

  seen_keys := make([]string, 0, len(taskmap))
  for _,key := range keys {
    if !slices.Contains(seen_keys, key) {
      taskmap[key].PrintMarkdown(taskmap, 0, &seen_keys, key, filter)
    }
  }
}

func (taskmap TaskMap) NextKey() string {
  i := 1
  for {
    key := "t"+strconv.Itoa(i);
    if _, ok := taskmap[key]; ok {
      i++
    } else {
      return key
    }
  }
}

type TfidfConfig struct {
  Enabled bool `yaml:",omitempty"`
  MinTaskCount int `yaml:"min-task-count,omitempty"`
}

func (taskmap TaskMap) TfidfNextKey(task string, config TfidfConfig, current_key string) string {
  if config.Enabled && len(taskmap) >= config.MinTaskCount {
    words := strings.Split(task, " ")
    word_count_in_docs := make(map[string]int, len(words))
    for key, task := range taskmap {
      if key == current_key {
        continue
      }
      for _, word := range words {
        count := word_count_in_docs[word]
        if strings.Contains(task.Task, word) {
          word_count_in_docs[word] = count + 1
        }
      }
    }
    num_tasks := len(taskmap)+1
    idf_map := make(map[string]float64, len(words))
    for _, word := range words {
      count := word_count_in_docs[word]+1
      idf_map[word] = math.Log(float64(num_tasks/count))
    }
    word_count_in_current := make(map[string]int, len(words))
    for _, word := range words {
      word_count_in_current[word] += 1
    }
    tfidf_map := make(map[string]float64, len(words))
    for _, word := range words {
      tfidf_map[word] = float64(word_count_in_current[word])/float64(len(words)) * idf_map[word]
    }
    slices.SortFunc(words, func(a,b string) int {
      return cmp.Compare(tfidf_map[b], tfidf_map[a])
    })
    slog.Info("Tfidf calculated","tfidf", tfidf_map)
    for _, word := range words {
      if _, ok := taskmap[word]; !ok {
        return word
      } else {
        if word == current_key {
          return word
        }
      }
    }
  }
  slog.Info("Tfidf fallback to NextKey.")
  return taskmap.NextKey()
}

func (taskmap TaskMap) PrintKeys() {
  for key := range taskmap {
    fmt.Printf("%s\n", key)
  }
}

func (taskmap TaskMap) MustHaveTask(key string) {
  if !taskmap.HasTask(key) {
    slog.Error("No such task")
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
  if err != nil {
    slog.Error("Failed converting the tasks to yaml.")
    os.Exit(1)
  }
  if os.WriteFile(path, content, 0644) != nil {
    slog.Error("Failed writing the tasks to file.", "path", path)
    os.Exit(1)
  }
  slog.Info("Wrote to file", "path", path)
}

func LoadTaskMap(path string) TaskMap {
  slog.Info("Task file loaded.", "path", path)
  taskmap := make(TaskMap)
  content, _ := os.ReadFile(path)
  ParseYaml(taskmap, content)
  return taskmap
}
