package core

import "fmt"

type Type interface {
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

func GetDataTypeByIdent(ident string) (Primitive, error) {
	switch ident {
	case "int":
		return Int, nil
	case "float":
		return Float, nil
	case "string":
		return String, nil
	case "bool":
		return Bool, nil
	}
	return 0, fmt.Errorf("unregistered type: %s", ident)
}
