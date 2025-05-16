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
	task := Task{}
	ParseYaml(&task, []byte(DATA))
	expected := Task{
		Task: TemplateBase{Template: "buy groceries"},
		Deps: []string{"2"},
		Done: true,
	}
	assert.Equal(t, task, expected)
}

func TestIsDone(t *testing.T) {
	task := Task{}
	ParseYaml(&task, []byte(DATA))
	assert.True(t, task.IsDone(nil, time.Now()))
	task.Done = false
	assert.False(t, task.IsDone(nil, time.Now()))
}

func TestIsNotDone(t *testing.T) {
	task := Task{}
	ParseYaml(&task, []byte(DATA))
	assert.False(t, task.IsNotDone(nil, time.Now()))
	task.Done = false
	assert.True(t, task.IsNotDone(nil, time.Now()))
}

func ExampleTask_PrintMarkdown() {
	task := &Task{}
	ParseYaml(task, []byte(DATA))
	task.Deps = []string{}
	task.CreatedAt = time.Now()
	config := MarkdownConfig{Indent: 3, Now: time.Now()}
	config.Init()
	task.PrintMarkdown(nil, 0, nil, "", &config, nil)
	task.Done = false
	task.DoneAt = time.Now().Add(-time.Hour * 24)
	task.PrintMarkdown(nil, 0, nil, "", &config, nil)
	task.Done = true
	task.PrintMarkdown(nil, 0, nil, "", &config, nil)
	task.Due = TemplateDate{Date: time.Now().Add(-time.Hour * 24 * 2)}
	task.PrintMarkdown(nil, 0, nil, "", &config, (*Task).IsNotDone)
	config_limit := config
	config_limit.Limit = 1
	task.Deps = []string{"2"}
	task.PrintMarkdown(nil, 0, nil, "", &config_limit, (*Task).IsNotDone)
	task.Due = TemplateDate{Date: time.Now().AddDate(0, 0, -2)}
	task.PrintMarkdown(nil, 0, nil, "", &config, nil)
	task.Undo(nil, time.Now())
	task.Due = TemplateDate{Date: time.Now().AddDate(10000, 0, 0)}
	task.PrintMarkdown(nil, 0, nil, "", &config, (*Task).IsNotDone)
	taskmap := TaskMap{}
	taskmap["2"] = &Task{}
	task.PrintMarkdown(taskmap, 0, nil, "", &config, (*Task).IsNotDone)
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
	task := &Task{}
	ParseYaml(task, []byte(TEMPLATE))
	assert.Equal(t, "do something till {{ .Due }}", task.Task.Template)
	old_description := task.Description.Value()
	task.ResolveTemplates(time.Now())
	assert.Equal(t, "do something till 2025-12-14T00:00:00+03:30", task.Task.Value())
	assert.Equal(t, old_description, task.Description.Value())
}

const TEMPLATE = `task: do something till {{ .Due }}
due: 2025-12-14T00:00:00+03:30
description: |-
  ok some description might say stuff. some might not.
  idk bro i have no idea.
`
