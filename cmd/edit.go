package cmd

import (
	"fmt"
	"log/slog"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/hooks"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

// edit flags
var (
	new_key       string
	remove_deps   bool
	key_regen     bool
	remove_dep_to bool
	auto_complete bool
)

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().BoolVarP(&key_regen, "key-regen", "K", false, "regen key using the automatic next key generator (respects the config file)")
	editCmd.Flags().BoolVarP(&flagTask.AutoComplete, "auto-complete", "a", false, "toggle auto complete for the task")
	editCmd.Flags().BoolVar(&remove_dep_to, "clean-dep-to", false, "remove previous 'dependent to' for the task. using this with --dep-to causes to replace 'dependent to's")
	editCmd.Flags().StringVarP(&flagTask.Description.Template, "description", "e", "", "new description of the task")
	editCmd.Flags().BoolVar(&remove_deps, "remove-deps", false, "remove previous dependencies for the task. using this with --deps causes to replace dependencies")
	editCmd.Flags().StringVarP(&new_key, "key", "k", "", "new key to the task")
	editCmd.RegisterFlagCompletionFunc("key", KeyCompletion)

	editCmd.Flags().StringArrayVarP(&flagTask.Tags, "tag", "T", []string{}, "tag(s) for the task")
	editCmd.RegisterFlagCompletionFunc("tag", TagCompletion)

	editCmd.Flags().StringArrayVarP(&flagTask.Deps, "deps", "d", []string{}, "append dependencies for the task")
	editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletionFilter(nil))

	editCmd.Flags().StringVarP(&flagTask.Until.Base.Template, "until", "U", "", "specify until (task is ignored after that date) for the tasks to print")
	editCmd.RegisterFlagCompletionFunc("until", DueCompletion)

	editCmd.Flags().StringVarP(&flagTask.Due.Base.Template, "due", "u", "", "specify due for the tasks to print")
	editCmd.RegisterFlagCompletionFunc("due", DueCompletion)

	editCmd.Flags().StringArrayVarP(&dep_tos, "dep-to", "D", []string{}, "append task keys for this task to be dependent to")
	editCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletionFilter(nil))

	editCmd.Flags().StringVarP(&flagTask.Schedule.Base.Template, "schedule", "s", "", "specify schedule for the tasks to print")
	editCmd.RegisterFlagCompletionFunc("schedule", DueCompletion)

	editCmd.Flags().StringVarP(&flagTask.Recur, "recur", "r", "", "duration of in which the ask recurs")
	editCmd.RegisterFlagCompletionFunc("recur", DurationCompletion)

	editCmd.ValidArgsFunction = TaskKeyCompletionOnFirst

}

var editCmd = &cobra.Command{
	Aliases: []string{"e"},
	Use:     "edit [key (optional)] [new task message (optional)]",
	Short:   "edit a task, or open the data file in favorite editor",
	Long:    "edit a task provided a key. give no keys to open the yaml file in your favorite editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			c, err := utils.EditorCmd(tasks_path)
			if err != nil {
				return err
			}
			utils.CmdStdOs(c)
			return c.Run()
		}
		for _, edit_key := range taskmap.RegexpMatchingKeys(args[0], config.Regexp) {
			task, err := taskmap.GetTask(edit_key)
			if err != nil {
				return err
			}
			if due != "" {
				if due_date, err := utils.ParseDue(due, now); err == nil {
					task.Due = core.NewTemplateDate(due_date)
				}
			}
			if new_task, err := TaskTitleFromArgs(args[1:]); err == nil {
				task.Task = core.NewTemplateBase(new_task)
				events = append(events, hooks.Event{
					Type: hooks.Edit,
					Key: edit_key,
					Secondary: new_task,
				})
			}
			if !flagTask.Description.IsZero() {
				task.Description = flagTask.Description
			}
			if key_regen {
				new_key = taskmap.TfidfNextKey(task.Task.Value(), config.Tfidf, edit_key)
				if new_key != edit_key {
					events = append(events, hooks.Event{Type: hooks.UpdateKey, Key: new_key, Secondary: edit_key})
				}
			}
			for _, dep := range task.Deps {
				if _, err := taskmap.GetTask(dep); err != nil {
					return err
				}
			}
			task.Recur = flagTask.Recur

			taskDateFields := task.DateFields()
			for i, item := range flagTask.DateFields() {
				if !item.IsZero() {
					*taskDateFields[i] = *item
				}
			}
			if err := checkTask(task); err != nil {
				return err
			}

			if remove_deps {
				task.Deps = make([]string, 0)
			}
			if len(flagTask.Tags) > 0 {
				task.Tags = flagTask.Tags
			}
			if remove_dep_to {
				taskmap.WipeDependenciesToKey(edit_key)
			}
			for _, dep_key := range dep_tos {
				task, event, err := taskmap.AddDep(dep_key, edit_key)
				if err != nil {
					return err
				}
				taskmap[dep_key] = task
				events = append(events, event)
			}
			event, edit_key := taskmap.ReplaceKeyInDeps(edit_key, new_key)
			events = append(events, event)
			if auto_complete {
				task.AutoComplete = !task.AutoComplete
				events = append(events, hooks.Event{
					Type:    hooks.Edit,
					Key: fmt.Sprintf("toggle auto complete of %q", edit_key),
				})
			}
			task.Deps = append(task.Deps, flagTask.Deps...)
			for _, dep := range flagTask.Deps {
				events = append(events, hooks.Event{
					Type:             hooks.AddDep,
					Key:          edit_key,
					Secondary: dep,
				})
			}
			taskmap[edit_key] = task
			slog.Info("Task edited", "task", task)
		}
		return nil
	},
	PostRunE: SaveChanges,
}
