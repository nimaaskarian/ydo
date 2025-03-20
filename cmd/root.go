package cmd

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func TaskKeyCompletionFilter(filter func(core.Task, core.TaskMap) bool) func (*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
  return func (cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    taskmap = core.LoadTaskMap(tasks_path)
    keys := make([]string, len(taskmap))
    i := 0
    for key := range taskmap {
      if (filter == nil || filter(taskmap[key], taskmap)) && strings.HasPrefix(key, toComplete) {
        keys[i] = key
        i++
      }
    }
    return keys, cobra.ShellCompDirectiveDefault
  }
}

func TaskKeyCompletionOnFirst(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  if len(args) > 0 {
    return []string{}, cobra.ShellCompDirectiveDefault
  }
  return TaskKeyCompletionFilter(nil)(cmd, args, toComplete)
}

func PrintMarkdown(md_config *core.MarkdownConfig) error {
  switch md_config.Mode {
  case "todo":
    return taskmap.PrintMarkdown(core.Task.IsNotDone, md_config)
  case "md":
    return taskmap.PrintMarkdown(nil, md_config)
  default:
    if len(taskmap) >= 10 {
      return taskmap.PrintMarkdown(core.Task.IsNotDone, md_config)
    } else {
      return taskmap.PrintMarkdown(nil, md_config)
    }
  }
}

type Config struct  {
  // files to look for if --file option is not present
  Files []string `yaml:",omitempty"`
  LogLevel string `yaml:",omitempty"`
  Tfidf core.TfidfConfig `yaml:",omitempty"`
  Markdown core.MarkdownConfig `yaml:",omitempty"`
}

func (config *Config) ReadFile(path string) {
  content, _ := os.ReadFile(path)
  err := yaml.Unmarshal([]byte(content), config)
  if err != nil {
    slog.Error("Error reading config file", "err", err)
  }
  if config.Markdown.Indent == 0 {
    config.Markdown.Indent = 3
  }
}

func (config *Config) FirstFileAvailable() (string, error) {
  for _, file := range config.Files {
    if _, err := os.Stat(file); err == nil {
      return file, nil
    }
  }
  if _, err := os.Stat(config_dir); err == nil {
    path := filepath.Join(config_dir, "tasks.yaml")
    slog.Info("No tasks file available. Using task file in default path", "path", path)
    return path, nil
  }
  return "", errors.New("No file available")
}

func (config *Config) SlogLevel() slog.Level {
  switch config.LogLevel {
    case "debug":
    return slog.LevelDebug
    case "warn":
    return slog.LevelWarn
    case "info":
    return slog.LevelInfo
    default:
    return slog.LevelError
  }
}

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
      return PrintMarkdown(&config.Markdown)
  },
  PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
    if !reflect.DeepEqual(old_taskmap, taskmap) {
      slog.Debug("TaskMap has changed. Writing to file.", "old", old_taskmap, "new", taskmap)
      if err := PrintMarkdown(&config.Markdown); err != nil {
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
