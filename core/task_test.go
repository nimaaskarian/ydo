package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const DATA = `
task: buy groceries
deps: [2]
done: true
`

func TestParseYaml(t *testing.T) {
  task := Task{};
  ParseYaml(&task, []byte(DATA));
  expected := Task {
    Task: TemplateBase{template: "buy groceries"},
    Deps: []string{"2"},
    Done: true,
  };
  assert.Equal(t, task, expected)
}

func TestIsDone(t *testing.T) {
  task := Task{};
  ParseYaml(&task, []byte(DATA));
  assert.True(t, task.IsDone(nil, time.Now()))
  task.Done = false
  assert.False(t, task.IsDone(nil, time.Now()))
}

func TestIsNotDone(t *testing.T) {
  task := Task{};
  ParseYaml(&task, []byte(DATA));
  assert.False(t, task.IsNotDone(nil, time.Now()))
  task.Done = false
  assert.True(t, task.IsNotDone(nil, time.Now()))
}


func ExampleTask_PrintMarkdown() {
  task := &Task{};
  ParseYaml(task, []byte(DATA));
  task.Deps = []string{};
	config := MarkdownConfig{Indent: 3, Now: time.Now()}
  config.Init()
  task.PrintMarkdown(nil, 0, nil, "", &config)
  task.Done = false;
  task.DoneAt = time.Now().Add(-time.Hour*24)
  task.PrintMarkdown(nil, 0, nil, "", &config)
  task.Done = true;
  task.PrintMarkdown(nil, 0, nil, "", &config)
  task.Due = TemplateDate{ date: time.Now().Add(-time.Hour*24*2) }
  config_done := config
  config_done.Filter = (*Task).IsNotDone
  task.PrintMarkdown(nil, 0, nil, "", &config_done)
  config_limit := config
  config_limit.Limit = 1
  task.Deps = []string{"2"};
  task.PrintMarkdown(nil, 0, map[string]bool{}, "", &config_limit)
  task.Undo(nil, time.Now())
  task.Due = TemplateDate{ date: time.Now().AddDate(10000, 0, 0) }
  task.PrintMarkdown(nil, 0, nil, "", &config)
  taskmap := TaskMap{}
  taskmap["2"] = &Task{}
  task.PrintMarkdown(taskmap, 0, nil, "", &config)
  // Output:
  // - [x] buy groceries
  // - [ ] buy groceries
  // - [x] buy groceries (1d ago)
  // - [x] buy groceries (1d ago, 1d overdue)
  // - [ ] buy groceries (10000y)
  // - [ ] buy groceries (10000y)
  //    - [ ] 2:
}

func TestInit(t *testing.T) {
  task := &Task{};
  ParseYaml(task, []byte(TEMPLATE));
  assert.Equal(t, "do something till {{ .Due }}", task.Task)
  old_description := task.Description
  task.ResolveTemplates()
  assert.Equal(t, "do something till 2025-12-14 00:00:00 +0330 +0330", task.Task)
  assert.Equal(t, old_description, task.Description)
}

const TEMPLATE = `task: do something till {{ .Due }}
due: 2025-12-14T00:00:00+03:30
description: |-
  ok some description might say stuff. some might not.
  idk bro i have no idea.
`
