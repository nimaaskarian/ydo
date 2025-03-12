package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the todo to add")
}

var editCmd = &cobra.Command{
  Use: "edit",
  Short: "edit a task",
  Run: func(cmd *cobra.Command, args []string) {
    utils.MustHaveTask(taskmap, key)
    task := strings.Join(args, " ")
    if task == "" {
      task = taskmap[key].Task
    }
    for _,key := range deps {
      if _, ok := taskmap[key]; !ok {
        log.Fatalf("No such task %q\n",key)
      }
    }
    taskmap[key] = core.Task {Task: task, Deps: deps}
    utils.WriteTaskmap(taskmap, path)
    fmt.Printf("Task %q created\n", key)
  },
}
