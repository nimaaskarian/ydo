package cmd

import (
  "github.com/spf13/cobra"
  "github.com/nimaaskarian/ydo/utils"
  "github.com/nimaaskarian/ydo/core"
  "os"
  "fmt"
)

var rootCmd = &cobra.Command{
  Use:   "ydo",
  Short: "ydo is a frictionless and fast todo app",
  Long: `Fast and featureful todo app 
made with love by nimaaskarian`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

var key string
var todomap core.TodoMap
func Execute() {
  dir := utils.ConfigDir()
  todomap = utils.LoadTodos(dir)
  rootCmd.PersistentFlags().StringVarP(&key, "key","k", todomap.NextKey(), "todo key")
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
