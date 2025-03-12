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
    todomap = utils.LoadTodos(path)
    if _, ok := todomap[key]; ok {
      log.Fatalf("Task %q already exists", key)
    }
    if todomap == nil {
      todomap = make(core.TodoMap)
    }
    task := strings.Join(args, " ")
    if task == "" {
      log.Fatalln("Task can't be empty")
    }
    for _,key := range deps {
      if _, ok := todomap[key]; !ok {
        log.Fatalf("No such todo %q\n",key)
      }
    }
    todomap[key] = core.Todo {Task: task, Deps: deps}
    str, err := yaml.Marshal(todomap)
    utils.Check(err)
    err = os.WriteFile(path, str, 0644)
    utils.Check(err)
    fmt.Printf("Task %q created\n", key)
  },
}
