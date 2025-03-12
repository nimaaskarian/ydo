package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
  rootCmd.AddCommand(editCmd)
  editCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the todo to add")
}

var editCmd = &cobra.Command{
  Use: "edit",
  Short: "edit a todo",
  Run: func(cmd *cobra.Command, args []string) {
    taskmap = utils.LoadTasks(path)
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
    taskmap[key] = core.Todo {Task: task, Deps: deps}
    str, err := yaml.Marshal(taskmap)
    utils.Check(err)
    err = os.WriteFile(path, str, 0644)
    utils.Check(err)
    fmt.Printf("Task %q created\n", key)
  },
}
