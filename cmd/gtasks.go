//go:build gtasks
// +build gtasks

package cmd

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/nimaaskarian/ydo/gtasks"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/tasks/v1"
)

var gtasksFlags gtasks.Gtasks

func init() {
	config_dir = utils.ConfigDir()
	rootCmd.AddCommand(gtasksCmd)
	gtasksCmd.AddCommand(gtasksSyncCmd)
	gtasksCmd.AddCommand(gtasksAddCmd)
	gtasksCmd.AddCommand(gtasksListCmd)
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.CredentialsFile, "credentials", filepath.Join(config_dir, "credentials.json"), "path to credentials file (credentials to your google tasks app)")
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.TokenFile, "token", filepath.Join(config_dir, "token.json"), "path to token file (token to your login info)")
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.CacheFile, "cache", filepath.Join(config_dir, "cache.yaml"), "path to cache file")
	gtasksCmd.PersistentFlags().StringVar(&gtasksFlags.CacheExpire, "cache-expire", "1d", "cache expire duration")
	gtasksCmd.RegisterFlagCompletionFunc("cache-expire", DurationCompletion)
  gtasksAddCmd.ValidArgsFunction = TasklistIdCompletionOnFirst

}

func TasklistIdCompletionOnFirst(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  if len(args) > 0 {
		return []string{}, cobra.ShellCompDirectiveDefault
  }
  gtasksFlags.CacheExpire = "1000y"
  gtasksFlags.ReadCache(now)
  ids := utils.Keys(gtasksFlags.Cache.Tasklists)
  return ids, cobra.ShellCompDirectiveNoFileComp
}

var gtasksCmd = &cobra.Command{
	Use:   "gtasks",
	Short: "handle google tasks",
	Long:  "handle google tasks. running any subcommand of this command will start the google tasks login process",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := rootCmd.PersistentPreRunE(cmd, args); err != nil {
			return err
		}
		UpdateOldTaskMap(cmd, args)
		if err := loginGtasks(); err != nil {
			return err
		}
		return nil
	},
	PostRunE: SaveChanges,
}

var gtasksSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync cache with google tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
    gtasksFlags.Load(now)
		gtasksFlags.Sync(now)
		return gtasksFlags.PrintMarkdown(&config.Markdown)
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
		gtasksFlags.AddTaskCache(task, list_id)
		return gtasksFlags.PrintMarkdown(&config.Markdown)
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return gtasksFlags.SaveCache(now)
	},
}

var gtasksListCmd = &cobra.Command{
	Use:   "md",
	Short: "see google tasks as markdown",
	RunE: func(cmd *cobra.Command, args []string) error {
		gtasksFlags.Load(now)
		if err := gtasksFlags.PrintMarkdown(&config.Markdown); err != nil {
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
