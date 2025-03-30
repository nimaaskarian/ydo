package cmd

import (
	"fmt"
	"slices"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
  "github.com/fatih/color"
)

var height int

func init() {
  rootCmd.AddCommand(disciplineCmd)
  disciplineCmd.Flags().IntVarP(&height, "height", "H", 5, "specify height for the discipline graph")
}

var disciplineCmd = &cobra.Command{
  Use: "discipline [start] [end]",
  Short: "get a discipline graph",
  Long: " get a graph of discipline from the start (which defaults to first task created) to end (which defaults to today)",
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
      end = now
    }
    data := taskmap.TrackDisciplineDaily(start, end, now)
    discipline_color := asciigraph.Blue
    if color.NoColor {
      discipline_color = asciigraph.Default
    }
    
    graph := asciigraph.Plot(
      data, asciigraph.Precision(3),
      asciigraph.Height(height),
      asciigraph.SeriesLegends("Discipline"),
      asciigraph.SeriesColors(discipline_color),
      )

    fmt.Println(graph)
    return nil
  },
}
