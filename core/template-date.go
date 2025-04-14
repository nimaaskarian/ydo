package core

import (
	"time"

	"github.com/nimaaskarian/ydo/utils"
	"gopkg.in/yaml.v3"
)

type TemplateDate struct {
	Base TemplateBase
	Date time.Time
}

const DATE_PARSE_LAYOUT = time.RFC3339

func date(date TemplateDate, duration string) string {
	if date.Date.IsZero() {
		return ""
	}
	t, err := utils.ParseDuration(duration, date.Date)
	if err != nil {
		return ""
	}
	return t.Format(DATE_PARSE_LAYOUT)
}

func (td *TemplateDate) Resolve(task *Task) error {
	err := td.Base.Resolve(task)
	if err != nil {
		return err
	}
	td.Date, err = time.Parse(DATE_PARSE_LAYOUT, td.Base.resolved)
	return err
}

func (td *TemplateDate) Value() time.Time {
	return td.Date
}

func (td TemplateDate) String() string {
	return td.Date.Format(DATE_PARSE_LAYOUT)
}

func NewTemplateDate(date time.Time) TemplateDate {
	td := TemplateDate{Date: date}
	if !td.Date.IsZero() {
		td.Base.Template = td.Date.Format(DATE_PARSE_LAYOUT)
	}
	return td
}

func TemplateDateFromTemplate(template string, task *Task) (TemplateDate, error) {
	td := TemplateDate{Base: TemplateBase{Template: template}}
	err := td.Resolve(task)
	return td, err
}

func (td TemplateDate) MarshalYAML() (any, error) {
	return td.Base.Template, nil
}

func (td TemplateDate) IsZero() bool {
	return td.Base.IsZero()
}

func (td *TemplateDate) UnmarshalYAML(node *yaml.Node) error {
	return td.Base.UnmarshalYAML(node)
}
