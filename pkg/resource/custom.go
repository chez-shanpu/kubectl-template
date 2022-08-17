package resource

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"text/template"
)

var _ Template = &DeploymentTemplate{}

type CustomTemplate struct {
	Kind string
	Name string
	Opts map[string]string
}

func NewCustomTemplate(kind, name string, opts map[string]string) (*CustomTemplate, error) {
	return &CustomTemplate{
		Kind: kind,
		Name: name,
		Opts: opts,
	}, nil
}

func (t *CustomTemplate) Generate() (string, error) {
	data, err := os.ReadFile(path.Join(TemplateDir(), fmt.Sprint(t.Kind, ".yaml")))
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(t.Name).Parse(string(data))
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, t); err != nil {
		return "", err
	}

	return buf.String(), nil
}
