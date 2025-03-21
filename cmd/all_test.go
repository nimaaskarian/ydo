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

func TestCmdRootMutate(t *testing.T) {
  assert.Nil(t, taskmap)
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "-n","edit", "tests", "we have edited the message!"})
  assert.Nil(t, rootCmd.Execute())
  taskmap = nil
  rootCmd.SetArgs([]string{"-f", "../tests/tasks.yaml", "-Y","del"})
  assert.Error(t, rootCmd.Execute(), "No tasks")
}
