package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func TaskKeyCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  taskmap = utils.LoadTasks(path)

  keys := make([]string, len(taskmap))

  i := 0
  for key := range taskmap {
		if strings.HasPrefix(key, toComplete) {
      keys[i] = key
      i++
    }
  }
  return keys, cobra.ShellCompDirectiveDefault
}


var (
  // flags
  key string
  path string
  // global state
  taskmap core.TaskMap

  rootCmd = &cobra.Command{
  Use:   "ydo",
  Short: "ydo is a frictionless and fast to-do app",
  Long: `Fast, featurefull and frictionless to-do app with a graph structure`,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    taskmap = utils.LoadTasks(path)
    if taskmap == nil {
      taskmap = make(core.TaskMap)
    }
  },
  Run: func(cmd *cobra.Command, args []string) {
    if task, ok := taskmap[key]; ok {
      task.PrintMarkdown(taskmap, 1, []string{key}, nil)
    } else {
      taskmap.PrintMarkdown()
    }
  },
}
)

func init() {
  dir := utils.ConfigDir()
  rootCmd.PersistentFlags().StringVarP(&path, "file","f",filepath.Join(dir, "tasks.yaml"), "path to task file")
  rootCmd.PersistentFlags().StringVarP(&key, "key","k", "", "task key")
  rootCmd.RegisterFlagCompletionFunc("key", TaskKeyCompletion)

}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
