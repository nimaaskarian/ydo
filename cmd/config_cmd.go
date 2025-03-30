package cmd

import (
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
  Use: "config [tasks (optional)]",
  Short: "edit ydo's config file",
  Long: "edit ydo's config file in your favorite EDITOR",
  RunE: func(cmd *cobra.Command, keys []string) error {
    c, err := utils.EditorCmd(config_path)
    utils.CmdStdOs(c)

    if err != nil {
      return err
    }
    return c.Run()
  },
}
