package core

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
  data := `
task: buy groceries
deps: [2]
done: true
`
  todo := Todo{};
  ParseYaml(&todo, []byte(data));
  expected := Todo {
    Task: "buy groceries",
    Deps: []string{"2"},
    Done: true,
  };
  assert.Equal(t, todo, expected)
}

