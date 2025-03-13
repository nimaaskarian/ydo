package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
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

type Config struct  {
  // files to look for if --file option is not present
  Files []string `yaml:",omitempty"`
  LogLevel string `yaml:",omitempty"`
}

func (config *Config) ReadFile(path string) {
  content, _ := os.ReadFile(path)
  err := yaml.Unmarshal([]byte(content), config)
  if err != nil {
    slog.Error("Error reading config file", "err", err)
  }
  slog.Info("Config file loaded", "path", path)
}

func (config *Config) FirstFileAvailable() string {
  for _, file := range config.Files {
    if _, err := os.Stat(file); err == nil {
      return file
    }
  }
  slog.Warn("No tasks file available. Falling back to default")
  return filepath.Join(config_dir, "tasks.yaml")
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
  // flags
  key string
  tasks_path string
  config_path string
  // global state
  taskmap core.TaskMap
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
    slog.Info("Log level set", "loglevel", loglevel)
    if tasks_path == "" {
      tasks_path = config.FirstFileAvailable()
    }
    taskmap = core.LoadTaskMap(tasks_path)
    if taskmap == nil {
      taskmap = make(core.TaskMap)
    }
  },
  Run: func(cmd *cobra.Command, args []string) {
    if task, ok := taskmap[key]; ok {
      task.PrintMarkdown(taskmap, 1, []string{key}, nil)
    } else {
      taskmap.PrintMarkdown()
    }
  },
}
)

func init() {
  config_dir = utils.ConfigDir()
  rootCmd.PersistentFlags().StringVarP(&tasks_path, "file","f","", "path to tasks file")
  rootCmd.PersistentFlags().StringVarP(&config_path, "config","c",filepath.Join(config_dir, "config.yaml"), "path to config file")
  rootCmd.PersistentFlags().StringVarP(&key, "key","k", "", "task key")
  rootCmd.RegisterFlagCompletionFunc("key", TaskKeyCompletion)
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
