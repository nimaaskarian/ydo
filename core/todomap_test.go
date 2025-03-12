package core

import (
  "testing"
  "strconv"
  "github.com/stretchr/testify/assert"
)

const GROCERIES = `t1:
  task: buy groceries
  deps: [t3, t2]
  donedeps: true
t2:
  task: buy milk
t3:
  task: buy bread
`

func TestTodoMapParseYaml(t *testing.T) {
  todomap := make(TaskMap)
  ParseYaml(todomap, []byte(GROCERIES));
  expected_tasks := []string{"buy groceries", "buy milk", "buy bread"};
  for i:=range 3 {
    assert.Equal(t, expected_tasks[i], todomap["t"+strconv.Itoa(i+1)].Task)
  }
}

func TestNextKey(t *testing.T) {
  todomap := make(TaskMap)
  ParseYaml(todomap, []byte(GROCERIES));
  for i:=range 20 {
    assert.Equal(t, todomap.NextKey(), "t"+strconv.Itoa(4+i))
    todomap["t"+strconv.Itoa(4+i)] = Task{};
  }
  todomap["t25"] = Task{};
  assert.Equal(t, "t24", todomap.NextKey())
  todomap["t24"] = Task{};
  assert.Equal(t, "t26", todomap.NextKey())
}

func TestDo(t *testing.T) {
  todomap := make(TaskMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.False(t, todomap["t2"].IsDone(todomap))
  todomap.Do("t2")
  assert.True(t, todomap["t2"].IsDone(todomap))
  assert.False(t, todomap["t3"].IsDone(todomap))
  todomap.Do("t3")
  assert.True(t, todomap["t3"].IsDone(todomap))
}

func TestDepIsDone(t *testing.T) {
  todomap := make(TaskMap)
  ParseYaml(todomap, []byte(GROCERIES))
  assert.False(t, todomap["t1"].IsDone(todomap))
  todomap.Do("t2")
  assert.False(t, todomap["t1"].IsDone(todomap))
  todomap.Do("t3")
  assert.True(t, todomap["t1"].IsDone(todomap))
}

func ExamplePrintYaml() {
  todomap := make(TaskMap)
  ParseYaml(todomap, []byte(GROCERIES))
  PrintYaml(todomap)
  todomap = make(TaskMap)
  ParseYaml(todomap, []byte(HOMEWORKS))
  PrintYaml(todomap)
  // Output:
  // t1:
  //     task: buy groceries
  //     deps: [t3, t2]
  //     donedeps: true
  // t2:
  //     task: buy milk
  // t3:
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

func ExampleTaskMap_PrintKeys() {
  todomap := make(TaskMap)
  ParseYaml(todomap, []byte(GROCERIES))
  todomap.PrintKeys()
  todomap = make(TaskMap)
  ParseYaml(todomap, []byte(HOMEWORKS))
  todomap.PrintKeys()
  // Unordered output:
  // t1
  // t2
  // t3
  // homework
  // study
  // project
}
