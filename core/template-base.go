package core

import (
	"bufio"
	"bytes"
	"text/template"

	"gopkg.in/yaml.v3"
)

var funcs = template.FuncMap{"date": date }

// a base TemplateField with return type of string
// Template* structs have to have yaml.Marshaler and yaml.Unmarshaler
// interfaces.
type TemplateBase struct {
  Template, resolved string
}

func (tb TemplateBase) GoString() string {
  return tb.Template
}

func (tb TemplateBase) String() string {
  return tb.Value()
}

func NewTemplateBase(task string) TemplateBase {
  return TemplateBase { Template: task }
}

func (tb *TemplateBase) Resolve(task *Task) error {
  tmpl := template.New("task-field").Funcs(funcs)
  var buffer bytes.Buffer
  tmpl, err := tmpl.Parse(tb.Template)
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

func (tb *TemplateBase) Value() string {
  if tb.resolved == "" {
    return tb.Template
  }
  return tb.resolved
}

func (tb TemplateBase) MarshalYAML() (any, error) {
  return tb.Template, nil
}

func (tb TemplateBase) IsZero() bool {
  return tb.Template == ""
}

func (tb *TemplateBase) UnmarshalYAML(node *yaml.Node) error {
  var raw string
  err := node.Decode(&raw)
  if err != nil {
    return err
  }
  tb.Template = raw
  return nil
}
