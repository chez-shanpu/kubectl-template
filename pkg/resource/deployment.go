package resource

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"text/template"
)

//go:embed templates/deployment.yaml
var deployment string

var _ Template = &DeploymentTemplate{}

type DeploymentTemplate struct {
	Name string
	Opts *DeploymentTemplateOptions
}

type DeploymentTemplateOptions struct {
	Image    string
	Replicas int64
}

type DeploymentTemplateOption func(options *DeploymentTemplateOptions)

func NewDeploymentTemplateOptions(setters ...DeploymentTemplateOption) *DeploymentTemplateOptions {
	opts := &DeploymentTemplateOptions{}
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}

func DeploymentImageName(name string) DeploymentTemplateOption {
	return func(opts *DeploymentTemplateOptions) {
		opts.Image = name
	}
}

func DeploymentReplicas(replicas int64) DeploymentTemplateOption {
	return func(opts *DeploymentTemplateOptions) {
		opts.Replicas = replicas
	}
}

func NewDeploymentTemplate(name string, opts map[string]string) (*DeploymentTemplate, error) {
	setters, err := ParseMapOpts(opts)
	if err != nil {
		return nil, err
	}

	return &DeploymentTemplate{
		Name: name,
		Opts: NewDeploymentTemplateOptions(setters...),
	}, nil
}

func (t *DeploymentTemplate) Generate() (string, error) {
	tmpl, err := template.New("deployment").Parse(deployment)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, t); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ParseMapOpts(opts map[string]string) ([]DeploymentTemplateOption, error) {
	var setters []DeploymentTemplateOption

	for key, val := range opts {
		switch key {
		case "image":
			setters = append(setters, DeploymentImageName(val))
		case "replicas":
			replicas, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			setters = append(setters, DeploymentReplicas(replicas))
		default:
			fmt.Printf("%s option is not supported\n", key)
		}
	}
	return setters, nil
}
