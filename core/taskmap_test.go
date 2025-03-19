package core

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const GROCERIES = `t1:
  task: buy groceries
  deps: [t3, t2]
  auto-complete: true
t2:
  task: buy milk
  done: true
  done-at: 2025-03-14T00:00:00+03:30
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
    assert.Equal(t, taskmap.NextKey(""), "t"+strconv.Itoa(4+i))
    taskmap["t"+strconv.Itoa(4+i)] = Task{};
  }
  taskmap["t25"] = Task{};
  assert.Equal(t, "t24", taskmap.NextKey(""))
  taskmap["t24"] = Task{};
  assert.Equal(t, "t26", taskmap.NextKey(""))
}

func TestDo(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
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
  assert.True(t, taskmap["t2"].DoneAt.IsZero())
  taskmap.Do("t3")
  assert.True(t, taskmap["t3"].IsDone(taskmap))
  taskmap.Undo("t3")
  assert.False(t, taskmap["t3"].IsDone(taskmap))
  assert.True(t, taskmap["t3"].DoneAt.IsZero())
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
  //     done: true
  //     done-at: 2025-03-14T00:00:00+03:30
  // t3:
  //     task: buy bread
  // homework:
  //     task: do uni practice
  //     deps: [study, project]
  // milk:
  //     task: buy some milk
  // project:
  //     task: do the hobby project
  // study:
  //     task: study for the uni exam
}

const HOMEWORKS = `homework:
    task: do uni practice
    deps: [study, project]
study:
    task: study for the uni exam
project:
    task: do the hobby project
milk:
    task: buy some milk
`

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

func  TestTfidfNextKey(t *testing.T) {
  config := TfidfConfig {Enabled: true};
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(HOMEWORKS))
  key := taskmap.TfidfNextKey("buy some laptop for uni", config, "")
  assert.Equal(t, "laptop", key)
  config = TfidfConfig {Enabled: false};
  key = taskmap.TfidfNextKey("buy some laptop for uni", config, "")
  assert.Equal(t, "t1", key)
}

func TestReplaceKeyInDeps(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  fmt.Println(taskmap["t1"].Deps)
  taskmap.ReplaceKeyInDeps("t2", "milk")
  expected := []string{"t3", "milk"}
  assert.Equal(t, expected, taskmap["t1"].Deps)
}

func TestWriteAndLoad(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  tm.Write("test")
  tm2 := LoadTaskMap("test")
  assert.Equal(t, tm, tm2)
  os.Remove("test")
}

func TestMustHave(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  if os.Getenv("BE_CRASHER") == "1" {
    tm.MustHave("t8")
    return
  }
  cmd := exec.Command(os.Args[0], "-test.run=TestMustHave")
  cmd.Env = append(os.Environ(), "BE_CRASHER=1")
  err := cmd.Run()
  if e, ok := err.(*exec.ExitError); ok && !e.Success() {
    return
  }
  t.Fatal("Crasher didn't crash")
}

func TestMustNotHave(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  if os.Getenv("BE_CRASHER") == "1" {
    tm.MustNotHave("t3")
    return
  }
  cmd := exec.Command(os.Args[0], "-test.run=TestMustNotHave")
  cmd.Env = append(os.Environ(), "BE_CRASHER=1")
  err := cmd.Run()
  if e, ok := err.(*exec.ExitError); ok && !e.Success() {
    return
  }
  t.Fatal("Crasher didn't crash.", "err:", err)
}

func TestFindDoneAt(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  tm.Do("t3")
  assert.Equal(t, tm["t3"].DoneAt, tm["t1"].FindDoneAt(tm))
}
