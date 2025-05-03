package cmd

import (
	"slices"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mdCmd)
	mdCmd.Flags().StringVarP(&due, "due", "u", "", "specify due for the tasks to print")
	mdCmd.RegisterFlagCompletionFunc("due", DueCompletion)
	mdCmd.ValidArgsFunction = TaskKeyCompletionFilter(nil)

	mdCmd.Flags().StringArrayVarP(&flagTask.Tags, "tag", "T", []string{}, "tag(s) for the task")
	mdCmd.RegisterFlagCompletionFunc("tag", TagCompletion)
}

var mdCmd = &cobra.Command{
  Aliases: []string {
    "markdown",
    "all",
    "la",
  },
	Use:   "md [tasks (optional)]",
	Short: "output tasks as markdown",
	Long:  "output tasks as markdown. run with no args to list all tasks",
	RunE: func(cmd *cobra.Command, keys []string) error {
		due_time, err := utils.ParseDue(due, now)
		if err != nil {
			return err
		}
		md_config := config.Markdown
		md_config.Limit = 0
		markdownFilter := func(task *core.Task, taskmap core.TaskMap, now time.Time) bool {
			return (due_time.IsZero() || utils.NaiveDateEqual(task.Due.Value(), due_time)) &&
				(len(flagTask.Tags) == 0 || slices.ContainsFunc(flagTask.Tags, func(tag string) bool {
					return slices.Contains(task.Tags, tag)
				}))
		}
		if len(keys) == 0 {
			taskmap.PrintMarkdown(&md_config, markdownFilter)
		} else {
			md_config.Description = true
			seen_keys := make(map[string]bool, len(keys))
			for _, key := range keys {
				task, err := taskmap.GetTask(key)
				if err != nil {
					return err
				}
				task.PrintMarkdown(taskmap, 0, seen_keys, key, &md_config, markdownFilter)
			}
		}
		return nil
	},
}
