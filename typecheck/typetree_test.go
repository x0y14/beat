package typecheck

import (
	"github.com/stretchr/testify/assert"
	"github.com/x0y14/beat/core"
	"testing"
)

func M(a, b any) []any {
	return []any{a, b}
}

func TestTypeTree_01(t *testing.T) {
	tc := NewTypeChecker()
	var main_ *Function
	var forMain *TypeTree
	global_ := &TypeTree{
		Parent: nil,
		F: map[string]*Function{
			"main": main_,
		},
		V:     nil,
		Lower: nil,
	}
	main_ = &Function{
		Params:  nil,
		Returns: nil,
		TypeTree: &TypeTree{
			Parent: nil,
			F:      nil,
			V: map[string]*Variable{
				"x": {
					Name: "x",
					Type: core.Int,
				},
			},
			Lower: []*TypeTree{forMain},
		},
	}
	forMain = &TypeTree{
		Parent: main_.TypeTree,
		F:      nil,
		V: map[string]*Variable{
			"i": {
				Name: "i",
				Type: core.Int,
			},
		},
		Lower: nil,
	}

	tc.Tree = global_
	tc.curtNest = 2 // global(0)->main(1)->for(2)
	tc.curtTree = forMain
	// forの中から
	assert.Equal(t, M(tc.FindVariableInCurrent("i")), []any{&Variable{"i", core.Int}, true})
	// forの中からは直接xは見えない(遡らなければいけない)
	assert.Equal(t, M(tc.FindVariableInCurrent("x")), []any{(*Variable)(nil), false})
	// 遡った
	assert.Equal(t, M(tc.FindVariableConsiderNest("x")), []any{&Variable{"x", core.Int}, true})
	// main(1)からfor(2)が見えてないことを確認
	tc.curtNest = 1
	tc.curtTree = main_.TypeTree
	assert.Equal(t, M(tc.FindVariableInCurrent("x")), []any{&Variable{"x", core.Int}, true})
	assert.Equal(t, M(tc.FindVariableInCurrent("i")), []any{(*Variable)(nil), false})
}

func TestTypeChecker_02(t *testing.T) {
	tc := NewTypeChecker()
	tc.SetFunction("f1", NewFunction(nil, nil), false)                                        // 関数f1をグローバルに作成
	tc.SetFunction("main", NewFunction(nil, nil), true)                                       // 関数mainをグローバルに作成 注目
	assert.Equal(t, []any{(*Function)(nil), false}, M(tc.FindFunctionInCurrent("f1", false))) // mainから動かずにf1を検索
	parented := NewFunction(nil, nil)
	parented.Parent = tc.Tree
	assert.Equal(t, []any{parented, true}, M(tc.FindFunctionConsiderNest("f1", false))) // mainから動いて検索
}

func TestTypeTree_AppendLowerWithParentInfo(t *testing.T) {
	parent := NewTypeTree()
	child := NewTypeTree()
	parent.AppendLowerWithParentInfo(child)
	assert.Equal(t, parent, child.Parent)
}
