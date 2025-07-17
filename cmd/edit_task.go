package cmd

import (
	"fmt"
	"log/slog"

	"github.com/nimaaskarian/ydo/core"
	"github.com/nimaaskarian/ydo/hooks"
	"github.com/nimaaskarian/ydo/utils"
)

type EditTask struct {
	flagTask core.Task
	newKey string
	removeDeps bool
	removeDepsTo bool
	keyRegen bool
	autoComplete bool
	newTitle string
}

func (et * EditTask) ApplyEdit(tm core.TaskMap, editKey string) error {
	task, err := taskmap.GetTask(editKey)
	if err != nil {
		return err
	}
	if due != "" {
		if due_date, err := utils.ParseDue(due, now); err == nil {
			task.Due = core.NewTemplateDate(due_date)
		}
	}
	if et.newTitle != "" {
		task.Task = core.NewTemplateBase(et.newTitle)
		events = append(events, hooks.Event{
			Type:      hooks.Edit,
			Key:       editKey,
			Secondary: et.newTitle,
		})
	}
	if !et.flagTask.Description.IsZero() {
		task.Description = et.flagTask.Description
	}
	if et.keyRegen {
		et.newKey = taskmap.TfidfNextKey(task.Task.Value(), config.Tfidf, editKey)
		if et.newKey != editKey {
			events = append(events, hooks.Event{Type: hooks.UpdateKey, Key: et.newKey, Secondary: editKey})
		}
	}
	for _, dep := range task.Deps {
		if _, err := taskmap.GetTask(dep); err != nil {
			return err
		}
	}
	task.Recur = et.flagTask.Recur

	taskDateFields := task.DateFields()
	for i, item := range et.flagTask.DateFields() {
		if !item.IsZero() {
			*taskDateFields[i] = *item
		}
	}
	if err := checkTask(task); err != nil {
		return err
	}

	if et.removeDeps {
		task.Deps = make([]string, 0)
	}
	if len(et.flagTask.Tags) > 0 {
		task.Tags = et.flagTask.Tags
	}
	if et.removeDeps {
		taskmap.WipeDependenciesToKey(editKey)
	}
	for _, dep_key := range dep_tos {
		task, event, err := taskmap.AddDep(dep_key, editKey)
		if err != nil {
			return err
		}
		taskmap[dep_key] = task
		events = append(events, event)
	}
	event, editKey := taskmap.ReplaceKeyInDeps(editKey, et.newKey)
	events = append(events, event)
	if et.autoComplete {
		task.AutoComplete = !task.AutoComplete
		events = append(events, hooks.Event{
			Type: hooks.Edit,
			Key:  fmt.Sprintf("toggle auto complete of %q", editKey),
		})
	}
	task.Deps = append(task.Deps, et.flagTask.Deps...)
	for _, dep := range et.flagTask.Deps {
		events = append(events, hooks.Event{
			Type:      hooks.AddDep,
			Key:       editKey,
			Secondary: dep,
		})
	}
	taskmap[editKey] = task
	slog.Info("Task edited", "task", task)
	return nil
}
