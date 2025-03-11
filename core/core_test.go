package core

import (
  "testing"
  "reflect"
)

func assertEq(t *testing.T, actual any, expected any, format string) {
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf(format+" (actual) != "+format+" (expected)", actual, expected)
  }
}

func TestParse(t *testing.T) {
  data := `
task: buy groceries
priority: 1
deps: [2]
done: true
`
  todo := Todo{};
  todo.YamlParse(data);
  expected := Todo {
    Task: "buy groceries",
    Priority: 1,
    Deps: []int{2},
    Done: true,
  };
  assertEq(t, todo, expected, "%v")
}

const GROCERIES = `0:
  task: buy groceries
  priority: 1
  deps: [2, 3]
  done: true
1:
  task: buy milk
2:
  task: buy bread
  priority: 12
`

func TestParseMap(t *testing.T) {
  todomap := make(TodoMap)
  todomap.YamlParseMap(GROCERIES)
  expected_tasks := []string{"buy groceries", "buy milk", "buy bread"};
  for i:=0; i < 3; i++  {
    assertEq(t, todomap[i].Task, expected_tasks[i], "%q")
  }

  expected_priorities := []int{1,0,12};
  for i:=0; i < 3; i++  {
    assertEq(t, todomap[i].Priority, expected_priorities[i], "%d")
  }
}

func TestNextId(t *testing.T) {
  todomap := make(TodoMap)
  todomap.YamlParseMap(GROCERIES)
  for i:=0; i < 20; i++ {
    assertEq(t, todomap.NextId(), 3+i, "%d")
    todomap[3+i] = Todo{};
  }
  todomap[24] = Todo{};
  assertEq(t, todomap.NextId(), 23, "%d")
  todomap[23] = Todo{};
  assertEq(t, todomap.NextId(), 25, "%d")
}
