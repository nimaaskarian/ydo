package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCmdRoot(t *testing.T) {
  assert.Nil(t, taskmap)
  // system specific. fails if you don't have any tasks on your system
  assert.Nil(t, rootCmd.Execute())
  assert.NotNil(t, taskmap)
  taskmap = nil
  rootCmd.SetArgs([]string{"-f", os.DevNull})
  assert.ErrorContains(t, rootCmd.Execute(), "No tasks")
  assert.NotNil(t, taskmap)
  taskmap = nil
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml"})
  assert.Nil(t, rootCmd.Execute())
  taskmap = nil
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "acommandthatwonteverexit"})
  assert.ErrorContains(t, rootCmd.Execute(), "unknown command")
  taskmap = nil
}

func TestCmdRootRm(t *testing.T) {
  assert.Nil(t, taskmap)
  assert.Error(t, rootCmd.Execute(), "No tasks")
  file := []string{"-f", "../tests/tasks.yaml"}
  rootCmd.SetArgs(append(file, "edit","-n", "tests", "we have edited the message!"))
  assert.Nil(t, rootCmd.Execute())
  assert.Equal(t,"../tests/tasks.yaml", tasks_path)
  assert.NotNil(t, taskmap)
  assert.Nil(t, rootCmd.Execute())
  rootCmd.SetArgs(append(file, "rm","willneverexist"))
  assert.ErrorContains(t, rmCmd.Execute(), "No such task")
  taskmap = nil
  rootCmd.SetArgs(append(file, "rm","tests"))
  assert.Nil(t, rmCmd.Execute(), "No such task")
  assert.Empty(t, taskmap["tests"])
  taskmap = nil
}

func TestCmdRegenKey(t *testing.T) {
  assert.Nil(t, taskmap)
  base := []string{"-f", "../tests/tasks.yaml", "-n", "-Y", "regen-key"}
  rootCmd.SetArgs(append(base, "taskthatwillneverexist"))
  assert.ErrorContains(t, rootCmd.Execute(), "No such task")
  rootCmd.SetArgs(base)
  assert.Nil(t, rootCmd.Execute())
  taskmap = nil
  rootCmd.SetArgs(append(base, "tests"))
  assert.Nil(t, rootCmd.Execute())
  taskmap = nil
}

func TestCmdAdd(t *testing.T) {
  assert.Nil(t, taskmap)
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add"})
  assert.ErrorContains(t, rootCmd.Execute(), "at least 1 arg")
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add", "", "", ""})
  assert.ErrorContains(t, rootCmd.Execute(), "Task cannot be empty")
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add", "-k", "lowkey", "some", "task", "msg", "which is", "kinda odd"})
  assert.Nil(t, rootCmd.Execute())
  assert.Equal(t, "some task msg which is kinda odd", taskmap["lowkey"].Task)
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add", "some", "task", "msg", "which is", "kinda odd"})
  assert.Nil(t, rootCmd.Execute())
}

func TestMutatingCmdGotRightFuncs(t *testing.T) {
  cmds := []*cobra.Command{
    addCmd,
    editCmd,
    rmCmd,
    doCmd,
    undoCmd,
    regenKeyCmd,
  }
  save_changes := reflect.ValueOf(SaveChanges)
  save_old_map := reflect.ValueOf(UpdateOldTaskMap)
  for _, cmd := range cmds {
    if !assert.Equal(t, save_changes, reflect.ValueOf(cmd.PostRunE)) {
      defer fmt.Printf("ERROR The command %q doesn't have the SaveChanges PostRunE\n", cmd.Name())
      break
    }
    if !assert.Equal(t, save_old_map, reflect.ValueOf(cmd.PreRun)) {
      defer fmt.Printf("ERROR The command %q doesn't have the UpdateOldMap PreRun\n", cmd.Name())
      break
    }
  }
}

func TestDontAssignToOldMap(t *testing.T) {
  err := filepath.Walk("../cmd", func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    if filepath.Ext(path) != ".go" || filepath.Base(path) == "root.go" || filepath.Base(path) == "cmd_test.go" {
      return nil
    }
    content, err := os.ReadFile(path)
    if err != nil {
      return err
    }
    if !assert.False(t, bytes.Contains(content, []byte("old_taskmap ="))) {
      return fmt.Errorf("%s has old_taskmap assignment\n", path)
    }
    return nil
  })
  assert.Nil(t, err)
}

func TestCompletions(t *testing.T) {
  DueCompletion(nil, nil, "")
}
