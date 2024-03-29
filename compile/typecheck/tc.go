package typecheck

import (
	"github.com/x0y14/beat/compile/core"
)

type TFunction struct {
	Params  []*TVariable
	Returns []*TVariable
	*TypeTree
}

func NewFunction(params, returns []*TVariable) *TFunction {
	return &TFunction{
		Params:   params,
		Returns:  returns,
		TypeTree: &TypeTree{},
	}
}

type TVariable struct {
	Name string
	Type core.Type
}

func NewVariable(name string, typ core.Type) *TVariable {
	return &TVariable{
		Name: name,
		Type: typ,
	}
}

type TypeTree struct {
	// Nest 0がグローバル, 1が関数, 2が関数内のForとか, 3がForの中のIFとか, ...
	Parent *TypeTree
	F      map[string]*TFunction
	V      map[string]*TVariable
	Lower  []*TypeTree
}

func NewTypeTree() *TypeTree {
	return &TypeTree{
		Parent: nil,
		F:      map[string]*TFunction{},
		V:      map[string]*TVariable{},
		Lower:  []*TypeTree{},
	}
}

type TypeChecker struct {
	curtNest int       // 深さ
	curtTree *TypeTree // 現在参照してるやつ
	Tree     *TypeTree // 全体
}

func NewTypeChecker() *TypeChecker {
	t := NewTypeTree()
	return &TypeChecker{
		curtNest: 0,
		Tree:     t,
		curtTree: t,
	}
}

func (tc *TypeChecker) GoGlobal() {
	tc.curtTree = tc.Tree
	tc.curtNest = 0
}

func (tc *TypeChecker) GoParent() {
	tc.curtTree = tc.curtTree.Parent
}

// GoShallower ネストを浅く, {{}} -> {}
func (tc *TypeChecker) GoShallower() {
	tc.curtNest--
}

// GoDeeper ネストを深く, {} -> {{}}
func (tc *TypeChecker) GoDeeper() {
	tc.curtNest++
}

func (tc *TypeChecker) SetFunction(name string, f *TFunction, focus bool) {
	// 親の設定
	f.Parent = tc.curtTree
	// 名前対応表に追加
	tc.curtTree.F[name] = f
	// 注目する場合は差し替える
	if focus {
		tc.curtTree = f.TypeTree
		tc.curtNest++
	}
}

func (tc *TypeChecker) SetVariable(name string, v *TVariable) {
	tc.curtTree.V[name] = v
}

func (tc *TypeChecker) FindFunctionInCurrent(name string, focus bool) (*TFunction, bool) {
	// todo: test
	for fName, f := range tc.curtTree.F {
		if fName == name {
			if focus {
				tc.curtTree = f.TypeTree
			}
			return f, true
		}
	}
	return nil, false
}

func (tc *TypeChecker) FindFunctionConsiderNest(name string, focus bool) (*TFunction, bool) {
	// mainから検索を書けたなら、mainになかった場合、グローバルへ探索範囲を広げる
	// todo: test
	f, ok := tc.FindFunctionInCurrent(name, focus)
	if ok {
		return f, ok
	}
	tmp := tc.curtTree
	for n := tc.curtNest - 1; 0 <= n; n-- {
		tc.curtTree = tc.curtTree.Parent
		f, ok = tc.FindFunctionInCurrent(name, focus)
		if ok {
			tc.curtNest = n // 特にここバグになりやすそう
			tc.curtTree = f.TypeTree
			return f, ok
		}
	}
	tc.curtTree = tmp
	return nil, false
}

func (tc *TypeChecker) FindVariableInCurrent(name string) (*TVariable, bool) {
	for _, v := range tc.curtTree.V {
		if v.Name == name {
			return v, true
		}
	}
	return nil, false
}

func (tc *TypeChecker) FindVariableConsiderNest(name string) (*TVariable, bool) {
	// 現在の深さで探してみる
	v, ok := tc.FindVariableInCurrent(name)
	if ok {
		return v, ok
	}
	// なかったので浅瀬へ進みながら探す
	// 親を参照して行くので検索終了後復元
	tmp := tc.curtTree
	// 1までが関数内部, 現在から1つ引いた数分まで上昇(浅瀬へ)
	for n := tc.curtNest - 1; 0 <= n; n-- {
		tc.curtTree = tc.curtTree.Parent
		v, ok = tc.FindVariableInCurrent(name)
		if ok {
			return v, ok
		}
	}
	// 復元
	tc.curtTree = tmp
	return nil, false
}
