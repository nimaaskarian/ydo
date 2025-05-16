package cmd

import (
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var (
	force bool
)

func init() {
	rootCmd.AddCommand(doCmd)
	doCmd.Flags().BoolVarP(&force, "force", "F", false, "Force do task, ignore if its already done or not. Using this you might override the DoneAt data.")
	doCmd.ValidArgsFunction = TaskKeyCompletionFilter(func(t *core.Task, tm core.TaskMap, now time.Time) bool {
		return !t.AutoComplete && !t.IsDone(tm, now) && !t.IsDeleted(now)
	})
}

var doCmd = &cobra.Command{
	Use:   "do [tasks]",
	Short: "set tasks as completed",
	RunE: func(cmd *cobra.Command, keys []string) error {
		if len(keys) > 0 {
			for _, key := range keys {
				if err := taskmap.Do(key, now, force); err != nil {
					return err
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
