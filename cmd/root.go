package cmd

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/hooks"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var (
	// global flags
	tasks_path   string
	config_path  string
	dry_run      bool
	always_yes   bool
	exact_match  bool
	regexp       bool
	now          time.Time
	now_str      string
	color_option string
	// global state
	taskmap core.TaskMap
	events  hooks.Events

	config_dir          string
	config              Config
	tasksMarkdownFilter core.TaskFilter

	rootCmd = &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "ydo",
		Short:         "ydo is a frictionless and fast to-do app",
		Long:          `Fast, featurefull and frictionless to-do app with a graph structure`,
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
			if color_option != "" {
				config.Color = color_option
			}
			config.Init()
			if exact_match {
				config.Regexp = false
			}
			if regexp {
				config.Regexp = true
			}
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
			taskmap = core.LoadTaskMap(tasks_path, now)
			tasksMarkdownFilter = makeTasksMarkdownFilter(&config.Markdown)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return taskmap.PrintMarkdown(&config.Markdown, tasksMarkdownFilter)
		},
	}
)

func SaveChanges(cmd *cobra.Command, _ []string) error {
	if !events.Empty() {
		slog.Debug("TaskMap has changed. Writing to file.", "events", events)
		if dry_run {
			taskmap.DryWrite(tasks_path)
		} else {
			taskmap.Write(tasks_path)
		}
		for _, hook := range config.Hooks {
			for i, arg := range hook {
				if arg == "{}" {
					hook[i] = events.String()
				}
			}
			cmd := exec.Command(hook[0], hook[1:]...)
			cmd.Dir = filepath.Dir(tasks_path)
			if err := cmd.Run(); err != nil {
				return err
			}
		}
		if err := taskmap.PrintMarkdown(&config.Markdown, tasksMarkdownFilter); err != nil {
			return err
		}
	}
	return nil
}

func interactiveHelper(name string, include_func func(*core.Task, core.TaskMap, time.Time) bool) (map[string]bool, error) {
	temp, err := os.CreateTemp("", name)
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, 0)
	for _, key := range taskmap.SortedKeys() {
		task := taskmap[key]
		if include_func == nil || include_func(task, taskmap, now) {
			temp.WriteString(key + "\n")
			out[key] = true
		}
	}
	temp_name := temp.Name()
	if err := temp.Close(); err != nil {
		return nil, err
	}
	c, err := utils.EditorCmd(temp_name)
	if err != nil {
		return nil, err
	}
	utils.CmdStdOs(c)
	if err := c.Run(); err != nil {
		return nil, err
	}
	content, err := os.ReadFile(temp_name)
	if err != nil {
		return nil, err
	}
	os.Remove(temp_name)
	for line := range bytes.Lines(content) {
		key := string(bytes.TrimSpace(line))
		if _, ok := out[key]; !ok {
			return nil, fmt.Errorf("Key not found: %q", key)
		}
		out[key] = false
	}
	return out, nil
}

func init() {
	config_dir = utils.ConfigDir()
	rootCmd.PersistentFlags().StringVarP(&tasks_path, "file", "f", "", "path to tasks file")
	rootCmd.PersistentFlags().StringVarP(&config_path, "config", "c", filepath.Join(config_dir, "config.yaml"), "path to config file")
	rootCmd.PersistentFlags().BoolVarP(&dry_run, "dry-run", "n", false, "perform a trial run with no changes made")
	rootCmd.PersistentFlags().BoolVarP(&always_yes, "always-yes", "Y", false, "answer yes to all the yes/no questions")
	rootCmd.PersistentFlags().BoolVar(&exact_match, "exact-match", false, "exact match instead of using regexp to match the keys")
	rootCmd.PersistentFlags().BoolVar(&regexp, "regexp", false, "regexp instead of using exact match to match the keys")
	rootCmd.MarkFlagsMutuallyExclusive("regexp", "exact-match")

	rootCmd.PersistentFlags().StringVarP(&now_str, "now", "N", "", "current time of operations (defaults to current system time)")
	rootCmd.RegisterFlagCompletionFunc("now", DueCompletion)
	rootCmd.PersistentFlags().StringVar(&color_option, "color", "", "when to print in color (defaults to auto, overrides config's markdown.color option)")
	rootCmd.RegisterFlagCompletionFunc("color", cobra.FixedCompletions([]string{"always", "never", "auto"}, cobra.ShellCompDirectiveNoFileComp))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
