package typecheck

import (
	"github.com/stretchr/testify/assert"
	"github.com/x0y14/beat/core"
	"github.com/x0y14/beat/parse"
	"github.com/x0y14/beat/tokenize"
	"testing"
)

func M(a, b any) []any {
	return []any{a, b}
}

func TestTypeChecker_01(t *testing.T) {
	checker := NewTypeChecker()
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

	checker.Tree = global_
	checker.curtNest = 2 // global(0)->main(1)->for(2)
	checker.curtTree = forMain
	// forの中から
	assert.Equal(t, M(checker.FindVariableInCurrent("i")), []any{&Variable{"i", core.Int}, true})
	// forの中からは直接xは見えない(遡らなければいけない)
	assert.Equal(t, M(checker.FindVariableInCurrent("x")), []any{(*Variable)(nil), false})
	// 遡った
	assert.Equal(t, M(checker.FindVariableConsiderNest("x")), []any{&Variable{"x", core.Int}, true})
	// main(1)からfor(2)が見えてないことを確認
	checker.curtNest = 1
	checker.curtTree = main_.TypeTree
	assert.Equal(t, M(checker.FindVariableInCurrent("x")), []any{&Variable{"x", core.Int}, true})
	assert.Equal(t, M(checker.FindVariableInCurrent("i")), []any{(*Variable)(nil), false})
}

func TestTypeChecker_02(t *testing.T) {
	checker := NewTypeChecker()
	checker.SetFunction("f1", NewFunction(nil, nil), false)                                        // 関数f1をグローバルに作成
	checker.SetFunction("main", NewFunction(nil, nil), true)                                       // 関数mainをグローバルに作成 注目
	assert.Equal(t, []any{(*Function)(nil), false}, M(checker.FindFunctionInCurrent("f1", false))) // mainから動かずにf1を検索
	parented := NewFunction(nil, nil)
	parented.Parent = checker.Tree
	assert.Equal(t, []any{parented, true}, M(checker.FindFunctionConsiderNest("f1", false))) // mainから動いて検索
}

func TestTypeChecker_AppendLowerWithParentInfo(t *testing.T) {
	parent := NewTypeTree()
	child := NewTypeTree()
	parent.AppendLowerWithParentInfo(child)
	assert.Equal(t, parent, child.Parent)
}

func TestTypeChecker_check_01(t *testing.T) {
	code := `
func main() int {
	return 1
}
`
	tokens, err := tokenize.Tokenize(code)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := parse.Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := TypeCheck(nodes)
	if err != nil {
		t.Fatal(err)
	}

	expectTree := NewTypeTree()
	fMain := NewFunction(nil, []*Variable{NewVariable("", core.Int)})
	fMain.Parent = expectTree
	expectTree.F["main"] = fMain

	assert.Equal(t, expectTree, tree)
}

func TestTypeChecker_check_02(t *testing.T) {
	code := `
func main() int {
	var x int = 1
	return x
}
`
	tokens, err := tokenize.Tokenize(code)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := parse.Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := TypeCheck(nodes)
	if err != nil {
		t.Fatal(err)
	}

	expectTree := NewTypeTree()
	fMain := NewFunction(nil, []*Variable{NewVariable("", core.Int)})
	fMain.TypeTree = NewTypeTree()
	fMain.TypeTree.V["x"] = NewVariable("x", core.Int)
	fMain.Parent = expectTree
	expectTree.F["main"] = fMain
	assert.Equal(t, expectTree, tree)
}

func TestTypeChecker_check_03(t *testing.T) {
	code := `
func f() int {
	return 1
}
func main() int {
	var x int = f()
	return x
}
`
	tokens, err := tokenize.Tokenize(code)
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := parse.Parse(tokens)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := TypeCheck(nodes)
	if err != nil {
		t.Fatal(err)
	}

	expectTree := NewTypeTree()

	fF := NewFunction(nil, []*Variable{NewVariable("", core.Int)})
	fF.TypeTree = NewTypeTree()

	fMain := NewFunction(nil, []*Variable{NewVariable("", core.Int)})
	fMain.TypeTree = NewTypeTree()
	fMain.TypeTree.V["x"] = NewVariable("x", core.Int)

	fF.Parent = expectTree
	expectTree.F["f"] = fF
	fMain.Parent = expectTree
	expectTree.F["main"] = fMain
	assert.Equal(t, expectTree, tree)
}
