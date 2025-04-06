package core

import (
	"time"

	"gopkg.in/yaml.v3"
)

type TemplateDate struct {
  tb TemplateBase
  date time.Time
}

const DATE_PARSE_LAYOUT = "2006-01-02T15:04:05.999999999-07:00"
func (td * TemplateDate) Resolve(task *Task) error {
  err := td.tb.Resolve(task)
  if err != nil {
    return err
  }
  td.date, err = time.Parse(DATE_PARSE_LAYOUT, td.tb.resolved) 
  return err
}

func (td *TemplateDate) ToValue() time.Time {
  return td.date
}

func (td TemplateDate) String() string {
  return td.date.Format(DATE_PARSE_LAYOUT)
}

func NewTemplateDate(date time.Time) TemplateDate {
  td := TemplateDate { date: date }
  td.tb.template = td.date.Format(DATE_PARSE_LAYOUT)
  return td
}

func TemplateDateFromTemplate(template string, task *Task) (TemplateDate, error) {
  td := TemplateDate { tb: TemplateBase { template: template } }
  err := td.Resolve(task)
  return td, err
}


func (td TemplateDate) MarshalYAML() (any, error) {
  return td.tb.template, nil
}

func (td TemplateDate) IsZero() bool {
  return td.tb.IsZero()
}

func (td *TemplateDate) UnmarshalYAML(node *yaml.Node) error {
  var raw string
  err := node.Decode(&raw)
  if err != nil {
    return err
  }
  td.tb.template = raw
  return nil
}
