package typecheck

import (
	"fmt"
	"github.com/x0y14/beat/core"
)

type Block struct {
	Nest  int
	Vars  []*Variable
	Lower map[int][]*Block
}

func (b *Block) FindVarById(id Identifier) (*Variable, bool) {
	for _, v := range b.Vars {
		if v.Id == id {
			return v, true
		}
	}
	return nil, false
}
func (b *Block) FindVarByName(name string) (*Variable, bool) {
	// todo
	for _, v := range b.Vars {
		if v.Name == name {
			return v, true
		}
	}
	return nil, false
}

func (b *Block) SetVar(name string, typ core.Type) (Identifier, error) {
	if _, ok := b.FindVarByName(name); ok {
		return "", fmt.Errorf("定義済み: %s", name)
	}
	id := GenerateId()
	b.Vars = append(b.Vars, &Variable{
		Id:   id,
		Name: name,
		Type: typ,
	})
	return id, nil
}

type Variable struct {
	Id   Identifier // 未確認を扱うため
	Name string
	Type core.Type
}

type Function struct {
	Id      Identifier // 未確認を扱うため
	Params  []*Variable
	Returns []*Variable
	*Block
}

//func (f *Function) FindVarByIdConsiderNest(id Identifier) (*Variable, bool) {
//	for _, v := range f.Vars {
//		if f.Id == id {
//			return v, true
//		}
//	}
//	return nil, false
//}

//func (f *Function) FindVarByNameConsiderNest(name string) (*Variable, bool) {
//	// 1, 2
//	//for i, b := range f.
//}

type TypeMap struct {
	Functions        map[string]*Function
	Variables        map[string]*Variable
	UnknownFunctions []*Function
	UnknownVariables []*Variable
}

func (t *TypeMap) FindFnById(id Identifier) (*Function, bool) {
	for _, f := range t.Functions {
		if f.Id == id {
			return f, true
		}
	}
	return nil, false
}

func (t *TypeMap) FindFnByName(name string) (*Function, bool) {
	f, ok := t.Functions[name]
	if !ok {
		return nil, true
	}
	return f, false
}

func (t *TypeMap) FindVarByName(name string) (*Variable, bool) {
	v, ok := t.Variables[name]
	if ok {
		return v, true
	}
	return nil, false
}

func (t *TypeMap) SetFn(name string, f *Function) {
	t.Functions[name] = f
}
