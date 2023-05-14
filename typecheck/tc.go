package typecheck

import "github.com/x0y14/beat/core"

type TFunction struct {
	Params  []*Variable
	Returns []*Variable
	*TypeTree
}
type TVariable struct {
	Name string
	Type core.Type
}
type TypeTree struct {
	// Nest 0がグローバル, 1が関数, 2が関数内のForとか, 3がForの中のIFとか, ...
	Parent *TypeTree
	F      map[string]*TFunction
	V      map[string]*TVariable
	Lower  []*TypeTree
}

type TypeChecker struct {
	curtNest int       // 深さ
	curtTree *TypeTree // 現在参照してるやつ
	Tree     *TypeTree // 全体
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		curtNest: 0,
		Tree: &TypeTree{
			Parent: nil,
			F:      map[string]*TFunction{},
			V:      map[string]*TVariable{},
			Lower:  []*TypeTree{},
		},
	}
}

// GoShallower ネストを浅く, {{}} -> {}
func (tc *TypeChecker) GoShallower() {
	tc.curtNest--
}

// GoDeeper ネストを深く, {} -> {{}}
func (tc *TypeChecker) GoDeeper() {
	tc.curtNest++
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
	for n := tc.curtNest - 1; 0 < n; n-- {
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
