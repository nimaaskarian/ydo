package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nimaaskarian/ydo/utils"
)

type Task struct {
  Task string                   `yaml:",omitempty"`
  Description string            `yaml:",omitempty"`
  Deps []string                 `yaml:",omitempty,flow"`
  Done bool                     `yaml:",omitempty"`
  AutoComplete bool             `yaml:"auto-complete,omitempty"`
  CreatedAt time.Time           `yaml:"created-at,omitempty"`
  Due time.Time                 `yaml:",omitempty"`
  DoneAt time.Time              `yaml:"done-at,omitempty"`
  OldDoneAtList []time.Time     `yaml:"old-done-at-list,omitempty"`
  Recur string                  `yaml:",omitempty"`
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
  if t, err := utils.ParseDuration(task.Recur, task.DoneAt); err == nil && !t.IsZero() && time.Now().After(t) {
    task.Done = false
  } 
  return task.Done
}

func (task *Task) Do() {
  if !task.Done && !task.AutoComplete {
    task.Done = true
    if task.Recur != "" && !task.DoneAt.IsZero() {
      task.OldDoneAtList = append(task.OldDoneAtList, task.DoneAt)
    }
    task.DoneAt = time.Now()
  }
}

func (task *Task) Undo() {
  if task.Done && !task.AutoComplete {
    task.Done = false
    task.DoneAt = time.Time{}
  }
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

// copy a task, delete() the key, run this.
// runs over dependencies listed within the task itself.
func (task Task) CascadeOrphanDeps(taskmap TaskMap) {
  for _, dep := range task.Deps {
    if !taskmap.HasKeyInDeps(dep) {
      delete(taskmap, dep)
    }
  }
}

func (task Task) PrintMarkdown(taskmap TaskMap, depth uint, seen_keys map[string]bool, key string, config *MarkdownConfig) (count int) {
  if config.Filter != nil && !config.Filter(taskmap[key], taskmap) {
    return 0
  }
  if config.Limit != 0 && len(seen_keys) >= config.Limit {
    return 1
  }
  for range depth*config.Indent {
    fmt.Print(" ")
  }
  var print_key string
  if key != "" {
    print_key = key+": "
  }
  var recur string
  if task.IsDone(taskmap) {
    done_at := task.FindDoneAt(taskmap)
    if !done_at.IsZero() {
      overdue := ""
      if !task.Due.IsZero() && done_at.After(task.Due) {
        overdue += ", " + utils.FormatDuration(done_at.Sub(task.Due)) + " overdue"
      }
      if task.Recur != "" {
        recur = ", each "+task.Recur
      }
      fmt.Printf("- [x] %s%s (%s ago%s%s)\n", print_key,task.Task, utils.FormatDuration(time.Now().Sub(done_at)), overdue, recur)
    } else {
      if task.Recur != "" {
        recur = " (each "+task.Recur + ")"
      }
      fmt.Printf("- [x] %s%s%s\n", print_key,task.Task, recur)
    }
  } else {
    if task.Recur != "" {
      recur = " (each "+task.Recur + ")"
    }
    due_print := ""
    if !task.Due.IsZero() {
      now := time.Now()
      diff := task.Due.Sub(now)
      due_print = " ("
      if task.Due.Add(-diff).Compare(now) != 0 {
        due_print += strconv.Itoa(task.Due.Year() - now.Year()) + "y"
      } else {
        if diff < 0 {
          due_print = " (-"
          diff = - diff
        }
        due_print += utils.FormatDuration(diff)
      }
      due_print += ")"
    }
    fmt.Printf("- [ ] %s%s%s%s\n", print_key,task.Task, due_print, recur)
  }
  if config.Description && task.Description != "" {
    for line := range strings.Lines(task.Description) {
      for range (depth+1)*config.Indent {
        fmt.Print(" ")
      }
      fmt.Print(line)
    }
    fmt.Println()
  }
  if seen_keys != nil  {
    if value, ok := seen_keys[key]; ok && value {
      return 0
    }
    seen_keys[key] = true
  }
  for _, key := range task.Deps {
    count += taskmap[key].PrintMarkdown(taskmap, depth+1, seen_keys, key, config)
  }
  return 1+count
}
