package cmd

import (
	"log/slog"
	"os"

	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(yamlCmd)
  yamlCmd.ValidArgsFunction = TaskKeyCompletion
}

var yamlCmd = &cobra.Command{
  Use: "yaml",
  Short: "output tasks as yaml",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    if len(keys) == 0 {
      core.PrintYaml(taskmap)
    } else {
        tmp_map := make(core.TaskMap, len(keys))
        for _, key := range keys {
          task, ok := taskmap[key];
          if !ok {
            slog.Error("No such task", "key", key)
            os.Exit(1)
          }
          tmp_map[key] = task
        }
        core.PrintYaml(tmp_map)
    }
  },
}
