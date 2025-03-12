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

var deps []string
func init() {
  rootCmd.AddCommand(addCmd)
  addCmd.Flags().StringArrayVarP(&deps, "deps", "d", []string{}, "dependencies for the todo to add")
}

var addCmd = &cobra.Command{
  Use: "add",
  Short: "add a todo",
  Run: func(cmd *cobra.Command, args []string) {
    taskmap = utils.LoadTasks(path)
    utils.MustNotHaveTask(taskmap, key)
    task := strings.Join(args, " ")
    if task == "" {
      log.Fatalln("Task can't be empty")
    }
    for _,key := range deps {
      utils.MustHaveTask(taskmap, key)
    }
    taskmap[key] = core.Todo {Task: task, Deps: deps}
    str, err := yaml.Marshal(taskmap)
    utils.Check(err)
    err = os.WriteFile(path, str, 0644)
    utils.Check(err)
    fmt.Printf("Task %q created\n", key)
  },
}
