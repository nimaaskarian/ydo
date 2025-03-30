package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nimaaskarian/ydo/utils"
  "github.com/fatih/color"
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
  DoneAtArchive []time.Time     `yaml:"done-at-archive,omitempty"`
  Recur string                  `yaml:",omitempty"`
  Tags []string                 `yaml:",omitempty"`
}

func (task *Task) IsDone(taskmap TaskMap, now time.Time) bool {
  if task.AutoComplete {
    for _,key := range task.Deps {
      if !taskmap[key].IsDone(taskmap, now) {
        return false
      }
    }
    return true
  }
  if t, err := utils.ParseDuration(task.Recur, task.DoneAt); err == nil && !t.IsZero() && now.After(t) {
    task.Done = false
  } 
  return task.Done && (task.DoneAt.IsZero() || !task.DoneAt.Before(task.DoneAt))
}

func (task *Task) Do(taskmap TaskMap, now time.Time) {
  if !task.IsDone(taskmap, now) && !task.AutoComplete {
    task.Done = true
    if task.Recur != "" && !task.DoneAt.IsZero() {
      task.DoneAtArchive = append(task.DoneAtArchive, task.DoneAt)
    }
    task.DoneAt = now
  }
}

func (task *Task) Undo(taskmap TaskMap, now time.Time) {
  if task.IsDone(taskmap, now) && !task.AutoComplete {
    task.Done = false
    if task.Recur != "" && len(task.DoneAtArchive) > 0 {
      length := len(task.DoneAtArchive)
      task.DoneAt = task.DoneAtArchive[length-1]
      task.DoneAtArchive = task.DoneAtArchive[:length-1]
    } else {
      task.DoneAt = time.Time{}
    }
  }
}

func (task *Task) FindDoneAt(taskmap TaskMap) time.Time {
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

func (task *Task) IsNotDone(taskmap TaskMap, now time.Time) bool {
  return !task.IsDone(taskmap, now)
}

// copy a task, delete() the key, run this.
// runs over dependencies listed within the task itself.
func (task *Task) CascadeOrphanDeps(taskmap TaskMap) {
  for _, dep := range task.Deps {
    if !taskmap.HasKeyInDeps(dep) {
      delete(taskmap, dep)
    }
  }
}

func (task *Task) PrintMarkdown(taskmap TaskMap, depth uint, seen_keys map[string]bool, key string, config *MarkdownConfig) (count int) {
  if task == nil {
    return 0
  }
  if config.Filter != nil && !config.Filter(task, taskmap, config.Now) {
    return 0
  }
	if config.Now.Before(task.CreatedAt) {
		return 0
	}
  if config.Limit != 0 && len(seen_keys) >= config.Limit {
    return 1
  }
  printIndent(depth, config)
  if task.IsDone(taskmap, config.Now) {
    fmt.Print("- [x] ")
    printKey(key, config)
    printDoneTask(task, taskmap, config)
  } else {
    fmt.Print("- [ ] ")
    printKey(key, config)
    printPendingTask(task, config)
  }
  printTags(task, config)
  fmt.Println()
  printDescription(depth, task, config)
  if seen_keys != nil  {
    if _, ok := seen_keys[key]; ok {
      return 0
    }
    seen_keys[key] = true
  }
  for _, dep_key := range task.Deps {
    count += taskmap[dep_key].PrintMarkdown(taskmap, depth+1, seen_keys, dep_key, config)
  }
  return 1+count
}

func printDoneTask(task *Task, taskmap TaskMap, config *MarkdownConfig) {
  var recur string
  done_at := task.FindDoneAt(taskmap)
  if !done_at.IsZero() {
    overdue := ""
    if !task.Due.IsZero() && done_at.After(task.Due) {
      overdue += ", " + utils.FormatDuration(done_at.Sub(task.Due)) + " overdue"
    }
    if task.Recur != "" {
      recur = ", each "+task.Recur
    }
    fmt.Printf("%s (%s ago%s%s)", task.Task, utils.FormatDuration(config.Now.Sub(done_at)), overdue, recur)
  } else {
    if task.Recur != "" {
      recur = " (each "+task.Recur + ")"
    }
    fmt.Printf("%s%s", task.Task, recur)
  }
}

func printPendingTask(task *Task, config *MarkdownConfig) {
  var recur string
  if task.Recur != "" {
    recur = " (each "+task.Recur + ")"
  }
  due_print := ""
  if !task.Due.IsZero() {
    diff := task.Due.Sub(config.Now)
    due_print = " ("
    if task.Due.Add(-diff).Compare(config.Now) != 0 {
      due_print += strconv.Itoa(task.Due.Year() - config.Now.Year()) + "y"
    } else {
      if diff < 0 {
        due_print = " (-"
        diff = - diff
      }
      due_print += utils.FormatDuration(diff)
    }
    due_print += ")"
  }
  fmt.Printf("%s%s%s", task.Task, due_print, recur)
}

func printKey(key string, config *MarkdownConfig) {
  if key != "" {
    if config.HasColor(){
      d := color.New(color.Bold)
      fmt.Printf("%s: ", d.Sprint(key))
    } else {
      fmt.Printf("%s: ", key)
    }
  }
}

func printTags(task * Task, config *MarkdownConfig) {
  for _, tag := range task.Tags {
    if config.HasColor() {
      d := color.New(color.Underline)
      fmt.Print(" ")
      d.Printf("#%s", tag)
    } else {
      fmt.Printf(" #%s", tag)
    }
  }
}

func printIndent(depth uint, config *MarkdownConfig) {
  for range depth*config.Indent {
    fmt.Print(" ")
  }
}

func printDescription(depth uint, task *Task, config *MarkdownConfig) {
  if config.Description && task.Description != "" {
    for line := range strings.Lines(task.Description) {
      printIndent(depth+1, config)
      fmt.Print(line)
    }
    fmt.Println()
  }
}
