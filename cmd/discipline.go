package cmd

import (
	"fmt"
	"slices"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var height int

func init() {
  rootCmd.AddCommand(disciplineCmd)
  disciplineCmd.Flags().IntVarP(&height, "height", "H", 10, "specify height for the discipline graph")

  disciplineCmd.ValidArgsFunction = TaskKeyCompletionFilter(func(t core.Task, tm core.TaskMap) bool {return !t.AutoComplete && !t.IsDone(tm) })
}

var disciplineCmd = &cobra.Command{
  Use: "discipline [start] [end]",
  Short: "get a discipline graph from start (defaults to first task created) to end (defaults to today)",
  RunE: func(cmd *cobra.Command, args []string) error {
    var start, end time.Time
    var err error
    if len(args) >= 1 {
      start, err = time.ParseInLocation("2006-01-02", args[0], time.Local)
      if err != nil {
        return err
      }
    } else {
      min_created_at_key := slices.MinFunc(utils.Keys(taskmap), func(a, b string) int {
        return taskmap[a].CreatedAt.Compare(taskmap[b].CreatedAt)
      })
      start = taskmap[min_created_at_key].CreatedAt
    }
    if len(args) == 2 {
      end, err = time.ParseInLocation("2006-01-02", args[1], time.Local)
      if err != nil {
        return err
      }
    } else {
      end = time.Now()
    }
    data := taskmap.TrackDisciplineDaily(start, end)
    
    graph := asciigraph.Plot(
      data, asciigraph.Precision(3),
      asciigraph.Height(height),
      asciigraph.SeriesLegends("Discipline"),
      asciigraph.SeriesColors(asciigraph.Blue),
      )

    fmt.Println(graph)
    return nil
  },
  PostRunE: SaveChanges,
  PreRun: UpdateOldTaskMap,
}
