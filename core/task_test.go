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

func ExampleTask_PrintMarkdown() {
  task := Task{};
  ParseYaml(&task, []byte(DATA));
  task.Deps = []string{};
  task.PrintMarkdown(nil, 0, nil, "", nil, 3)
  task.Done = false;
  task.DoneAt = time.Now().Add(-time.Hour*24)
  task.PrintMarkdown(nil, 0, nil, "", nil, 3)
  task.Done = true;
  task.PrintMarkdown(nil, 0, nil, "", nil, 3)
  task.Due = time.Now().Add(-time.Hour*24*2)
  task.Deps = []string{"2"};
  task.PrintMarkdown(nil, 0, nil, "", nil, 3)
  // Output:
  // - [x] buy groceries
  // - [ ] buy groceries
  // - [x] buy groceries (1d ago)
  // - [x] buy groceries (1d ago, 1d overdue)
  //    - [ ] 2:
}
