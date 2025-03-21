package core

import (
	"os"
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
    Task: "buy groceries",
    Deps: []string{"2"},
    Done: true,
  };
  assert.Equal(t, task, expected)
}

func TestIsDone(t *testing.T) {
  task := Task{};
  ParseYaml(&task, []byte(DATA));
  assert.True(t, task.IsDone(nil))
  task.Done = false
  assert.False(t, task.IsDone(nil))
}

func TestIsNotDone(t *testing.T) {
  task := Task{};
  ParseYaml(&task, []byte(DATA));
  assert.False(t, task.IsNotDone(nil))
  task.Done = false
  assert.True(t, task.IsNotDone(nil))
}


func ExampleTask_PrintMarkdown() {
  task := Task{};
  ParseYaml(&task, []byte(DATA));
  task.Deps = []string{};
  config := MarkdownConfig{Indent: 3, file: os.Stdout}
  var count uint
  task.PrintMarkdown(nil, 0, nil, "", &count, &config)
  task.Done = false;
  task.DoneAt = time.Now().Add(-time.Hour*24)
  task.PrintMarkdown(nil, 0, nil, "", &count, &config)
  task.Done = true;
  task.PrintMarkdown(nil, 0, nil, "", &count, &config)
  task.Due = time.Now().Add(-time.Hour*24*2)
  config_done := config
  config_done.Filter = Task.IsDone
  task.PrintMarkdown(nil, 0, nil, "", &count, &config_done)
  config_limit := config
  config_limit.Limit = 1
  task.Deps = []string{"2"};
  count = 0
  task.PrintMarkdown(nil, 0, nil, "", &count, &config_limit)
  task.PrintMarkdown(nil, 0, nil, "", &count, &config)
  // Output:
  // - [x] buy groceries
  // - [ ] buy groceries
  // - [x] buy groceries (1d ago)
  // - [x] buy groceries (1d ago, 1d overdue)
  // - [x] buy groceries (1d ago, 1d overdue)
  //    - [ ] 2:
}
