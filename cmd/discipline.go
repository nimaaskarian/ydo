package cmd

import (
	"fmt"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/nimaaskarian/ydo/core"
	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(disciplineCmd)
  disciplineCmd.ValidArgsFunction = TaskKeyCompletionFilter(func(t core.Task, tm core.TaskMap) bool {return !t.AutoComplete && !t.IsDone(tm) })
}

var disciplineCmd = &cobra.Command{
  Use: "discipline start [end]",
  Short: "get a discipline graph from start to end (defaults to today)",
  Args: cobra.MinimumNArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
    start, err := time.ParseInLocation("2006-01-02", args[0], time.Local)
    if err != nil {
      return err
    }
    var end time.Time
    if len(args) == 2 {
      end, err = time.ParseInLocation("2006-01-02", args[1], time.Local)
      if err != nil {
        return err
      }
    } else {
      end = time.Now()
    }
    data := taskmap.TrackDisciplineDaily(start, end)
    fmt.Println(data)
    graph := asciigraph.Plot(data)

    fmt.Println(graph)
    return nil
  },
  PostRunE: SaveChanges,
  PreRun: UpdateOldTaskMap,
}
