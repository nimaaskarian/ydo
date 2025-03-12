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
  rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
  Use: "add",
  Short: "add a todo",
  Run: func(cmd *cobra.Command, args []string) {
    dir := utils.ConfigDir()
    todomap := utils.LoadTodos(dir)
    if _, ok := todomap[key]; ok {
      log.Fatalf("Task %q already exists", key)
    }
    todomap[key] = core.Todo {Task: strings.Join(args, " ")}
    fmt.Printf("Task %q created\n", key)
  },
}
