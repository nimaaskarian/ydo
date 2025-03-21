package core

import (
	"fmt"
	"os"
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

func TestTaskMapParseYamlFail(t *testing.T) {
  taskmap := make(TaskMap)
  assert.Panics(t, func () { ParseYaml(taskmap, []byte(GROCERIES+"laksjdflj;alsdrandombytesidk")); })
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
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  assert.True(t, tm["t2"].IsDone(tm))
  tm.Undo("t2")
  assert.False(t, tm["t2"].IsDone(tm))
}

func TestUndoDoneAt(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  tm.Undo("t2")
  fmt.Println(tm["t2"].DoneAt)
  assert.True(t, tm["t2"].DoneAt.IsZero())
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
  //     created-at: 2024-01-01T00:00:00+03:30
  // milk:
  //     task: buy some milk
  //     created-at: 2024-01-01T11:30:00+03:30
  // project:
  //     task: do the hobby project
  //     created-at: 2024-01-01T11:00:00+03:30
  // study:
  //     task: study for the uni exam
  //     created-at: 2024-01-01T10:00:00+03:30
}

const HOMEWORKS = `homework:
    task: do uni practice
    deps: [study, project]
    created-at: 2024-01-01T00:00:00+03:30
study:
    task: study for the uni exam
    created-at: 2024-01-01T10:00:00+03:30
project:
    task: do the hobby project
    created-at: 2024-01-01T11:00:00+03:30
milk:
    task: buy some milk
    created-at: 2024-01-01T11:30:00+03:30
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
  tm := make(TaskMap)
  ParseYaml(tm, []byte(HOMEWORKS))
  msg := "buy some laptop for uni"
  key := tm.TfidfNextKey(msg, config, "")
  tm[key] = Task { Task: msg}
  assert.Equal(t, "laptop", key)
  key = tm.TfidfNextKey("buy some milk (fresh)", config, "milk")
  assert.Equal(t, "milk", key)
  config = TfidfConfig {Enabled: false};
  key = tm.TfidfNextKey("buy some laptop for uni", config, "")
  assert.Equal(t, "t1", key)
}

func TestReplaceKeyInDeps(t *testing.T) {
  taskmap := make(TaskMap)
  ParseYaml(taskmap, []byte(GROCERIES))
  fmt.Println(taskmap["t1"].Deps)
  key := taskmap.ReplaceKeyInDeps("t2", "milk")
  expected := []string{"t3", "milk"}
  assert.Equal(t, expected, taskmap["t1"].Deps)
  assert.Equal(t, "milk", key)

  key = taskmap.ReplaceKeyInDeps("t3", "")
  assert.Equal(t, expected, taskmap["t1"].Deps)
  assert.Equal(t, "t3", key)
}

func TestWriteAndLoad(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  tm.Write("test")
  tm2 := LoadTaskMap("test")
  assert.Equal(t, tm, tm2)
  assert.Error(t, tm.Write(""))
  os.Remove("test")
}

func TestFindDoneAt(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  tm.Do("t3")
  assert.Equal(t, tm["t3"].DoneAt, tm["t1"].FindDoneAt(tm))
}

func TestWipeDependenciesToKey(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(GROCERIES))
  tm.WipeDependenciesToKey("t3")
  assert.Equal(t, []string{"t2"}, tm["t1"].Deps)
}

func BenchmarkPrintMarkdown(b *testing.B) {
  tm := LoadTaskMap("../tests/tasks.yaml")
  for b.Loop() {
    tm.PrintMarkdown(&MarkdownConfig{Indent: 4})
  }
}

func TestAddDep(t *testing.T) {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(HOMEWORKS))
  task := tm["homework"]
  assert.Error(t, task.AddDep(tm, "coco"))
  task.AddDep(tm, "milk")
  assert.Equal(t,[]string{"study", "project", "milk"}, task.Deps)
}

func TestEmptyTaskMapMarkdownError(t *testing.T) {
  tm := make(TaskMap)
  assert.Error(t, tm.PrintMarkdown(&MarkdownConfig{Indent: 4}))
}

func ExampleTaskMap_PrintMarkdown() {
  tm := make(TaskMap)
  ParseYaml(tm, []byte(HOMEWORKS))
  tm.Do("study")
  tm.PrintMarkdown(&MarkdownConfig{Indent: 4})
  // Output:
  // - [ ] homework: do uni practice
  //     - [x] study: study for the uni exam (0s ago)
  //     - [ ] project: do the hobby project
  // - [ ] milk: buy some milk
}
