package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var (
  // global flags
  tasks_path string
  config_path string
  dry_run bool
  always_yes bool
  now time.Time
  now_str string
  color string
  // global state
  old_taskmap map[string]core.Task;
  taskmap core.TaskMap
  
  config_dir string
  config Config

  rootCmd = &cobra.Command{
  SilenceErrors: true,
  SilenceUsage: true,
  Use:   "ydo",
  Short: "ydo is a frictionless and fast to-do app",
  Long: `Fast, featurefull and frictionless to-do app with a graph structure`,
  PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    now = time.Now()
    if now_str != "" {
      var err error
      now, err = utils.ParseDue(now_str, now)
      if err != nil {
        return err
      }
    }
    config = Config{}
    config.ReadFile(config_path)
    config.Markdown.Now = now
    if color != "" {
      config.Markdown.Color = color
    }
    config.Markdown.Init()
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
    config.Markdown.Filter = MarkdownFilter(&config.Markdown)
    return nil
  },
  RunE: func(cmd *cobra.Command, args []string) error {
      return taskmap.PrintMarkdown(&config.Markdown)
  },
}
)

func SaveChanges(cmd *cobra.Command, args []string) error {
  if !reflect.DeepEqual(old_taskmap, utils.DeepCopyMap(taskmap)) {
    slog.Debug("TaskMap has changed. Writing to file.", "old", old_taskmap, "new", taskmap)
    if dry_run {
      taskmap.DryWrite(tasks_path)
    } else {
      taskmap.Write(tasks_path)
    }
    if err := taskmap.PrintMarkdown(&config.Markdown); err != nil {
      return err
    }
  }
  return nil
}

func UpdateOldTaskMap(cmd *cobra.Command, args []string) {
  old_taskmap = utils.DeepCopyMap(taskmap)
}

func init() {
  config_dir = utils.ConfigDir()
  rootCmd.PersistentFlags().StringVarP(&tasks_path, "file","f","", "path to tasks file")
  rootCmd.PersistentFlags().StringVarP(&config_path, "config","c",filepath.Join(config_dir, "config.yaml"), "path to config file")
  rootCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run","n", false, "perform a trial run with no changes made")
  rootCmd.PersistentFlags().BoolVarP(&always_yes, "always-yes","Y", false, "answer yes to all the yes/no questions")
  
  rootCmd.PersistentFlags().StringVarP(&now_str, "now","N", "", "current time of operations (defaults to current system time)")
  rootCmd.RegisterFlagCompletionFunc("now", DueCompletion)
  rootCmd.PersistentFlags().StringVar(&color, "color", "", "color to print (defaults to auto, overrides config's markdown.color option)")
  rootCmd.RegisterFlagCompletionFunc("color", cobra.FixedCompletions([]string{"always", "never", "auto"},cobra.ShellCompDirectiveNoFileComp))
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
