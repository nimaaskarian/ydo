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
  todomap core.TodoMap

  rootCmd = &cobra.Command{
  Use:   "ydo",
  Short: "ydo is a frictionless and fast todo app",
  Long: `Fast and featureful todo app 
made with love by nimaaskarian`,
  Run: func(cmd *cobra.Command, args []string) {
    todomap = utils.LoadTodos(path)
    core.PrintYaml(todomap)
  },
}
)

func init() {
  dir := utils.ConfigDir()
  rootCmd.PersistentFlags().StringVarP(&path, "file","f",filepath.Join(dir, "todos.yaml"), "path to todo file")
  rootCmd.PersistentFlags().StringVarP(&key, "key","k", todomap.NextKey(), "todo key")
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
