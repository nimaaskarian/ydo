package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/nimaaskarian/ydo/hooks"
	"github.com/nimaaskarian/ydo/utils"
)

type Task struct {
	Task          TemplateBase `yaml:",omitempty"`
	Description   TemplateBase `yaml:",omitempty"`
	Deps          []string     `yaml:",omitempty,flow"`
	Done          bool         `yaml:",omitempty"`
	AutoComplete  bool         `yaml:"auto-complete,omitempty"`
	CreatedAt     time.Time    `yaml:"created-at,omitempty"`
	Due           TemplateDate `yaml:",omitempty"`
	Until         TemplateDate `yaml:",omitempty"`
	Schedule      TemplateDate `yaml:",omitempty"`
	DoneAt        time.Time    `yaml:"done-at,omitempty"`
	DoneAtArchive []time.Time  `yaml:"done-at-archive,omitempty"`
	Recur         string       `yaml:",omitempty"`
	Tags          []string     `yaml:",omitempty"`
}

func (task *Task) DateFields() [3]*TemplateDate {
	return [...]*TemplateDate{
		&task.Due,
		&task.Until,
		&task.Schedule,
	}
}

func (task *Task) ResolveTemplates(now time.Time) {
	date_fields := task.DateFields()
	for _, item := range date_fields {
		date, err := time.Parse(DATE_PARSE_LAYOUT, item.Base.resolved)
		if err == nil {
			*item = NewTemplateDate(date)
		}
	}
	for _, item := range date_fields {
		item.Resolve(task)
	}
	task.Task.Resolve(task)
	task.Description.Resolve(task)
}

func (task *Task) IsDone(taskmap TaskMap, now time.Time) bool {
	if task.AutoComplete {
		for _, key := range task.Deps {
			if !taskmap[key].IsDone(taskmap, now) {
				return false
			}
		}
		return true
	}
	if t, err := utils.ParseDuration(task.Recur, task.DoneAt); err == nil && !t.IsZero() && !now.Before(t) {
		task.Done = false
	}
	return task.Done && (task.DoneAt.IsZero() || !task.DoneAt.After(now))
}

func (task *Task) Do(taskmap TaskMap, now time.Time, force bool) error {
	if task.AutoComplete {
		return errors.New("Task is auto-completed")
	}
	if !force && task.IsDone(taskmap, now) {
		return errors.New("Task is already done")
	}
	task.Done = true
	if task.Recur != "" && !task.DoneAt.IsZero() {
		task.DoneAtArchive = append(task.DoneAtArchive, task.DoneAt)
	}
	task.DoneAt = now
	return nil
}

func (task *Task) Undo(taskmap TaskMap, now time.Time) {
	if task.IsDone(taskmap, now) && !task.AutoComplete {
		task.Done = false
		if task.Recur != "" && len(task.DoneAtArchive) > 0 {
			length := len(task.DoneAtArchive)
			task.DoneAt = task.DoneAtArchive[length-1]
			task.DoneAtArchive = task.DoneAtArchive[:length-1]
		} else {
			task.DoneAt = time.Time{}
		}
	}
}

func (task *Task) FindDoneAt(taskmap TaskMap) time.Time {
	if task.AutoComplete {
		max_doneat := time.Time{}
		for _, key := range task.Deps {
			task := taskmap[key]
			doneat := task.FindDoneAt(taskmap)
			if doneat.After(max_doneat) {
				max_doneat = doneat
			}
			return max_doneat
		}
	}
	return task.DoneAt
}

func (task *Task) IsNotDone(taskmap TaskMap, now time.Time) bool {
	return !task.IsDone(taskmap, now)
}

// delete() the key, run this.
// runs over dependencies listed within the task itself.
func (task *Task) CascadeOrphanDeps(taskmap TaskMap) (out hooks.Events) {
	for _, dep := range task.Deps {
		if !taskmap.HasKeyInDeps(dep) {
			delete(taskmap, dep)
			out = append(out, hooks.Event{Type: hooks.Delete, Key: dep})
		}
	}
	return
}

func (task *Task) IsDeleted(date time.Time) bool {
	return date.Before(task.CreatedAt) ||
		(!task.Until.Value().IsZero() && task.Until.Value().Before(date)) ||
		(!task.Schedule.Value().IsZero() && task.Schedule.Value().After(date))
}

func (task *Task) PrintMarkdown(taskmap TaskMap, depth uint, seen_keys map[string]bool, key string, config *MarkdownConfig, filter TaskFilter) (count int) {
	if task == nil ||
		(filter != nil && !filter(task, taskmap, config.Now)) ||
		task.IsDeleted(config.Now) {
		return 0
	}

	if config.Limit != 0 && len(seen_keys) >= config.Limit {
		return 1
	}

	config.PrintIndent(depth)
	if task.IsDone(taskmap, config.Now) {
		config.PrintDonePrefix()
		printKey(key, config)
		printDoneTask(task, taskmap, config)
	} else {
		config.PrintUndonePrefix()
		printKey(key, config)
		printPendingTask(task, taskmap, config)
	}
	printTags(task, config)
	fmt.Fprintln(config.File)
	printDescription(depth, task, config)
	if seen_keys != nil {
		if _, ok := seen_keys[key]; ok {
			return 0
		}
		seen_keys[key] = true
		if config.Keys != nil {
			config.Keys = append(config.Keys, key)
		}
	}
	for _, dep_key := range task.Deps {
		count += taskmap[dep_key].PrintMarkdown(taskmap, depth+1, seen_keys, dep_key, config, filter)
	}
	return 1 + count
}

func printDoneTask(task *Task, taskmap TaskMap, config *MarkdownConfig) {
	var recur string
	done_at := task.FindDoneAt(taskmap)
	fmt.Fprintf(config.File, "%s", task.Task)
	if *config.AdditionalInfo {
		if !done_at.IsZero() {
			overdue := ""
			if !task.Due.Value().IsZero() && done_at.After(task.Due.Value()) {
				overdue += ", " + utils.FormatDuration(done_at.Sub(task.Due.Value())) + " overdue"
			}
			if task.Recur != "" {
				recur = ", each " + task.Recur
			}
			fmt.Fprintf(config.File, "(%s ago%s%s)", utils.FormatDuration(config.Now.Sub(done_at)), overdue, recur)
		} else {
			if task.Recur != "" {
				recur = " (each " + task.Recur + ")"
			}
			fmt.Fprintf(config.File, " (each %s)", task.Recur)
		}
	}
}

func printPendingTask(task *Task, taskmap TaskMap, config *MarkdownConfig) {
	fmt.Fprintf(config.File, "%s", task.Task)
	if *config.AdditionalInfo {
		var recur string
		if task.Recur != "" {
			recur = " (each " + task.Recur
			done_at := task.FindDoneAt(taskmap)
			if !done_at.IsZero() {
				if date, err := utils.ParseDuration(task.Recur, done_at); err == nil && date.Before(config.Now) {
					recur += ", " + utils.FormatDuration(config.Now.Sub(date)) + " overdue"
				}
			}
			recur += ")"
		}
		due_print := ""
		if !task.Due.Value().IsZero() {
			diff := task.Due.Value().Sub(config.Now)
			due_print = " ("
			if task.Due.Value().Add(-diff).Compare(config.Now) != 0 {
				due_print += strconv.Itoa(task.Due.Value().Year()-config.Now.Year()) + "y"
			} else {
				if diff < 0 {
					due_print = " (-"
					diff = -diff
				}
				due_print += utils.FormatDuration(diff)
			}
			due_print += ")"
		}
		fmt.Fprintf(config.File, "%s%s", due_print, recur)
	}
}

var bold = color.New(color.Bold)

func printKey(key string, config *MarkdownConfig) {
	if key != "" {
		fmt.Fprintf(config.File, "%s: ", bold.Sprint(key))
	}
}

var underline = color.New(color.Underline)

func printTags(task *Task, config *MarkdownConfig) {
	for _, tag := range task.Tags {
		fmt.Fprint(config.File, " ")
		underline.Fprintf(config.File, "#%s", tag)
	}
}

type TaskFilter func(task *Task, taskmap TaskMap, now time.Time) bool

type MarkdownConfig struct {
	Indent         uint      `yaml:",omitempty"`
	Mode           string    `yaml:",omitempty"`
	Description    bool      `yaml:",omitempty"`
	Limit          int       `yaml:",omitempty"`
	Beautify       *bool     `yaml:",omitempty"`
	AdditionalInfo *bool     `yaml:"additional-info,omitempty"`
	File           io.Writer `yaml:"-"`
	Keys           []string  `yaml:"-"`
	Now            time.Time
}

func (config *MarkdownConfig) Init() {
	if config.Beautify == nil {
		beautify := !color.NoColor
		config.Beautify = &beautify
	}
	if config.AdditionalInfo == nil {
		additional_info := true
		config.AdditionalInfo = &additional_info
	}
	if config.File == nil {
		config.File = os.Stdout
	}
	if config.Indent == 0 {
		config.Indent = 3
	}
}
func (mc *MarkdownConfig) PrintIndent(depth uint) {
	for range depth * mc.Indent {
		fmt.Fprint(mc.File, " ")
	}
}

func (mc *MarkdownConfig) PrintDonePrefix() {
	mc.PrintListPrefix()
	if *mc.Beautify {
		fmt.Fprint(mc.File, "[✓] ")
	} else {
		fmt.Fprint(mc.File, "[x] ")
	}
}

func (mc *MarkdownConfig) PrintListPrefix() {
	if !*mc.Beautify {
		fmt.Fprint(mc.File, "-")
	}
	fmt.Fprint(mc.File, " ")
}

func (mc *MarkdownConfig) PrintUndonePrefix() {
	mc.PrintListPrefix()
	fmt.Fprint(mc.File, "[ ] ")
}

func printDescription(depth uint, task *Task, config *MarkdownConfig) {
	if config.Description && task.Description.Value() != "" {
		for line := range strings.Lines(task.Description.Value()) {
			config.PrintIndent(depth + 1)
			fmt.Fprint(config.File, line)
		}
		fmt.Fprintln(config.File)
	}
}
