package core

type Types interface {
	String() string
}

type Primitive int

const (
	Int Primitive = iota
	Float
	String
	Bool
)

func (p Primitive) String() string {
	types := []string{
		Int:   "Int",
		Float: "Float",
		Bool:  "Bool",
	}
	return types[p]
}
