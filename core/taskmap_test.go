package core

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"
	"time"

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
	ParseYaml(taskmap, []byte(GROCERIES))
	expected_tasks := []string{"buy groceries", "buy milk", "buy bread"}
	for i := range 3 {
		assert.Equal(t, expected_tasks[i], taskmap["t"+strconv.Itoa(i+1)].Task.String())
	}
}

func TestTaskMapParseYamlFail(t *testing.T) {
	taskmap := make(TaskMap)
	assert.Panics(t, func() { ParseYaml(taskmap, []byte(GROCERIES+"laksjdflj;alsdrandombytesidk")) })
}

func TestNextKey(t *testing.T) {
	taskmap := make(TaskMap)
	ParseYaml(taskmap, []byte(GROCERIES))
	for i := range 20 {
		assert.Equal(t, taskmap.NextKey(""), "t"+strconv.Itoa(4+i))
		taskmap["t"+strconv.Itoa(4+i)] = &Task{}
	}
	taskmap["t25"] = &Task{}
	assert.Equal(t, "t24", taskmap.NextKey(""))
	taskmap["t24"] = &Task{}
	assert.Equal(t, "t26", taskmap.NextKey(""))
}

func TestDo(t *testing.T) {
	taskmap := make(TaskMap)
	ParseYaml(taskmap, []byte(GROCERIES))
	assert.True(t, taskmap["t2"].IsDone(taskmap, time.Now()))
	assert.False(t, taskmap["t3"].IsDone(taskmap, time.Now()))
	taskmap.Do("t3", time.Now(), false)
	assert.True(t, taskmap["t3"].IsDone(taskmap, time.Now()))
}

func TestUndo(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(GROCERIES))
	assert.True(t, tm["t2"].IsDone(tm, time.Now()))
	tm.Undo("t2", time.Now())
	assert.False(t, tm["t2"].IsDone(tm, time.Now()))
}

func TestUndoDoneAt(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(GROCERIES))
	tm.Undo("t2", time.Now())
	fmt.Println(tm["t2"].DoneAt)
	assert.True(t, tm["t2"].DoneAt.IsZero())
}

func TestDepIsDone(t *testing.T) {
	taskmap := make(TaskMap)
	ParseYaml(taskmap, []byte(GROCERIES))
	assert.False(t, taskmap["t1"].IsDone(taskmap, time.Now()))
	taskmap.Do("t2", time.Now(), false)
	assert.False(t, taskmap["t1"].IsDone(taskmap, time.Now()))
	taskmap.Do("t3", time.Now(), false)
	assert.True(t, taskmap["t1"].IsDone(taskmap, time.Now()))
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

func TestTfidfNextKey(t *testing.T) {
	config := TfidfConfig{Enabled: true}
	tm := make(TaskMap)
	ParseYaml(tm, []byte(HOMEWORKS))
	msg := "buy some laptop for uni"
	key := tm.TfidfNextKey(msg, config, "")
	tm[key] = &Task{Task: TemplateBase{Template: msg}}
	assert.Equal(t, "laptop", key)
	key = tm.TfidfNextKey("buy some milk (fresh)", config, "milk")
	assert.Equal(t, "milk", key)
	config = TfidfConfig{Enabled: false}
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
	now := time.Now()
	tm.ResolveAllTemplates(now)
	tm2 := LoadTaskMap("test", now)
	assert.Equal(t, tm, tm2)
	assert.Error(t, tm.Write(""))
	os.Remove("test")
}

func TestFindDoneAt(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(GROCERIES))
	tm.Do("t3", time.Now(), false)
	assert.Equal(t, tm["t3"].DoneAt, tm["t1"].FindDoneAt(tm))
}

func TestWipeDependenciesToKey(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(GROCERIES))
	tm.WipeDependenciesToKey("t3")
	assert.Equal(t, []string{"t2"}, tm["t1"].Deps)
}

func BenchmarkPrintMarkdown(b *testing.B) {
	tm := LoadTaskMap("../tests/tasks.yaml", time.Now())
	for b.Loop() {
		tm.PrintMarkdown(&MarkdownConfig{Indent: 4}, nil)
	}
}

func TestAddDep(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(HOMEWORKS))
	task, err := tm.AddDep("homework", "coco")
	assert.Error(t, err)
	task, err = tm.AddDep("homework", "milk")
	assert.Nil(t, err)
	assert.Equal(t, []string{"study", "project", "milk"}, task.Deps)
}

func TestEmptyTaskMapMarkdownError(t *testing.T) {
	tm := make(TaskMap)
	assert.Error(t, tm.PrintMarkdown(&MarkdownConfig{Indent: 4}, nil))
}

func TestCascadeDeps(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(HOMEWORKS))

	tm["homework"].CascadeOrphanDeps(tm)
	assert.Equal(t, 4, len(tm))
	task := tm["homework"]
	delete(tm, "homework")
	assert.Equal(t, 3, len(tm))
	task.CascadeOrphanDeps(tm)
	assert.Equal(t, 1, len(tm))
	for key := range tm {
		assert.Equal(t, "milk", key)
	}
}
func TestSortedKeys(t *testing.T) {
	tm := make(TaskMap)
	tm["ydo"] = &Task{Task: TemplateBase{Template: "make ydo usable"}}
	tm["milk"] = &Task{Task: TemplateBase{Template: "buy milk"}, Due: TemplateDate{Date: time.Now().Add(time.Hour * 2)}}
	tm["workout"] = &Task{Task: TemplateBase{Template: "workout"}, Due: TemplateDate{Date: time.Now().AddDate(1000, 0, 0)}}
	tm["homework"] = &Task{Task: TemplateBase{Template: "do homework"}, Due: TemplateDate{Date: time.Now().Add(time.Hour * 12)}}

	assert.Equal(t, []string{"milk", "homework", "workout", "ydo"}, tm.SortedKeys())
}

func ExampleTaskMap_PrintMarkdown() {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(HOMEWORKS))
	tm.Do("study", time.Now(), false)
	task := tm["homework"]
	task.Due = TemplateDate{Date: time.Now().Add(time.Hour * 24)}
	tm["homework"] = task
	task = tm["milk"]
	task.Due = TemplateDate{Date: time.Now().Add(time.Minute * 12)}
	tm["milk"] = task
	config := MarkdownConfig{Indent: 4, Now: time.Now()}
	config.Init()
	tm.PrintMarkdown(&config, nil)
	// Output:
	// - [ ] milk: buy some milk (12min)
	// - [ ] homework: do uni practice (1d)
	//     - [x] study: study for the uni exam (0s ago)
	//     - [ ] project: do the hobby project
}

const SINGLE_DISCIPLINE = `workout:
  task: workout
  created-at: 2025-03-20T00:20:52.625601175+03:30
  done-at: 2025-03-26T08:46:50.967817015+03:30
  done-at-archive:
    - 2025-03-23T08:46:50.967817015+03:30
    - 2025-03-25T08:46:50.967817015+03:30
  recur: 1d
`

func TestTrackDisciplineDailySingle(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(SINGLE_DISCIPLINE))
	start, _ := time.ParseInLocation("2006-01-02", "2025-03-20", time.Local)
	end, _ := time.ParseInLocation("2006-01-02", "2025-03-27", time.Local)
	assert.Equal(t, []float64{0, 0, 0, 1, 0, 1, 1, 0}, tm.TrackDisciplineDaily(start, end, time.Now()))
}

const MULTIPLE_DISCIPLINE = `clean:
    task: clean room
    done: true
    done-at: 2025-03-25T11:12:53.584720168+03:30
    created-at: 2025-03-21T08:00:00.587348085+03:30
    recur: 1w
library:
    task: call library to renew book
    done: true
    done-at: 2025-03-23T00:49:08.587348085+03:30
    created-at: 2025-03-21T08:00:00.587348085+03:30
    recur: 2w
plants-water:
    task: change plants water
    done: true
    created-at: 2025-03-26T12:27:19.204076503+03:30
    done-at: 2025-03-27T08:38:13.019924716+03:30
    recur: 3d
wax:
    task: wax the shoe
    done: true
    done-at: 2025-03-22T11:12:36.228830075+03:30
    created-at: 2025-03-21T08:00:00.587348085+03:30
    recur: 1m
workout:
    task: workout
    done: true
    created-at: 2025-03-26T00:20:52.625601175+03:30
    done-at: 2025-03-27T08:46:50.967817015+03:30
    recur: 1d
`

func TestTrackDisciplineDailyMultiple(t *testing.T) {
	tm := make(TaskMap)
	ParseYaml(tm, []byte(MULTIPLE_DISCIPLINE))
	start, _ := time.ParseInLocation("2006-01-02", "2025-03-21", time.Local)
	end, _ := time.ParseInLocation("2006-01-02", "2025-03-27", time.Local)
	assert.Equal(t, []float64{0, 1. / 3., 1. / 2., 0, 1, 0, 1}, tm.TrackDisciplineDaily(start, end, time.Now()))
}

const COMPLEX_DISCIPLINE = `t1:
  created-at: 2025-03-20T00:20:52.625601175+03:30
  done-at: 2025-03-21T00:20:52.625601175+03:30
  done: true
t2:
  created-at: 2025-03-21T00:20:52.625601175+03:30
  done-at: 2025-03-22T00:20:52.625601175+03:30
  done: true
t3:
  created-at: 2025-03-24T00:20:52.625601175+03:30
`

func TestTrackDisciplineDailyNoRecur(t *testing.T) {
	expected := []float64{0, 0.5, 1, math.NaN(), 0, 0}

	tm := make(TaskMap)
	ParseYaml(tm, []byte(COMPLEX_DISCIPLINE))
	start, _ := time.ParseInLocation("2006-01-02", "2025-03-20", time.Local)
	end, _ := time.ParseInLocation("2006-01-02", "2025-03-25", time.Local)
	assert.Equal(t, fmt.Sprint(expected), fmt.Sprint(tm.TrackDisciplineDaily(start, end, time.Now())))
}
