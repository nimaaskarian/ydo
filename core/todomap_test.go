package core

import (
  "testing"
  "strconv"
  "github.com/stretchr/testify/assert"
)

const GROCERIES = `0:
  task: buy groceries
  deps: [2, 1]
  donedeps: true
1:
  task: buy milk
2:
  task: buy bread
`

func TestTodoMapParseYaml(t *testing.T) {
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
    assert.Equal(t, todomap.NextKey(), strconv.Itoa(3+i))
    todomap[strconv.Itoa(3+i)] = Todo{};
  }
  todomap["24"] = Todo{};
  assert.Equal(t, todomap.NextKey(), "23")
  todomap["23"] = Todo{};
  assert.Equal(t, todomap.NextKey(), "25")
}

func TestDo(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.False(t, todomap["1"].IsDone(todomap))
  todomap.Do("1")
  assert.True(t, todomap["1"].IsDone(todomap))
  assert.False(t, todomap["2"].IsDone(todomap))
  todomap.Do("2")
  assert.True(t, todomap["2"].IsDone(todomap))
}

func TestDepIsDone(t *testing.T) {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.False(t, todomap["0"].IsDone(todomap))
  todomap.Do("1")
  assert.False(t, todomap["0"].IsDone(todomap))
  todomap.Do("2")
  assert.True(t, todomap["0"].IsDone(todomap))
}

func ExampleTodoMap_PrintYaml() {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  todomap.PrintYaml()
  // Output:
  // "0":
  //     task: buy groceries
  //     deps:
  //         - "2"
  //         - "1"
  //     donedeps: true
  // "1":
  //     task: buy milk
  // "2":
  //     task: buy bread
}

const HOMEWORKS = `homework:
    task: blah blah homework
    deps: [study, project]
study:
    task: study blah balh
project:
    task: blah balh project
`

func ExampleTodoMap_PrintKeys() {
  todomap := make(TodoMap)
  ParseYaml(todomap, []byte(GROCERIES))
  todomap.PrintKeys()
  todomap = make(TodoMap)
  ParseYaml(todomap, []byte(HOMEWORKS))
  todomap.PrintKeys()
  // Unordered output:
  // 0
  // 1
  // 2
  // homework
  // study
  // project
}
