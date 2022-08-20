package resource

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/cronjob.yaml
var cronjob string

var _ Template = &DeploymentTemplate{}

type CronjobTemplate struct {
	Name string
	Opts *CronjobTemplateOptions
}

type CronjobTemplateOptions struct {
	Image    string
	Schedule string
}

type CronjobTemplateOption func(options *CronjobTemplateOptions)

func NewCronjobTemplateOptions(setters ...CronjobTemplateOption) *CronjobTemplateOptions {
	opts := &CronjobTemplateOptions{}
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}

func CronjobImageName(name string) CronjobTemplateOption {
	return func(opts *CronjobTemplateOptions) {
		opts.Image = name
	}
}

func CronjobSchedule(schedule string) CronjobTemplateOption {
	return func(opts *CronjobTemplateOptions) {
		opts.Schedule = schedule
	}
}

func NewCronjobTemplate(name string, opts map[string]string) (*CronjobTemplate, error) {
	setters, err := ParseCronjobMapOpts(opts)
	if err != nil {
		return nil, err
	}

	return &CronjobTemplate{
		Name: name,
		Opts: NewCronjobTemplateOptions(setters...),
	}, nil
}

func (t *CronjobTemplate) Generate() (string, error) {
	tmpl, err := template.New("cronjob").Parse(cronjob)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, t); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ParseCronjobMapOpts(opts map[string]string) ([]CronjobTemplateOption, error) {
	var setters []CronjobTemplateOption

	for key, val := range opts {
		switch key {
		case "image":
			setters = append(setters, CronjobImageName(val))
		case "schedule":
			setters = append(setters, CronjobSchedule(val))
		default:
			fmt.Printf("%s option is not supported\n", key)
		}
	}
	return setters, nil
}
