package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var deps []string
func init() {
  rootCmd.AddCommand(addCmd)
  addCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the task to add")
}

var addCmd = &cobra.Command{
  Use: "add",
  Short: "add a task",
  Run: func(cmd *cobra.Command, args []string) {
    if key == "" {
      key = taskmap.NextKey()
    }
    utils.MustNotHaveTask(taskmap, key)
    task := strings.Join(args, " ")
    if task == "" {
      log.Fatalln("Task can't be empty")
    }
    for _,key := range deps {
      utils.MustHaveTask(taskmap, key)
    }
    taskmap[key] = core.Task {Task: task, Deps: deps}
    utils.WriteTaskmap(taskmap, path)
    fmt.Printf("Task %q created\n", key)
  },
}
