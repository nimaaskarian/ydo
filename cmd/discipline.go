package cmd

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fatih/color"
	"github.com/guptarohit/asciigraph"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

var height int
var thisweek bool
var thisday bool
var hourly bool
var start_of_the_week string

func init() {
	rootCmd.AddCommand(disciplineCmd)
	disciplineCmd.Flags().IntVarP(&height, "height", "H", 5, "specify height for the discipline graph")
	disciplineCmd.Flags().BoolVarP(&thisweek, "this-week", "w", false, "draw discipline for this week")
	disciplineCmd.Flags().BoolVarP(&thisday, "this-day", "d", false, "draw discipline for this day")
	disciplineCmd.Flags().BoolVar(&hourly, "hourly", false, "calculate discipline hourly, instead of daily")
	disciplineCmd.Flags().StringVar(&start_of_the_week, "start-of-the-week", "", "calculate discipline hourly, instead of daily")
	disciplineCmd.MarkFlagsMutuallyExclusive("this-day", "this-week")
}

var disciplineCmd = &cobra.Command{
	Use:               "discipline [start] [end]",
	Short:             "get a discipline graph",
	Args:              cobra.MaximumNArgs(2),
	ValidArgsFunction: DueCompletion,
	Long:              "get a graph of discipline from the start (which defaults to first task created) to end (which defaults to today)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var start, end time.Time
		var err error
    var first_weekday = time.Sunday
    if start_of_the_week != "" {
      if !thisweek {
        return errors.New(`flag "start-of-the-week" must be used with "this-week"`)
      }
      first_weekday, err = utils.ParseWeekday(start_of_the_week)
      if err != nil {
        return err
      }
    }
		if len(args) >= 1 {

			start, err = utils.ParseDue(args[0], now)
			if err != nil {
				return err
			}
		} else {
			if thisweek {
				start = utils.NaiveDate(now)
				for start.Weekday() != first_weekday {
					start = start.Add(-24 * time.Hour)
				}
			} else if thisday {
				start = utils.NaiveDate(now).AddDate(0, 0, -1)
			} else {
				min_created_at_key := slices.MinFunc(utils.Keys(taskmap), func(a, b string) int {
					return taskmap[a].CreatedAt.Compare(taskmap[b].CreatedAt)
				})
				start = utils.NaiveDate(taskmap[min_created_at_key].CreatedAt)
			}
		}
		if len(args) == 2 {
			end, err = utils.ParseDue(args[1], now)
			if err != nil {
				return err
			}
		} else {
			end = utils.NaiveDate(now).AddDate(0, 0, 1)
		}
		d := taskmap.TrackDiscipline(start, end, time.Hour*24)
		keys := utils.Keys(d)
		slices.SortFunc(keys, time.Time.Compare)

		data := taskmap.DisciplineSum(d)
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
