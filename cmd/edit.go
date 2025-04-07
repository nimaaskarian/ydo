package cmd

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

// edit flags
var (
	new_key          string
	remove_deps      bool
	key_regen        bool
	remove_dep_to    bool
	no_auto_complete bool
	auto_complete    bool
)

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().BoolVarP(&key_regen, "key-regen", "K", false, "regen key using the automatic next key generator (respects the config file)")
	editCmd.Flags().BoolVarP(&flagTask.AutoComplete, "auto-complete", "a", false, "enable auto complete for the task")
	editCmd.Flags().BoolVarP(&no_auto_complete, "no-auto-complete", "A", false, "disable auto complete for the task")
	editCmd.Flags().BoolVar(&remove_dep_to, "clean-dep-to", false, "remove previous 'dependent to' for the task. using this with --dep-to causes to replace 'dependent to's")
	editCmd.Flags().StringVarP(&flagTask.Description.Template, "description", "e", "", "new description of the task")
	editCmd.Flags().BoolVar(&remove_deps, "remove-deps", false, "remove previous dependencies for the task. using this with --deps causes to replace dependencies")
	editCmd.Flags().StringVarP(&new_key, "key", "k", "", "new key to the task")
	editCmd.RegisterFlagCompletionFunc("key", KeyCompletion)

	editCmd.Flags().StringArrayVarP(&flagTask.Tags, "tag", "T", []string{}, "tag(s) for the task")
	editCmd.RegisterFlagCompletionFunc("tag", TagCompletion)

	editCmd.Flags().StringArrayVarP(&flagTask.Deps, "deps", "d", []string{}, "append dependencies for the task")
	editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletionFilter(nil))

	editCmd.Flags().StringVarP(&flagTask.Due.Base.Template, "due", "u", "", "specify due for the tasks to print")
	editCmd.RegisterFlagCompletionFunc("due", DueCompletion)

	editCmd.Flags().StringArrayVarP(&dep_tos, "dep-to", "D", []string{}, "append task keys for this task to be dependent to")
	editCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletionFilter(nil))

	editCmd.Flags().StringVarP(&flagTask.Recur, "recur", "r", "", "duration of in which the ask recurs")
	editCmd.RegisterFlagCompletionFunc("recur", DurationCompletion)

	editCmd.MarkFlagsMutuallyExclusive("no-auto-complete", "auto-complete")
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
		task, err := taskmap.GetTask(args[0])
		if err != nil {
			return err
		}
		if due != "" {
			if due_date, err := utils.ParseDue(due, now); err == nil {
				task.Due = core.NewTemplateDate(due_date)
			}
		}
		edit_key := args[0]
		if new_task, err := TaskTitleFromArgs(args[1:]); err == nil {
			task.Task = core.NewTemplateBase(new_task)
		}
		if !flagTask.Description.IsZero() {
			task.Description = flagTask.Description
		}
		if recur != "" {
			if _, err := utils.ParseDuration(recur, now); err != nil {
				return err
			}

			task.Recur = recur
		}
		if key_regen {
			new_key = taskmap.TfidfNextKey(task.Task.Value(), config.Tfidf, edit_key)
		}
		if _, err := taskmap.GetTask(edit_key); err != nil {
			return err
		}
		for _, dep := range task.Deps {
			if _, err := taskmap.GetTask(dep); err != nil {
				return err
			}
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
			task, err := taskmap.AddDep(dep_key, edit_key)
			if err != nil {
				return err
			}
			taskmap[dep_key] = task
		}
		edit_key = taskmap.ReplaceKeyInDeps(edit_key, new_key)
		if auto_complete {
			task.AutoComplete = true
		}
		if no_auto_complete {
			task.AutoComplete = false
		}
		task.Deps = append(task.Deps, flagTask.Deps...)
		taskmap[edit_key] = task
		if reflect.DeepEqual(taskmap, old_taskmap) {
			return errors.New("Not edited")
		}
		slog.Info("Task edited", "task", task)
		return nil
	},
	PostRunE: SaveChanges,
	PreRun:   UpdateOldTaskMap,
}
