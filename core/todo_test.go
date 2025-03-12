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
  todo := Task{};
  ParseYaml(&todo, []byte(DATA));
  expected := Task {
    Task: "buy groceries",
    Deps: []string{"2"},
    Done: true,
  };
  assert.Equal(t, todo, expected)
}

func TestIsDone(t *testing.T) {
  todo := Task{};
  ParseYaml(&todo, []byte(DATA));
  assert.True(t, todo.IsDone(nil))
  todo.Done = false
  assert.False(t, todo.IsDone(nil))
}

func ExampleTask_PrintMarkdown() {
  todo := Task{};
  ParseYaml(&todo, []byte(DATA));
  todo.Deps = []string{};
  todo.PrintMarkdown(nil, 1)
  todo.Done = false;
  todo.PrintMarkdown(nil, 1)
  todo.Deps = []string{"2"};
  todo.Done = true;
  todo.PrintMarkdown(nil, 1)
  // Output:
  // - [x] buy groceries
  // - [ ] buy groceries
  // - [x] buy groceries
  //    - [ ] 
}
