package core

import (
  "testing"
  "strconv"
  "github.com/stretchr/testify/assert"
)

const GROCERIES = `t1:
  task: buy groceries
  deps: [t3, t2]
  auto-complete: true
t2:
  task: buy milk
t3:
  task: buy bread
`

func TestTaskMapParseYaml(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES));
  expected_tasks := []string{"buy groceries", "buy milk", "buy bread"};
  for i:=range 3 {
    assert.Equal(t, expected_tasks[i], taskmap["t"+strconv.Itoa(i+1)].Task)
  }
}

func TestNextKey(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES));
  for i:=range 20 {
    assert.Equal(t, taskmap.NextKey(), "t"+strconv.Itoa(4+i))
    taskmap["t"+strconv.Itoa(4+i)] = Task{};
  }
  taskmap["t25"] = Task{};
  assert.Equal(t, "t24", taskmap.NextKey())
  taskmap["t24"] = Task{};
  assert.Equal(t, "t26", taskmap.NextKey())
}

func TestDo(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  assert.False(t, taskmap["t2"].IsDone(taskmap))
  taskmap.Do("t2")
  assert.True(t, taskmap["t2"].IsDone(taskmap))
  assert.False(t, taskmap["t3"].IsDone(taskmap))
  taskmap.Do("t3")
  assert.True(t, taskmap["t3"].IsDone(taskmap))
}

func TestUndo(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  taskmap.Do("t2")
  assert.True(t, taskmap["t2"].IsDone(taskmap))
  taskmap.Undo("t2")
  taskmap.Do("t3")
  assert.True(t, taskmap["t3"].IsDone(taskmap))
  taskmap.Undo("t3")
  assert.False(t, taskmap["t3"].IsDone(taskmap))
}

func TestDepIsDone(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  assert.False(t, taskmap["t1"].IsDone(taskmap))
  taskmap.Do("t2")
  assert.False(t, taskmap["t1"].IsDone(taskmap))
  taskmap.Do("t3")
  assert.True(t, taskmap["t1"].IsDone(taskmap))
}

func ExamplePrintYaml() {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  PrintYaml(taskmap)
  taskmap = make(TaskMap)
  ParseYaml(taskmap, []byte(HOMEWORKS))
  PrintYaml(taskmap)
  // Output:
  // t1:
  //     task: buy groceries
  //     deps: [t3, t2]
  //     auto-complete: true
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
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  taskmap.PrintKeys()
  taskmap = make(TaskMap)
  ParseYaml(taskmap, []byte(HOMEWORKS))
  taskmap.PrintKeys()
  // Unordered output:
  // t1
  // t2
  // t3
  // homework
  // study
  // project
}

func TestHasTask(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  assert.False(t, taskmap.HasTask("homework"))
  assert.False(t, taskmap.HasTask("project"))
  assert.False(t, taskmap.HasTask(""))
  assert.True(t, taskmap.HasTask("t1"))
  assert.True(t, taskmap.HasTask("t2"))
  assert.True(t, taskmap.HasTask("t3"))
}
