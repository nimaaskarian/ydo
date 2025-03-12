package core

import (
  "testing"
  "strconv"
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

const GROCERIES = `0:
  task: buy groceries
  deps: [2, 1]
  donedeps: true
1:
  task: buy milk
2:
  task: buy bread
`

func TestParseMap(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES));
  expected_tasks := []string{"buy groceries", "buy milk", "buy bread"};
  for i:=range 3 {
    assert.Equal(t, expected_tasks[i], todomap[strconv.Itoa(i)].Task)
  }
}

func TestNextId(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES));
  for i:=range 20 {
    assert.Equal(t, todomap.NextId(), strconv.Itoa(3+i))
    todomap[strconv.Itoa(3+i)] = Todo{};
  }
  todomap["24"] = Todo{};
  assert.Equal(t, todomap.NextId(), "23")
  todomap["23"] = Todo{};
  assert.Equal(t, todomap.NextId(), "25")
}

func TestMapDo(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.Equal(t, false, todomap["1"].IsDone(todomap))
  todomap.Do("1")
  assert.Equal(t, true, todomap["1"].IsDone(todomap))
  assert.Equal(t, false, todomap["2"].IsDone(todomap))
  todomap.Do("2")
  assert.Equal(t, true, todomap["2"].IsDone(todomap))
}

func TestIsDone(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.Equal(t, false, todomap["0"].IsDone(todomap))
  todomap.Do("1")
  assert.Equal(t, false, todomap["0"].IsDone(todomap))
  todomap.Do("2")
  assert.Equal(t, true, todomap["0"].IsDone(todomap))
}
