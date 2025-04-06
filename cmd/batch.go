package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(batchCmd)
}

var batchCmd = &cobra.Command{
  Aliases: []string{"b"},
  Use: "batch [path to file (optional)]",
  Short: "add a list of tasks from file or stdin",
  Long: "read each line of stdin (or the given file) as a series of key: task (key defaults to auto), and add them. it also adds a microsecond to time of operation for each task created so they'd be in order.",
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
    count := 1
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
        err := taskmap.Add(strings.TrimSpace(key), &core.Task{Task: core.NewTemplateBase(task), CreatedAt: now.Add(time.Duration(count))})
        if err != nil {
          fmt.Println(err)
        } else {
          count+=1
        }
      }
    }
    return nil
  },
  PostRunE: SaveChanges,
  PreRun: UpdateOldTaskMap,
}
