package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

// add flags
var (
	dep_tos []string
	key     string
	due     string
	tfidf   bool
	taskmsg string
	recur   string
)

var flagTask core.Task

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringArrayVarP(&flagTask.Deps, "deps", "d", []string{}, "dependencies for the task to add")
	addCmd.Flags().StringVarP(&flagTask.Description.Template, "description", "e", "", "description of the task")
	addCmd.Flags().StringArrayVarP(&dep_tos, "dep-to", "D", []string{}, "task keys for this task to be dependent to")
	addCmd.Flags().BoolVarP(&flagTask.AutoComplete, "auto-complete", "a", false, "enable auto complete for the task (done when deps are done)")
	addCmd.Flags().BoolVarP(&tfidf, "tfidf", "t", false, "use tfidf for automatic key generation (overrides config file and --key flag)")
	addCmd.Flags().StringVarP(&key, "key", "k", "", "key of the new task")
	addCmd.RegisterFlagCompletionFunc("key", KeyCompletion)

	addCmd.Flags().StringArrayVarP(&flagTask.Tags, "tag", "T", []string{}, "tag(s) for the task")
	addCmd.RegisterFlagCompletionFunc("tag", TagCompletion)

	addCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletionFilter(nil))
	addCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletionFilter(nil))

	addCmd.Flags().StringVarP(&flagTask.Due.Base.Template, "due", "u", "", "specify due for the tasks to print")
	addCmd.RegisterFlagCompletionFunc("due", DueCompletion)

	addCmd.Flags().StringVarP(&flagTask.Until.Base.Template, "until", "U", "", "specify until (task is ignored after that date) for the tasks to print")
	addCmd.RegisterFlagCompletionFunc("until", DueCompletion)

	addCmd.Flags().StringVarP(&flagTask.Schedule.Base.Template, "schedule", "s", "", "specify schedule for the tasks to print")
	addCmd.RegisterFlagCompletionFunc("schedule", DueCompletion)

	addCmd.Flags().StringVarP(&flagTask.Recur, "recur", "r", "", "duration of in which the ask recurs")
	addCmd.RegisterFlagCompletionFunc("recur", DurationCompletion)
}

var addCmd = &cobra.Command{
	Aliases: []string{"a"},
	Use:     "add [your task here yay]",
	Short:   "add a task",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		taskmsg, err = TaskTitleFromArgs(args)
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if tfidf {
			key = taskmap.TfidfNextKey(taskmsg, core.TfidfConfig{Enabled: true}, "")
		} else {
			if key == "" {
				key = taskmap.TfidfNextKey(taskmsg, config.Tfidf, "")
			}
		}
		if _, err := utils.ParseDuration(recur, now); err != nil {
			return err
		}
		flagTask.Task = core.NewTemplateBase(taskmsg)
		flagTask.CreatedAt = now
    if err := checkFlagTaskDateFields(); err != nil {
      return err
    }
    err := taskmap.Add(key, &flagTask)
		if err != nil {
			return err
		}
		for _, dep_to := range dep_tos {
			task, err := taskmap.AddDep(dep_to, key)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			taskmap[dep_to] = task
		}
		fmt.Printf("Task %q added\n", key)
		slog.Debug("Added a task", "task", taskmap[key])
		return nil
	},
	PostRunE: SaveChanges,
	PreRun:   UpdateOldTaskMap,
}

func checkFlagTaskDateFields() error {
  date_fields := flagTask.DateFields()
  for _, item :=  range date_fields {
    date, err := utils.ParseDue(item.Base.Template, now)
    if err == nil {
      *item = core.NewTemplateDate(date)
    }
  }
  for _, item :=  range date_fields {
    if err := resolveTemplateDate(&flagTask, item); err != nil {
      return err
    }
  }
  return nil
}

func resolveTemplateDate(task *core.Task, item *core.TemplateDate) error {
	err := item.Resolve(task)
	if err != nil {
		date, err := utils.ParseDue(item.Base.Template, now)
		if err != nil {
			return err
		}
		*item = core.NewTemplateDate(date)
	}
	return nil
}

func TaskTitleFromArgs(args []string) (taskmsg string, err error) {
	has_non_empty := slices.ContainsFunc(args, func(arg string) bool {
		return arg != ""
	})
	if !has_non_empty {
		return "", errors.New("Task cannot be empty")
	}
	// in task title, replace all newlines with spaces
	return strings.Replace(strings.Join(args, " "), "\n", " ", 0), nil
}
