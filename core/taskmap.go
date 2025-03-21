package core

import (
	"cmp"
	"errors"
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
    panic("Failed unmarshaling yaml");
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

type NoSuchTask struct {}
func (e NoSuchTask) Error() string {
  return "No such task"
}

func (taskmap TaskMap) GetTask(key string) (Task, error) {
  task, ok := taskmap[key]
  if !ok {
    return task, NoSuchTask{}
  }
  return task, nil
}

func (taskmap TaskMap) HasTask(key string) bool {
  _, ok := taskmap[key]
  return ok
}

func (taskmap TaskMap) Do(key string) error {
  task, err := taskmap.GetTask(key)
  if err != nil{
    return err
  }
  task.Do()
  taskmap[key] = task
  slog.Info("Completed task","key" ,key)
  return nil
}

func (taskmap TaskMap) Undo(key string) error {
  task, err := taskmap.GetTask(key)
  if err != nil{
    return err
  }
  task.Undo()
  taskmap[key] = task
  slog.Info("Un-completed task","key" ,key)
  return nil
}

func PrintYaml(obj any) error {
  s,err:=yaml.Marshal(obj)
  if err != nil {
    slog.Error("Failed marshaling yaml", "err", err)
    return err
  }
  fmt.Printf("%s", s)
  return nil
}

type MarkdownFilter func(task Task, taskmap TaskMap) bool;

type MarkdownConfig struct {
  Indent uint `yaml:",omitempty"`
  Mode string `yaml:",omitempty"`
  Filter MarkdownFilter
  file *os.File
}

func (taskmap TaskMap) PrintMarkdown(config *MarkdownConfig) error {
  if len(taskmap) == 0 {
    return errors.New("No tasks found")
  }
  if config.file == nil {
    config.file = os.Stdout
  }
  keys := make([]string, 0 ,len(taskmap))
  for key := range taskmap {
    keys = append(keys, key)
  }
  slices.SortFunc(keys, func(k1, k2 string) int {
    return taskmap[k1].CreatedAt.Compare(taskmap[k2].CreatedAt)
  })

  seen_keys := make(map[string]bool, len(taskmap))
  for _,key := range keys {
    if value, ok := seen_keys[key]; !ok || !value {
      taskmap[key].PrintMarkdown(taskmap, 0, seen_keys, key, config)
    }
  }
  return nil
}

func (taskmap TaskMap) NextKey(current string) string {
  i := 1
  for {
    key := "t"+strconv.Itoa(i);
    if _, ok := taskmap[key]; (!ok || (current != "" && key == current)){
      return key
    } else {
      i++
    }
  }
}

type TfidfConfig struct {
  Enabled bool `yaml:",omitempty"`
  MinTaskCount int `yaml:"min-task-count,omitempty"`
}

// gets a config, if the config says the tfidf should be enabled, returns
// tfidf. 
// it skips current_key and also falls back to TaskMap.NextKey if config is
// says it should be disabled
func (taskmap TaskMap) TfidfNextKey(task string, config TfidfConfig, current_key string) string {
  if config.Enabled && len(taskmap) >= config.MinTaskCount {
    words := strings.Fields(task)
    word_count_in_docs := make(map[string]int, len(words))
    for key, task := range taskmap {
      if key == current_key {
        continue
      }
      for _, word := range words {
        if strings.Contains(task.Task, word) {
          word_count_in_docs[word] += 1
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
  return taskmap.NextKey(current_key)
}

// replaces keys in dependencies of the whole list. returns the new key if the 
// transition went good, the old key if not
func (taskmap TaskMap) ReplaceKeyInDeps(old_key string, new_key string) string {
  if new_key != "" && new_key != old_key && !taskmap.HasTask(new_key) {
    for dep_key, task := range taskmap {
      index := slices.Index(task.Deps, old_key)
      if index != -1 {
        task.Deps = slices.Replace(task.Deps, index, index+1, new_key)
        taskmap[dep_key] = task
      }
    }
    delete(taskmap, old_key)
    return new_key
  } else {
    return old_key
  }
}

func (taskmap TaskMap) Write(path string) error {
  content, err := yaml.Marshal(taskmap)
  if err != nil {
    slog.Error("Failed converting the tasks to yaml.")
    return err
  }
  if err := os.WriteFile(path, content, 0644); err != nil {
    slog.Error("Failed writing the tasks to file.", "path", path)
    return err
  }
  slog.Info("Wrote to file", "path", path)
  return nil
}

func (taskmap TaskMap) DryWrite(path string) {
  _, err := yaml.Marshal(taskmap)
  if err != nil {
    slog.Error("Failed converting the tasks to yaml.")
    os.Exit(1)
  }
  slog.Info("(dry) Wrote to file", "path", path)
}

func LoadTaskMap(path string) TaskMap {
  slog.Info("Task file loaded.", "path", path)
  taskmap := make(TaskMap)
  content, _ := os.ReadFile(path)
  ParseYaml(taskmap, content)
  return taskmap
}
