package generators

type Kind interface {
	String() string
}

type Generator interface {
	Generate() (string, error)
}

type GenSchema struct {
	kind Kind
}
