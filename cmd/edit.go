package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/fatih/color"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/nimaaskarian/ydo/utils"
	"github.com/spf13/cobra"
)

// edit flags
var (
	distance_threshold float64
	editTask EditTask
)

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.AddCommand(interactiveEditCmd)
	interactiveEditCmd.Flags().Float64VarP(&distance_threshold, "threshold", "t", 0.3, "threshold of difference to match a line as edited vs deleted")
	editCmd.Flags().BoolVarP(&editTask.keyRegen, "key-regen", "K", false, "regen key using the automatic next key generator (respects the config file)")
	editCmd.Flags().BoolVarP(&editTask.flagTask.AutoComplete, "auto-complete", "a", false, "toggle auto complete for the task")
	editCmd.Flags().BoolVar(&editTask.removeDepsTo, "remove-dep-to", false, "remove previous 'dependent to' for the task. using this with --dep-to causes to replace 'dependent to's")
	editCmd.Flags().StringVarP(&editTask.flagTask.Description.Template, "description", "e", "", "new description of the task")
	editCmd.Flags().BoolVar(&editTask.removeDeps, "remove-deps", false, "remove previous dependencies for the task. using this with --deps causes to replace dependencies")
	editCmd.Flags().StringVarP(&editTask.newKey, "key", "k", "", "new key to the task")
	editCmd.RegisterFlagCompletionFunc("key", KeyCompletion)

	editCmd.Flags().StringArrayVarP(&editTask.flagTask.Tags, "tag", "T", []string{}, "tag(s) for the task")
	editCmd.RegisterFlagCompletionFunc("tag", TagCompletion)

	editCmd.Flags().StringArrayVarP(&editTask.flagTask.Deps, "deps", "d", []string{}, "append dependencies for the task")
	editCmd.RegisterFlagCompletionFunc("deps", TaskKeyCompletionFilter(nil))

	editCmd.Flags().StringVarP(&editTask.flagTask.Until.Base.Template, "until", "U", "", "specify until (task is ignored after that date) for the tasks to print")
	editCmd.RegisterFlagCompletionFunc("until", DueCompletion)

	editCmd.Flags().StringVarP(&editTask.flagTask.Due.Base.Template, "due", "u", "", "specify due for the tasks to print")
	editCmd.RegisterFlagCompletionFunc("due", DueCompletion)

	editCmd.Flags().StringArrayVarP(&dep_tos, "dep-to", "D", []string{}, "append task keys for this task to be dependent to")
	editCmd.RegisterFlagCompletionFunc("dep-to", TaskKeyCompletionFilter(nil))

	editCmd.Flags().StringVarP(&editTask.flagTask.Schedule.Base.Template, "schedule", "s", "", "specify schedule for the tasks to print")
	editCmd.RegisterFlagCompletionFunc("schedule", DueCompletion)

	editCmd.Flags().StringVarP(&editTask.flagTask.Recur, "recur", "r", "", "duration of in which the ask recurs")
	editCmd.RegisterFlagCompletionFunc("recur", DurationCompletion)

	editCmd.ValidArgsFunction = TaskKeyCompletionOnFirst

}

var editCmd = &cobra.Command{
	Aliases: []string{"e"},
	Use:     "edit [key (optional)] [new task message (optional)]",
	Short:   "edit a task, or open the data file in favorite editor",
	Long:    "edit a task provided a key. give no keys to open the yaml file in your favorite editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			c, err := utils.EditorCmd(tasks_path)
			if err != nil {
				return err
			}
			utils.CmdStdOs(c)
			return c.Run()
		}
		var err error
		editTask.newTitle, err = TaskTitleFromArgs(args[1:])
		if err != nil {
			return err
		}
		for _, editKey := range taskmap.RegexpMatchingKeys(args[0], config.Regexp) {
			editTask.ApplyEdit(taskmap, editKey)
		}
		return nil
	},
	PostRunE: SaveChanges,
}

var interactiveEditCmd = &cobra.Command{
	Aliases: []string{"e"},
	Use:     "interactive",
	Short:   "interactively edit of your task list as markdown",
	Long:    "interactively edit of your task list as markdown in your EDITOR. add, edit or delete tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		temp, err := os.CreateTemp("", "ydo-interactive-edit-*.md")
		if err != nil {
			return err
		}
		config_copy := config.Markdown
		var before_buffer bytes.Buffer
		config_copy.File = io.Writer(&before_buffer)
		color.NoColor = true
		*config_copy.Beautify = false
		*config_copy.AdditionalInfo = false
		config_copy.Keys = make([]string, 0, len(taskmap))
		taskmap.PrintMarkdown(&config_copy, nil)
		before := before_buffer.Bytes()
		temp.Write(before)
		temp.Close()
		c, err := utils.EditorCmd(temp.Name())
		if err != nil {
			return err
		}
		utils.CmdStdOs(c)
		c.Run()
		after, err := os.ReadFile(temp.Name())
		if err != nil {
			return err
		}
		fmt.Println(config_copy.Keys)
		edits := myers.ComputeEdits(span.URIFromPath("before.md"), string(before), string(after))
		ud := UnifiedDiff{}
		ud.unified = gotextdiff.ToUnified("before.md", "after.md", string(before), edits)
		for ; ud.hunkIndex < len(ud.unified.Hunks); ud.hunkIndex += 1 {
			for ud.lineIndex < len(ud.lines()) {
				if ud.isLineEdit() {
					fmt.Println("haha this is an edit probably", ud.lines()[ud.lineIndex].Content)
				}
				// this adds to ud.lineIndex as well. yeah sorry this is a funciton that
				// does many things...
				ud.iterateAndSkipDeps(ud.handleDep)
			}
		}
		os.Remove(temp.Name())
		return nil
	},
	PostRunE: SaveChanges,
}

type UnifiedDiff struct {
	unified              gotextdiff.Unified
	hunkIndex, lineIndex int
}

func (ud *UnifiedDiff) LineNumber() int {
	return ud.hunk().FromLine + ud.lineIndex
}

func (ud *UnifiedDiff) isLineEdit() bool {
	i := ud.lineIndex
	lines := ud.lines()

	if i != 0 && lines[i].Kind == gotextdiff.Insert &&
		lines[i-1].Kind == gotextdiff.Delete {
		distance := levenshtein.ComputeDistance(lines[i].Content, lines[i-1].Content)
		max_len := float64(max(len(string(lines[i].Content)), len(lines[i-1].Content)))
		var normalized_distance float64 = float64(distance) / max_len
		return normalized_distance < distance_threshold
	}
	return false
}

func (ud *UnifiedDiff) iterateAndSkipDeps(callback func(string)) {
	prefix_size := 0
	lines := ud.lines()

	j := ud.lineIndex + 1
	for ; j < len(lines); j++ {
		trimmed := strings.TrimLeft(lines[j].Content, "\t ")
		if trimmed == lines[j].Content {
			break
		}
		cur_prefix_size := len(lines[j].Content) - len(trimmed)
		if prefix_size == 0 {
			prefix_size = cur_prefix_size
		} else if prefix_size != cur_prefix_size {
			break
		}
		callback(trimmed)
	}
	ud.lineIndex = j
}

func (ud *UnifiedDiff) lines() []gotextdiff.Line {
	return ud.hunk().Lines
}

func (ud *UnifiedDiff) hunk() *gotextdiff.Hunk {
	return ud.unified.Hunks[ud.hunkIndex]
}

func (ud *UnifiedDiff) handleDep(dep string) {
	fmt.Print("dep :> ", dep)
}
