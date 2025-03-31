package cmd

import (
	"bufio"
	"fmt"
	"os"

	// "os/signal"

	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(addBatchCmd)
}
var addBatchCmd = &cobra.Command{
  Aliases: []string{"ab"},
  Use: "add-batch [path to file (optional)]",
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
        task := string(line)
        key := taskmap.TfidfNextKey(task, config.Tfidf, "")
        taskmap[key] = &core.Task{Task: task}
      }
    }
    return nil
  },
  PostRunE: SaveChanges,
  PreRun: UpdateOldTaskMap,
}
