//go:build gtasks
// +build gtasks

package cmd

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/gtasks"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

var gtasksFlags gtasks.Gtasks
var gtasksMarkdownFilter gtasks.GtasksFilter

func init() {
	config_dir = utils.ConfigDir()
	rootCmd.AddCommand(gtasksCmd)
	gtasksCmd.AddCommand(gtasksSyncCmd)
	gtasksCmd.AddCommand(gtasksAddCmd)
	gtasksCmd.AddCommand(gtasksMdCmd)
	gtasksCmd.AddCommand(gtasksTodoCmd)
	gtasksCmd.AddCommand(gtasksDoCmd)
	gtasksCmd.AddCommand(gtasksEditCmd)
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.CredentialsFile, "credentials", filepath.Join(config_dir, "credentials.json"), "path to credentials file (credentials to your google tasks app)")
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.TokenFile, "token", filepath.Join(config_dir, "token.json"), "path to token file (token to your login info)")
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.CacheFile, "cache", filepath.Join(config_dir, "cache.yaml"), "path to cache file")
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.CacheExpire, "cache-expire", "1d", "cache expire duration")
	gtasksCmd.RegisterFlagCompletionFunc("cache-expire", DurationCompletion)
	gtasksAddCmd.ValidArgsFunction = TasklistIdCompletionOnFirst
	gtasksDoCmd.ValidArgsFunction = TasklistSlashTaskIdCompletionFilter(gtasks.IsPending)

}

func TasklistIdCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	gtasksFlags.CacheExpire = "1000y"
	gtasksFlags.ReadCache(now)
	ids := utils.Keys(gtasksFlags.Cache.Tasklists)
	return ids, cobra.ShellCompDirectiveNoFileComp
}

func TasklistIdCompletionOnFirst(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	return TasklistIdCompletion(cmd, args, toComplete)
}

func TaskIdCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	gtasksFlags.CacheExpire = "1000y"
	gtasksFlags.ReadCache(now)
	ids := make([]string, len(gtasksFlags.Cache.Tasks))
	for _, tasks := range gtasksFlags.Cache.Tasks {
		for _, task := range tasks {
			ids = append(ids, task.Id)
		}
	}
	return ids, cobra.ShellCompDirectiveNoFileComp
}

func TaskIdCompletionOnFirst(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	return TaskIdCompletion(cmd, args, toComplete)
}
func TasklistSlashTaskIdCompletionFilter(filter gtasks.GtasksFilter) cobra.CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		gtasksFlags.CacheExpire = "1000y"
		gtasksFlags.ReadCache(now)
		ids := make([]string, len(gtasksFlags.Cache.Tasks))
		for list_id, tasks := range gtasksFlags.Cache.Tasks {
			for _, task := range tasks {
				if filter != nil && !filter(task) {
					continue
				}
				ids = append(ids, list_id+"/"+task.Id)
			}
		}
		return ids, cobra.ShellCompDirectiveNoFileComp
	}
}

func makeGtasksMarkdownFilter(md_config *core.MarkdownConfig) gtasks.GtasksFilter {
	switch md_config.Mode {
	case "todo":
		return gtasks.IsPending
	case "md":
		return nil
	default:
		return func(task *tasks.Task) bool {
			if len(gtasksFlags.Cache.Tasks) >= md_config.Limit {
				return gtasks.IsPending(task)
			} else {
				return true
			}
		}
	}
}

var gtasksCmd = &cobra.Command{
	Use:   "gtasks",
	Short: "handle google tasks",
	Long:  "handle google tasks. running any subcommand of this command will start the google tasks login process",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := rootCmd.PersistentPreRunE(cmd, args); err != nil {
			return err
		}
		gtasksMarkdownFilter = makeGtasksMarkdownFilter(&config.Markdown)
		if err := loginGtasks(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		gtasksFlags.Load(now)
		return gtasksFlags.PrintMarkdown(&config.Markdown, gtasksMarkdownFilter)
	},
}

var gtasksSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync cache with google tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		gtasksFlags.Load(now)
		if err := gtasksFlags.Sync(now); err != nil {
			return err
		}
		return gtasksFlags.PrintMarkdown(&config.Markdown, gtasksMarkdownFilter)
	},
}

var gtasksAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add task",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		gtasksFlags.Load(now)
		list_id := args[0]
		task_title := strings.Join(args[1:], " ")
		task := &tasks.Task{Title: task_title}
		if err := gtasksFlags.AddTaskCache(task, list_id); err != nil {
			return err
		}
		return gtasksFlags.PrintMarkdown(&config.Markdown, gtasksMarkdownFilter)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return gtasksFlags.SaveCache(now)
	},
}

var gtasksTodoCmd = &cobra.Command{
	Use:   "todo",
	Short: "list all pending tasks as markdown",
	RunE: func(cmd *cobra.Command, args []string) error {
		gtasksFlags.Load(now)
		if err := gtasksFlags.PrintMarkdown(&config.Markdown, gtasks.IsPending); err != nil {
			return err
		}
		return nil
	},
}

var gtasksEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit cache file in your favorite editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := utils.EditorCmd(gtasksFlags.CacheFile)
		if err != nil {
			return err
		}
		utils.CmdStdOs(c)
		c.Run()
		return nil
	},
}

var gtasksDoCmd = &cobra.Command{
	Use:   "do [taskListId/taskId]",
	Short: "set tasks in tasklist ids as completed",
	RunE: func(cmd *cobra.Command, args []string) error {
		gtasksFlags.Load(now)
		for _, id := range args {
			ids := strings.Split(id, "/")
			if err := gtasksFlags.DoTaskCache(ids[0], ids[1]); err != nil {
				return err
			}
		}
		if len(args) > 0 {
			return gtasksFlags.PrintMarkdown(&config.Markdown, gtasksMarkdownFilter)
		}
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return gtasksFlags.SaveCache(now)
	},
}

var gtasksMdCmd = &cobra.Command{
	Use:   "md",
	Short: "see google tasks as markdown",
	RunE: func(cmd *cobra.Command, args []string) error {
		gtasksFlags.Load(now)
		if err := gtasksFlags.PrintMarkdown(&config.Markdown, nil); err != nil {
			return err
		}
		return nil
	},
}

func loginGtasks() error {
	slog.Info("Using token file", "token_file", gtasksFlags.TokenFile)
	slog.Info("Using credentials file", "credentials_file", gtasksFlags.CredentialsFile)
	if err := gtasksFlags.Login(); err != nil {
		return err
	}
	return nil
}
