package cmd

import (
	"log/slog"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var (
	force bool
)

func doCmdInclude(t *core.Task, tm core.TaskMap, now time.Time) bool {
	return !t.AutoComplete && !t.IsDone(tm, now) && !t.IsDeleted(now)
}

func init() {
	rootCmd.AddCommand(doCmd)
	doCmd.AddCommand(interactiveDoCmd)
	doCmd.PersistentFlags().BoolVarP(&force, "force", "F", false, "Force do task, ignore if its already done or not. Using this you might override the DoneAt data.")
	doCmd.ValidArgsFunction = TaskKeyCompletionFilter(doCmdInclude)
}

var doCmd = &cobra.Command{
	Use:   "do [tasks]",
	Short: "set tasks as completed",
	RunE: func(cmd *cobra.Command, keys []string) error {
		if len(keys) > 0 {
			for _, key := range keys {
				for _, key := range taskmap.RegexpMatchingKeys(key, config.Regexp) {
					if err := taskmap.Do(key, now, force); err != nil {
						return err
					}
				}
			}
		} else {
			if always_yes || utils.ReadYesNo("This will set all tasks as completed. ARE YOU REALLY SURE? (yes/no) ") {
				for key := range taskmap {
					taskmap.Do(key, now, force)
				}
			}
		}
		return nil
	},
	PostRunE: SaveChanges,
	PreRun:   UpdateOldTaskMap,
}

var interactiveDoCmd = &cobra.Command{
	Use:   "interactive",
	Short: "interactively set tasks as done",
	Long: "interactively set tasks as done using your EDITOR. all removed lines will be done",
	RunE: func(cmd *cobra.Command, keys []string) error {
		task_should_do, err := interactiveHelper("ydo-interactive-do", doCmdInclude)
		if err != nil {
			return err
		}
		for key, should_do := range task_should_do {
			if should_do {
				slog.Info("Interactively doing task", "key", key)
				taskmap.Do(key, now, force)
			}
		}
		return nil
	},
	PostRunE: SaveChanges,
	PreRun:   UpdateOldTaskMap,
}
