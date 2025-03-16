package cmd

import (
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

func TaskKeyCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  taskmap = core.LoadTaskMap(tasks_path)
  keys := make([]string, len(taskmap))
  i := 0
  for key := range taskmap {
		if strings.HasPrefix(key, toComplete) {
      keys[i] = key
      i++
    }
  }
  return keys, cobra.ShellCompDirectiveDefault
}

func TaskKeyCompletionOnFirst(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  if len(args) > 0 {
    return []string{}, cobra.ShellCompDirectiveDefault
  }
  return TaskKeyCompletion(cmd, args, toComplete)
}

func TaskKeyCompletionDone(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  taskmap = core.LoadTaskMap(tasks_path)
  keys := make([]string, len(taskmap))
  i := 0
  for key,task := range taskmap {
		if task.IsDone(taskmap) && strings.HasPrefix(key, toComplete) {
      keys[i] = key
      i++
    }
  }
  return keys, cobra.ShellCompDirectiveDefault
}

func TaskKeyCompletionNotDone(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
  taskmap = core.LoadTaskMap(tasks_path)
  keys := make([]string, len(taskmap))
  i := 0
  for key,task := range taskmap {
		if !task.IsDone(taskmap) && strings.HasPrefix(key, toComplete) {
      keys[i] = key
      i++
    }
  }
  return keys, cobra.ShellCompDirectiveDefault
}

var NeedKeysCmd = &cobra.Command{
	Use:   "base [keys]",
	Short: "Base command that takes non-flag arguments",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keys = args
	},
}

type Config struct  {
  // files to look for if --file option is not present
  Files []string `yaml:",omitempty"`
  LogLevel string `yaml:",omitempty"`
  Tfidf core.TfidfConfig `yaml:",omitempty"`
}

func (config *Config) ReadFile(path string) {
  content, _ := os.ReadFile(path)
  err := yaml.Unmarshal([]byte(content), config)
  if err != nil {
    slog.Error("Error reading config file", "err", err)
  }
}

func (config *Config) FirstFileAvailable() string {
  for _, file := range config.Files {
    if _, err := os.Stat(file); err == nil {
      return file
    }
  }
  if _, err := os.Stat(config_dir); err == nil {
    path := filepath.Join(config_dir, "tasks.yaml")
    slog.Info("No tasks file available. Using task file in default path", "path", path)
    return path
  }
  slog.Error("No tasks file available and the default path can't be used")
  slog.Error("Because config directory doesn't exist and can't be created", "dir", config_dir)
  os.Exit(1)
  return ""
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
  keys []string
  tasks_path string
  config_path string
  dry_run bool
  // global state
  old_taskmap, taskmap core.TaskMap
  
  config_dir string
  config Config

  rootCmd = &cobra.Command{
  Use:   "ydo",
  Short: "ydo is a frictionless and fast to-do app",
  Long: `Fast, featurefull and frictionless to-do app with a graph structure`,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    config = Config{};
    config.ReadFile(config_path)
    loglevel := config.SlogLevel()
    slog.SetLogLoggerLevel(loglevel)
    log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
    slog.Info("Config file loaded", "path", config_path)
    slog.Info("Log level set", "loglevel", loglevel)
    if tasks_path == "" {
      tasks_path = config.FirstFileAvailable()
    }
    taskmap = core.LoadTaskMap(tasks_path)
    if taskmap == nil {
      taskmap = core.TaskMap{}
    }
    old_taskmap = utils.DeepCopyMap(taskmap)
  },
  Run: func(cmd *cobra.Command, args []string) {
    if len(taskmap) >= 10 {
      taskmap.PrintMarkdown(core.Task.IsNotDone)
    } else {
      taskmap.PrintMarkdown(nil)
    }
  },
  PersistentPostRun: func(cmd *cobra.Command, args []string) {
    if !reflect.DeepEqual(old_taskmap, taskmap) {
      slog.Debug("TaskMap has changed. Writing to file.", "old", old_taskmap, "new", taskmap)
      taskmap.PrintMarkdown(nil)
      if dry_run {
        taskmap.DryWrite(tasks_path)
      } else {
        taskmap.Write(tasks_path)
      }
    }
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
