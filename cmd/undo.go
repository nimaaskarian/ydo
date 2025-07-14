package cmd

import (
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(undoCmd)
}

var undoCmd = &cobra.Command{
	Aliases: []string{"u"},
	Use:     "undo [keys]",
	Short:   "set tasks as not completed",
	ValidArgsFunction: TaskKeyCompletionFilter(func(t *core.Task, tm core.TaskMap, now time.Time) bool {
		return !t.AutoComplete && t.IsDone(tm, now) && !t.IsDeleted(now)
	}),
	RunE: func(cmd *cobra.Command, keys []string) error {
		if len(keys) > 0 {
			for _, key := range keys {
				for _, key := range taskmap.RegexpMatchingKeys(key, config.Regexp) {
					event, err := taskmap.Undo(key, now)
					if err != nil {
						return err
					}
					events = append(events, event)
				}
			}
		} else {
			if always_yes || utils.ReadYesNo("This will set all tasks as not completed. ARE YOU REALLY SURE? (yes/no) ") {
				for key := range taskmap {
					event, _ := taskmap.Undo(key, now)
					events = append(events, event)
				}
			}
		}
		return nil
	},
	PostRunE: SaveChanges,
}
