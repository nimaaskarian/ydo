package core

import (
  "testing"
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
  task.PrintMarkdown(nil, 1, nil, nil)
  task.Done = false;
  task.PrintMarkdown(nil, 1, nil, nil)
  task.Deps = []string{"2"};
  task.Done = true;
  task.PrintMarkdown(nil, 1, nil, nil)
  // Output:
  // - [x] buy groceries
  // - [ ] buy groceries
  // - [x] buy groceries
  //    - [ ] 
}
