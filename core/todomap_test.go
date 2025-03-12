package core

import (
  "testing"
  "strconv"
  "github.com/stretchr/testify/assert"
)

const GROCERIES = `t0:
  task: buy groceries
  deps: [t2, t1]
  donedeps: true
t1:
  task: buy milk
t2:
  task: buy bread
`

func TestTodoMapParseYaml(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES));
  expected_tasks := []string{"buy groceries", "buy milk", "buy bread"};
  for i:=range 3 {
    assert.Equal(t, expected_tasks[i], todomap["t"+strconv.Itoa(i)].Task)
  }
}

func TestNextKey(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES));
  for i:=range 20 {
    assert.Equal(t, todomap.NextKey(), "t"+strconv.Itoa(3+i))
    todomap["t"+strconv.Itoa(3+i)] = Todo{};
  }
  todomap["t24"] = Todo{};
  assert.Equal(t, "t23", todomap.NextKey())
  todomap["t23"] = Todo{};
  assert.Equal(t, "t25", todomap.NextKey())
}

func TestDo(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.False(t, todomap["t1"].IsDone(todomap))
  todomap.Do("t1")
  assert.True(t, todomap["t1"].IsDone(todomap))
  assert.False(t, todomap["t2"].IsDone(todomap))
  todomap.Do("t2")
  assert.True(t, todomap["t2"].IsDone(todomap))
}

func TestDepIsDone(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.False(t, todomap["t0"].IsDone(todomap))
  todomap.Do("t1")
  assert.False(t, todomap["t0"].IsDone(todomap))
  todomap.Do("t2")
  assert.True(t, todomap["t0"].IsDone(todomap))
}

func ExampleTodoMap_PrintYaml() {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  todomap.PrintYaml()
  todomap = make(TodoMap)
  ParseYaml(todomap, []byte(HOMEWORKS))
  todomap.PrintYaml()
  // Output:
  // t0:
  //     task: buy groceries
  //     deps: [t2, t1]
  //     donedeps: true
  // t1:
  //     task: buy milk
  // t2:
  //     task: buy bread
  // homework:
  //     task: blah blah homework
  //     deps: [study, project]
  // project:
  //     task: blah blah project
  // study:
  //     task: study blah blah
}

const HOMEWORKS = `homework:
    task: blah blah homework
    deps: [study, project]
study:
    task: study blah blah
project:
    task: blah blah project
`

func ExampleTodoMap_PrintKeys() {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  todomap.PrintKeys()
  todomap = make(TodoMap)
  ParseYaml(todomap, []byte(HOMEWORKS))
  todomap.PrintKeys()
  // Unordered output:
  // t0
  // t1
  // t2
  // homework
  // study
  // project
}
