package cmd

import (
	"slices"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var show_description bool

func init() {
	rootCmd.AddCommand(todoCmd)
	todoCmd.Flags().StringVarP(&due, "due", "u", "", "specify due for the tasks to print")
	todoCmd.RegisterFlagCompletionFunc("due", DueCompletion)

	todoCmd.Flags().BoolVarP(&show_description, "description", "d", false, "show descriptions")
	todoCmd.Flags().StringArrayVarP(&flagTask.Tags, "tag", "T", []string{}, "tag(s) for the task")
	todoCmd.RegisterFlagCompletionFunc("tag", TagCompletion)
}

var todoCmd = &cobra.Command{
	Aliases:           []string{"t"},
	Use:               "todo [tasks (optional)]",
	Short:             "output to-do as markdown",
	Long:              "output all unfinished tasks (to-dos) as markdown",
	ValidArgsFunction: TaskKeyCompletionFilter((*core.Task).IsNotDone),
	RunE: func(cmd *cobra.Command, keys []string) error {
		due_time, err := utils.ParseDue(due, now)
		if err != nil {
			return err
		}
		md_config := config.Markdown
		md_config.Limit = 0
		md_config.Filter = func(task *core.Task, taskmap core.TaskMap, now time.Time) bool {
			return (due_time.IsZero() || utils.NaiveDateEqual(due_time, task.Due.ToValue())) && task.IsNotDone(taskmap, now) && (len(flagTask.Tags) == 0 || slices.ContainsFunc(flagTask.Tags, func(tag string) bool {
				return slices.Contains(task.Tags, tag)
			}))
		}
		if len(keys) == 0 {
			taskmap.PrintMarkdown(&md_config)
		} else {
			md_config.Description = true
			seen_keys := make(map[string]bool, len(keys))
			for _, key := range keys {
				task, err := taskmap.GetTask(key)
				if err != nil {
					return err
				}
				task.PrintMarkdown(taskmap, 0, seen_keys, key, &md_config)
			}
		}
		return nil
	},
}
