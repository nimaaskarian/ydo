package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var (
  // global flags
  tasks_path string
  config_path string
  dry_run bool
  // global state
  old_taskmap, taskmap core.TaskMap
  
  config_dir string
  config Config

  rootCmd = &cobra.Command{
  SilenceErrors: true,
  SilenceUsage: true,
  Use:   "ydo",
  Short: "ydo is a frictionless and fast to-do app",
  Long: `Fast, featurefull and frictionless to-do app with a graph structure`,
  PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    config = Config{};
    config.ReadFile(config_path)
    loglevel := config.SlogLevel()
    slog.SetLogLoggerLevel(loglevel)
    log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
    slog.Info("Config file loaded", "path", config_path)
    slog.Info("Log level set", "loglevel", loglevel)
    if tasks_path == "" {
      var err error
      tasks_path, err = config.FirstFileAvailable()
      if err != nil {
        return err
      }
    }
    taskmap = core.LoadTaskMap(tasks_path)
    if taskmap == nil {
      taskmap = core.TaskMap{}
    }
    old_taskmap = utils.DeepCopyMap(taskmap)
    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
      return taskmap.PrintMarkdown(&config.Markdown)
  },
  PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
    if !reflect.DeepEqual(old_taskmap, taskmap) {
      slog.Debug("TaskMap has changed. Writing to file.", "old", old_taskmap, "new", taskmap)
      if err := taskmap.PrintMarkdown(&config.Markdown); err != nil {
        return err
      }
      if dry_run {
        taskmap.DryWrite(tasks_path)
      } else {
        taskmap.Write(tasks_path)
      }
    }
    return nil
  },
}
)

func init() {
  config_dir = utils.ConfigDir()
  rootCmd.PersistentFlags().StringVarP(&tasks_path, "file","f","", "path to tasks file")
  rootCmd.PersistentFlags().StringVarP(&config_path, "config","c",filepath.Join(config_dir, "config.yaml"), "path to config file")
  rootCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run","n", false, "perform a trial run with no changes made")
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
