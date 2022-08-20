package resource

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/persistentvolumeclaim.yaml
var persitentVolumeClaim string

var _ Template = &PersistentVolumeClaimTemplate{}

type PersistentVolumeClaimTemplate struct {
	Name string
	Opts *PersistentVolumeClaimTemplateOptions
}

type PersistentVolumeClaimTemplateOptions struct {
	Size string
}

type PersistentVolumeClaimTemplateOption func(options *PersistentVolumeClaimTemplateOptions)

func NewPersitentVolumeClaimTemplateOptions(setters ...PersistentVolumeClaimTemplateOption) *PersistentVolumeClaimTemplateOptions {
	opts := &PersistentVolumeClaimTemplateOptions{}
	for _, setter := range setters {
		setter(opts)
	}
	return opts
}

func PersistentVolumeClaimSize(size string) PersistentVolumeClaimTemplateOption {
	return func(opts *PersistentVolumeClaimTemplateOptions) {
		opts.Size = size
	}
}

func NewPersistentVolumeClaimTemplate(name string, opts map[string]string) (*PersistentVolumeClaimTemplate, error) {
	setters, err := ParsePersistentVolumeClaimMapOpts(opts)
	if err != nil {
		return nil, err
	}

	return &PersistentVolumeClaimTemplate{
		Name: name,
		Opts: NewPersitentVolumeClaimTemplateOptions(setters...),
	}, nil
}

func (t *PersistentVolumeClaimTemplate) Generate() (string, error) {
	tmpl, err := template.New("persistentvolumeclaim").Parse(persitentVolumeClaim)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, t); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ParsePersistentVolumeClaimMapOpts(opts map[string]string) ([]PersistentVolumeClaimTemplateOption, error) {
	var setters []PersistentVolumeClaimTemplateOption

	for key, val := range opts {
		switch key {
		case "size":
			setters = append(setters, PersistentVolumeClaimSize(val))
		default:
			fmt.Printf("%s option is not supported\n", key)
		}
	}
	return setters, nil
}
