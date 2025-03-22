package cmd

import (
	"os"
	"testing"

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
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "-n","edit", "tests", "we have edited the message!"})
  assert.Nil(t, rootCmd.Execute())
  taskmap = nil
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "-Y","rm"})
  assert.Error(t, rootCmd.Execute(), "No tasks")
  taskmap = nil
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "rm", "neverwillexist"})
  assert.ErrorContains(t, rootCmd.Execute(), "No such task")
  taskmap = nil
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "rm", "tests"})
  assert.Nil(t, rootCmd.Execute())
  taskmap = nil
}

func TestCmdRegenKey(t *testing.T) {
  assert.Nil(t, taskmap)
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "regen-key", "taskthatwillneverexist"})
  assert.ErrorContains(t, rootCmd.Execute(), "No such task")
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "regen-key"})
  assert.Nil(t, rootCmd.Execute())
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "regen-key", "tests"})
  assert.Nil(t, rootCmd.Execute())
  taskmap = nil
}

func TestCmdAdd(t *testing.T) {
  assert.Nil(t, taskmap)
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add"})
  assert.ErrorContains(t, rootCmd.Execute(), "at least 1 arg")
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add", "", "", ""})
  assert.ErrorContains(t, rootCmd.Execute(), "Task cannot be empty")
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add", "-k", "newkey", "some", "task", "msg", "which is", "kinda odd"})
  assert.Nil(t, rootCmd.Execute())
  assert.Equal(t, "some task msg which is kinda odd", taskmap["newkey"].Task)
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "add", "some", "task", "msg", "which is", "kinda odd"})
  assert.Nil(t, rootCmd.Execute())
}
