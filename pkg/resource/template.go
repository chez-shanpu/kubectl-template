package resource

type Template interface {
	Generate() (string, error)
}
