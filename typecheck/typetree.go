package typecheck

import (
	"github.com/x0y14/beat/core"
)

type Function struct {
	Params  []*Variable
	Returns []*Variable
	*TypeTree
}

func NewFunction(params, returns []*Variable) *Function {
	return &Function{
		Params:   params,
		Returns:  returns,
		TypeTree: NewTypeTree(),
	}
}

type Variable struct {
	Name string
	Type core.Type
}

func NewVariable(name string, typ core.Type) *Variable {
	return &Variable{
		Name: name,
		Type: typ,
	}
}

type TypeTree struct {
	// Nest 0がグローバル, 1が関数, 2が関数内のForとか, 3がForの中のIFとか, ...
	Parent *TypeTree
	F      map[string]*Function
	V      map[string]*Variable
	Lower  []*TypeTree
}

func NewTypeTree() *TypeTree {
	return &TypeTree{
		Parent: nil,
		F:      map[string]*Function{},
		V:      map[string]*Variable{},
		Lower:  []*TypeTree{},
	}
}

func (tt *TypeTree) AppendLowerWithParentInfo(t *TypeTree) {
	t.Parent = tt
	tt.Lower = append(tt.Lower, t)
}
