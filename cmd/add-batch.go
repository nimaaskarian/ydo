package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(addBatchCmd)
}
var addBatchCmd = &cobra.Command{
  Aliases: []string{"b"},
  Use: "batch [path to file (optional)]",
  Short: "add a list of tasks from file or stdin",
  Long: "read each line of stdin (or the given file) as a series of key: task (key defaults to auto), and add them",
  Args: cobra.MaximumNArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
    file := os.Stdin
    if len(args) > 0 {
      var err error
      file, err = os.Open(args[0])
      if err != nil {
        return err
      }
    }
    reader := bufio.NewReader(file)
    if file == os.Stdin {
      fi, _ := os.Stdin.Stat()
      if (fi.Mode() & os.ModeCharDevice) != 0 {
        fmt.Println("C-d to quit and save. C-c to quit and abort")
      }
    }
    for {
      line, _, err := reader.ReadLine()
      if err != nil {
        break
      }
      if len(line) != 0 {
        key_task := strings.Split(string(line), ":")
        var key, task string
        if len(key_task) >= 2 {
          task = strings.TrimSpace(strings.Join(key_task[1:], " "))
          key = key_task[0]
        } else {
          key = taskmap.TfidfNextKey(task, config.Tfidf, "")
          task = key_task[0]
        }
        taskmap[strings.TrimSpace(key)] = &core.Task{Task: task}
      }
    }
    return nil
  },
  PostRunE: SaveChanges,
  PreRun: UpdateOldTaskMap,
}
