package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(mdCmd)
  mdCmd.ValidArgsFunction = TaskKeyCompletion
}

var mdCmd = &cobra.Command{
  Use: "md [tasks (optional)]",
  Short: "output tasks as markdown (run with no args so it'd output all tasks like `ydo` does)",
  Run: func(cmd *cobra.Command, args []string) {
    NeedKeysCmd.Run(cmd, args)
    if len(keys) == 0 {
      taskmap.PrintMarkdown()
    } else {
        for _, key := range keys {
          task, ok := taskmap[key];
          if !ok {
            slog.Error("No such task", "key", key)
            os.Exit(1)
          }
          task.PrintMarkdown(taskmap, 1, &[]string{})
        }
    }
  },
}
