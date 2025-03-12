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
  todo := Todo{};
  ParseYaml(&todo, []byte(DATA));
  expected := Todo {
    Task: "buy groceries",
    Deps: []string{"2"},
    Done: true,
  };
  assert.Equal(t, todo, expected)
}

func TestIsDone(t *testing.T) {
  todo := Todo{};
  ParseYaml(&todo, []byte(DATA));
  assert.True(t, todo.IsDone(nil))
  todo.Done = false
  assert.False(t, todo.IsDone(nil))
}

func ExampleTodo_PrintMarkdown() {
  todo := Todo{};
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
