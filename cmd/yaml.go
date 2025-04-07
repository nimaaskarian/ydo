package cmd

import (
	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(yamlCmd)
	yamlCmd.ValidArgsFunction = TaskKeyCompletionFilter(nil)
}

var yamlCmd = &cobra.Command{
	Aliases: []string{"y"},
	Use:     "yaml [tasks (optional)]",
	Short:   "output tasks as yaml",
	RunE: func(cmd *cobra.Command, keys []string) error {
		if len(keys) == 0 {
			core.PrintYaml(taskmap)
		} else {
			tmp_map := make(core.TaskMap, len(keys))
			for _, key := range keys {
				task, err := taskmap.GetTask(key)
				if err != nil {
					return err
				}
				tmp_map[key] = task
			}
			core.PrintYaml(tmp_map)
		}
		return nil
	},
}
