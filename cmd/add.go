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
  rootCmd.AddCommand(addCmd)
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
    todomap[key] = core.Todo {Task: task}
    str, err := yaml.Marshal(todomap)
    utils.Check(err)
    err = os.WriteFile(path, str, 0644)
    utils.Check(err)
    fmt.Printf("Task %q created\n", key)
  },
}
