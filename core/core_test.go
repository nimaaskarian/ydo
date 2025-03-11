package core

import (
  "testing"
  "reflect"
  "strconv"
)

func assertEq(t *testing.T, actual any, expected any, format string) {
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf(format+" (actual) != "+format+" (expected)", actual, expected)
  }
}

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
  assertEq(t, todo, expected, "%v")
}

const GROCERIES = `0:
  task: buy groceries
  deps: [2, 3]
  done: true
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
    assertEq(t, todomap[strconv.Itoa(i)].Task, expected_tasks[i], "%q")
  }
}

func TestNextId(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES));
  for i:=range 20 {
    assertEq(t, todomap.NextId(), strconv.Itoa(3+i), "%d")
    todomap[strconv.Itoa(3+i)] = Todo{};
  }
  todomap["24"] = Todo{};
  assertEq(t, todomap.NextId(), "23", "%d")
  todomap["23"] = Todo{};
  assertEq(t, todomap.NextId(), "25", "%d")
}
