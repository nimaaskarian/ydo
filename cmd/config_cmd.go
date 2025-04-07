package cmd

import (
	"errors"
	"io/fs"
	"os"
	"text/template"

	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

// default config template string with a yaml format
const DEFAULT_CONFIG = `# list of files that ydo tries to go through to find your tasks.yaml file
# relative files are relative to your shells cwd when you run the program
files:
  - tasks.yaml
  - {{ . }}
# possible values: debug, info, warn, error (default)
loglevel: error
# possible values: always, never, auto (default)
color: auto
markdown:
  # number of spaces on each indent. 0 is replaced with 3 (default)
  indent: 0
  # show description in markdown output
  description: false
  # true or false. by default gets enabled when the color is enabled
  beautify: true
  # possible values: todo, md, smart (default)
  mode: smart
tfidf:
  enabled: false
  # minimum tasks needed to turn on tfidf (enabled has to be true for this to work)
  min-task-count: 0
`

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [tasks (optional)]",
	Short: "edit ydo's config file",
	Long:  "edit ydo's config file in your favorite EDITOR. Writes the default config if it founds no config",
	RunE: func(cmd *cobra.Command, keys []string) error {
		if err := ensureConfigFile(); err != nil {
			return err
		}
		c, err := utils.EditorCmd(config_path)
		utils.CmdStdOs(c)

		if err != nil {
			return err
		}
		return c.Run()
	},
}

func ensureConfigFile() error {
	if _, err := os.Stat(config_path); errors.Is(err, fs.ErrNotExist) {
		f, err := os.Create(config_path)
		if err != nil {
			return err
		}
		defer f.Close()

		tmpl, err := template.New("config").Parse(DEFAULT_CONFIG)
		if err != nil {
			panic(err)
		}
		tmpl.Execute(f, tasks_path)
	} else if err != nil { // some other error occurred
		return err
	}
	return nil
}
