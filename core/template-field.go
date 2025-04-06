package core

import (
	"bufio"
	"bytes"
	"text/template"

	"github.com/nimaaskarian/ydo/utils"
	"gopkg.in/yaml.v3"
)

type TemplateField[T any] interface {
  Resolve(*Task) error
  ToValue() T
  yaml.Marshaler
  yaml.Unmarshaler
}

func date(date TemplateDate, duration string) string {
  t, err := utils.ParseDue(duration, date.ToValue())
  if err != nil {
    return ""
  }
  return t.Format(DATE_PARSE_LAYOUT)
}

var funcs = template.FuncMap{"date": date }

// a base TemplateField with return type of string
type TemplateBase struct {
  template, resolved string
}

func (tb TemplateBase) GoString() string {
  return tb.template
}

func (tb TemplateBase) String() string {
  return tb.ToValue()
}

func NewTemplateBase(task string) TemplateBase {
  return TemplateBase { template: task }
}

func (tb *TemplateBase) Resolve(task *Task) error {
  tmpl := template.New("task-field").Funcs(funcs)
  var buffer bytes.Buffer
  tmpl, err := tmpl.Parse(tb.template)
  if err != nil {
    return err
  }
  writer := bufio.NewWriter(&buffer)
  err = tmpl.Execute(writer, task)
  writer.Flush()
  if err != nil {
    return err
  }
  tb.resolved = buffer.String()
  return nil
}

func (tb *TemplateBase) ToValue() string {
  if tb.resolved == "" {
    return tb.template
  }
  return tb.resolved
}

func (tb TemplateBase) MarshalYAML() (interface{}, error) {
  return tb.template, nil
}

func (tb TemplateBase) IsZero() bool {
  return tb.template == ""
}

func (tb *TemplateBase) UnmarshalYAML(node *yaml.Node) error {
  var raw string
  err := node.Decode(&raw)
  if err != nil {
    return err
  }
  tb.template = raw
  return nil
}
