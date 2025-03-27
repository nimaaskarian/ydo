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
	"time"

	"github.com/nimaaskarian/ydo/utils"
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

func (taskmap TaskMap) Delete(key string, cascade bool) error {
  task, err := taskmap.GetTask(key)
  if err != nil {
    return err
  }
  delete(taskmap, key)
  taskmap.WipeDependenciesToKey(key)
  if cascade {
    task.CascadeOrphanDeps(taskmap)
  }
  return nil
}

func (taskmap TaskMap) WipeDependenciesToKey(key string) error {
  for each_key, task := range taskmap {
    index := slices.Index(task.Deps, key)
    if index != -1 {
      task.Deps = slices.Delete(task.Deps, index, index+1)
      taskmap[each_key] = task
    }
  }
  return nil
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
  task.Do(taskmap)
  taskmap[key] = task
  slog.Info("Completed task","key" ,key)
  return nil
}

func (tm TaskMap) AddDep(key string, dep string) (Task, error) {
  task, err := tm.GetTask(key)
  if err != nil {
    return Task{}, err
  } else {
    _, err := tm.GetTask(dep)
    if err != nil {
      return Task{}, err
    }
  }
  if !slices.Contains(task.Deps, dep) {
    task.Deps = append(task.Deps, dep)
  }
  return task, nil
}

func (taskmap TaskMap) Undo(key string) error {
  task, err := taskmap.GetTask(key)
  if err != nil{
    return err
  }
  task.Undo(taskmap)
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
  Indent uint             `yaml:",omitempty"`
  Mode string             `yaml:",omitempty"`
  Description bool        `yaml:",omitempty"`
  Limit int               `yaml:",omitempty"`
  Filter MarkdownFilter
}

func (taskmap TaskMap) SortedKeys() []string {
  keys := utils.Keys(taskmap)
  slices.SortFunc(keys, func(k1, k2 string) int {
    t1, t2 := taskmap[k1], taskmap[k2]
    due_zero := 0
    if t1.Due.IsZero() && !t2.Due.IsZero(){
      due_zero = 1
    }  else if !t1.Due.IsZero() && t2.Due.IsZero() {
      due_zero = -1
    } else {
      due_zero = t1.Due.Compare(t2.Due)
    }
    return 2*due_zero+t1.CreatedAt.Compare(t2.CreatedAt)
  })
  return keys
}

func (taskmap TaskMap) PrintMarkdown(config *MarkdownConfig) error {
  if len(taskmap) == 0 {
    return errors.New("No tasks found")
  }
  keys := taskmap.SortedKeys()
  seen_keys := make(map[string]bool, len(taskmap))
  count := 0
  for _,key := range keys {
    if value, ok := seen_keys[key]; !ok || !value {
      count += taskmap[key].PrintMarkdown(taskmap, 0, seen_keys, key, config)
    }
  }
  shown := len(seen_keys)
  if count > shown {
    fmt.Printf("%d tasks, %d shown\n", count, shown)
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

func (taskmap TaskMap) HasKeyInDeps(key string) bool {
  for _, task := range taskmap {
    if slices.Contains(task.Deps, key) {
      return true
    }
  }
  return false
}

func as_days(d time.Duration) int {
  return int(d.Hours())/24
}

func addToIndexIfKeyOk(m map[time.Time][2]int, start, end time.Time, index int, key time.Time) bool {
  if key.After(end) {
    return true
  }
  if key.Before(start) {
    return false
  }
  arr := m[key]
  arr[index] += 1
  m[key] = arr
  return false
}

// start and end are included
func (taskmap TaskMap) TrackDisciplineDaily(start, end time.Time) []float64 {
  end = utils.NaiveDate(end)
  start = utils.NaiveDate(start)
  days := as_days(end.Sub(start))
  total_done_map := make(map[time.Time][2]int)
  for i := range days+1 {
    d := utils.NaiveDate(start).AddDate(0, 0, i)
    total_done_map[d] = [2]int{0, 0}
  }

  for _, task := range taskmap {
    // ignore legeacy tasks
    if task.IsDone(taskmap) && task.DoneAt.IsZero() {
      continue
    }
    created_date := utils.NaiveDate(task.CreatedAt)
    if task.Recur == "" {
      if !task.DoneAt.IsZero() {
        for i := range as_days(task.DoneAt.Sub(created_date))+1 {
          exists_date := created_date.AddDate(0, 0, i)
          if addToIndexIfKeyOk(total_done_map, start, end, 0, exists_date) {
            break
          }
        }
        done_date := utils.NaiveDate(task.DoneAt)
        arr, ok := total_done_map[done_date]
        if ok {
          arr[1] += 1
          total_done_map[done_date] = arr
        }
      } else {
        for i := range as_days(end.Sub(created_date))+1 {
          exists_date := created_date.AddDate(0, 0, i)
          if addToIndexIfKeyOk(total_done_map, start, end, 0, exists_date) {
            break
          }
        }
      }
    } else {
      done_dates := append(task.OldDoneAtList, task.DoneAt)
      created_dates := make([]time.Time, 1, len(done_dates))
      created_dates[0] = created_date
      for _, done_at := range task.OldDoneAtList {
        t, _ := utils.ParseDuration(task.Recur, utils.NaiveDate(done_at))
        created_dates = append(created_dates, t)
      }
      for i, created := range created_dates {
        upper_bound := done_dates[i]
        if done_dates[i].IsZero() {
          upper_bound = end
        }
        for j := range as_days(upper_bound.Sub(created))+1 {
          exists_date := created.AddDate(0, 0, j)
          if addToIndexIfKeyOk(total_done_map, start, end, 0, exists_date) {
            break
          }
        }
      }
      last_created, _ := utils.ParseDuration(task.Recur, utils.NaiveDate(task.DoneAt))
      for i := range as_days(end.Sub(last_created))+1 {
        exists_date := last_created.AddDate(0, 0, i)
        if addToIndexIfKeyOk(total_done_map, start, end, 0, exists_date) {
          break
        }
      }
      for _, done_at := range done_dates {
        if addToIndexIfKeyOk(total_done_map, start, end, 1, utils.NaiveDate(done_at)) {
          break
        }
      }
    }
  }
  discipline_arr := make([]float64, 0, len(total_done_map))
  sorted_keys := utils.Keys(total_done_map)
  slices.SortFunc(sorted_keys, time.Time.Compare)
  for _,key := range sorted_keys {
    item := total_done_map[key]
    discipline := float64(item[1])/float64(item[0])
    discipline_arr = append(discipline_arr, discipline)
  }
  return discipline_arr
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

func (taskmap TaskMap) DryWrite(path string) error {
  _, err := yaml.Marshal(taskmap)
  if err != nil {
    slog.Error("Failed converting the tasks to yaml.")
    return err
  }
  slog.Info("(dry) Wrote to file", "path", path)
  return nil
}

func LoadTaskMap(path string) TaskMap {
  slog.Info("Task file loaded.", "path", path)
  taskmap := TaskMap{}
  content, _ := os.ReadFile(path)
  ParseYaml(taskmap, content)
  return taskmap
}
