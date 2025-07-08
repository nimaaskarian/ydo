package cmd

import (
	"log/slog"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var cascade bool

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.AddCommand(interactiveRmCmd)
	rmCmd.PersistentFlags().BoolVarP(&cascade, "cascade", "C", false, "cascade remove the dependencies of this task that aren't a dependency to other tasks")
	rmCmd.ValidArgsFunction = TaskKeyCompletionFilter(nil)
}

var rmCmd = &cobra.Command{
	Aliases: []string{"del", "delete"},
	Use:     "rm [tasks]",
	Short:   "remove tasks completely from the data (also removes it from all the dependency lists)",
	RunE: func(cmd *cobra.Command, keys []string) error {
		if len(keys) > 0 {
			for _, key := range keys {
				for _, key := range taskmap.RegexpMatchingKeys(key, config.Regexp) {
					if err := taskmap.Delete(key, cascade); err != nil {
						return err
					}
				}
			}
		} else {
			if always_yes || utils.ReadYesNo("WARN This will DELETE ALL THE TASKS. ARE YOU REALLY SURE? (yes/no) ") {
				taskmap = make(core.TaskMap)
			}
		}
		return nil
	},
	PostRunE: SaveChanges,
	PreRun:   UpdateOldTaskMap,
}

var interactiveRmCmd = &cobra.Command{
	Use:   "interactive",
	Short: "interactively remove tasks",
	Long:  "interactively remove tasks using your EDITOR. all removed lines will be done",
	RunE: func(cmd *cobra.Command, keys []string) error {
		task_should_do, err := interactiveHelper("ydo-interactive-rm", nil)
		if err != nil {
			return err
		}
		for key, should_rm := range task_should_do {
			if should_rm {
				slog.Info("Interactively removing task", "key", key)
				taskmap.Delete(key, cascade)
			}
		}
		return nil
	},
	PostRunE: SaveChanges,
	PreRun:   UpdateOldTaskMap,
}
