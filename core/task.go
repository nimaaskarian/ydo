package core

import (
	"fmt"
	"time"

	"github.com/nimaaskarian/ydo/utils"
)

type Task struct {
  Task string           `yaml:",omitempty"`
  Deps []string         `yaml:",omitempty,flow"`
  Done bool             `yaml:",omitempty"`
  AutoComplete bool     `yaml:"auto-complete,omitempty"`
  CreatedAt time.Time   `yaml:"created-at,omitempty"`
  Due time.Time         `yaml:",omitempty"`
  DoneAt time.Time      `yaml:"done-at,omitempty"`
}

func (task Task) IsDone(taskmap TaskMap) bool {
  if task.AutoComplete {
    for _,key := range task.Deps {
      if !taskmap[key].IsDone(taskmap) {
        return false
      }
    }
    return true
  }
  return task.Done
}

func (task *Task) Do() {
  task.Done = true
  task.DoneAt = time.Now()
}

func (task *Task) Undo() {
  task.Done = false
  task.DoneAt = time.Time{}
}

func (task *Task) AddDep(tm TaskMap, key string) error {
  if _, err := tm.GetTask(key); err != nil {
    return err
  }
  task.Deps = append(task.Deps, key)
  return nil
}

func (task Task) FindDoneAt(taskmap TaskMap) time.Time {
  if task.AutoComplete {
    max_doneat := time.Time{}
    for _,key := range task.Deps {
      task := taskmap[key]
      doneat := task.FindDoneAt(taskmap)
      if doneat.After(max_doneat) {
        max_doneat = doneat
      }
      return max_doneat
    }
  }
  return task.DoneAt
}

func (task Task) IsNotDone(taskmap TaskMap) bool {
  return !task.IsDone(taskmap)
}

func (task Task) PrintMarkdown(taskmap TaskMap, depth uint, seen_keys map[string]bool, key string,filtered *int, config *MarkdownConfig) {
  if config.Filter != nil && !config.Filter(taskmap[key], taskmap) {
    if filtered != nil {
      *filtered += 1
    }
    return
  }
  if config.Limit != 0 && len(seen_keys) >= config.Limit {
    return
  }
  for range depth*config.Indent {
    fmt.Print(" ")
  }
  var print_key string
  if key != "" {
    print_key = key+": "
  }
  if task.IsDone(taskmap) {
    done_at := task.FindDoneAt(taskmap)
    if !done_at.IsZero() {
      overdue := ""
      if !task.Due.IsZero() && done_at.After(task.Due) {
        overdue += ", " + utils.FormatDuration(done_at.Sub(task.Due)) + " overdue"
      }
      fmt.Fprintf(config.file, "- [x] %s%s (%s ago%s)\n", print_key,task.Task, utils.FormatDuration(time.Now().Sub(done_at)), overdue)
    } else {
      fmt.Fprintf(config.file, "- [x] %s%s\n", print_key,task.Task)
    }
  } else {
    due_print := ""
    if !task.Due.IsZero() {
      diff := task.Due.Sub(time.Now())
      due_print = " ("
      if diff < 0 {
        due_print = " (-"
        diff = - diff
      }
      due_print += utils.FormatDuration(diff) + ")"
    }
    fmt.Fprintf(config.file, "- [ ] %s%s%s\n", print_key,task.Task, due_print)
  }
  if seen_keys != nil  {
    if value, ok := seen_keys[key]; ok && value {
      return
    }
    seen_keys[key] = true
  }
  for _, key := range task.Deps {
    taskmap[key].PrintMarkdown(taskmap, depth+1, seen_keys, key,filtered, config)
  }
}
