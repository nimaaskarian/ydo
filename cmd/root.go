package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var (
  // flags
  key string
  path string
  // global state
  taskmap core.TaskMap

  rootCmd = &cobra.Command{
  Use:   "ydo",
  Short: "ydo is a frictionless and fast todo app",
  Long: `Fast, featurefull and frictionless todo app with a graph structure`,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    taskmap = utils.LoadTasks(path)
    if taskmap == nil {
      taskmap = make(core.TaskMap)
    }
  },
  Run: func(cmd *cobra.Command, args []string) {
    if todo, ok := taskmap[key]; ok {
      todo.PrintMarkdown(taskmap, 1)
    } else {
      taskmap.PrintMarkdown()
    }
  },
}
)

func init() {
  dir := utils.ConfigDir()
  rootCmd.PersistentFlags().StringVarP(&path, "file","f",filepath.Join(dir, "todos.yaml"), "path to todo file")
  rootCmd.PersistentFlags().StringVarP(&key, "key","k", "", "todo key")
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
