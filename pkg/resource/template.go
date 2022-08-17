package resource

import (
	"os"
	"path"
)

func TemplateDir() string {
	home, _ := os.UserHomeDir()
	return path.Join(home, ".kube/templates")
}

type Template interface {
	Generate() (string, error)
}
