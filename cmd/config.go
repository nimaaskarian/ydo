package cmd

import (
	"path/filepath"
  "errors"
  "log/slog"
  "os"
	"github.com/nimaaskarian/ydo/core"

	"gopkg.in/yaml.v3"
)

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

func MarkdownFilter(md_config *core.MarkdownConfig) core.MarkdownFilter {
  switch md_config.Mode {
  case "todo":
    return core.Task.IsNotDone
  case "md":
    return nil
  default:
    if len(taskmap) >= md_config.Limit {
      return core.Task.IsNotDone
    } else {
      return nil
    }
  }
}
