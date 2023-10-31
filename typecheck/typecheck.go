package typecheck

import (
	"fmt"
	"github.com/x0y14/beat/core"
	"github.com/x0y14/beat/parse"
)

type TypeChecker struct {
	curtNest int       // 深さ
	curtTree *TypeTree // 現在参照してるやつ
	curtFn   *Function
	Tree     *TypeTree // 全体
}

func NewTypeChecker() *TypeChecker {
	t := NewTypeTree()
	return &TypeChecker{
		curtNest: 0,
		curtTree: t,
		curtFn:   nil,
		Tree:     t,
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

func (tc *TypeChecker) SetFunction(name string, f *Function, focus bool) {
	// 親の設定
	f.Parent = tc.curtTree
	// 名前対応表に追加
	tc.curtTree.F[name] = f
	// 注目する場合は差し替える
	if focus {
		tc.curtTree = f.TypeTree
		tc.curtFn = f
		tc.curtNest++
	}
}

func (tc *TypeChecker) SetVariable(name string, v *Variable) {
	tc.curtTree.V[name] = v
}

func (tc *TypeChecker) FindFunctionInCurrent(name string, focus bool) (*Function, bool) {
	// todo: test
	for fName, f := range tc.curtTree.F {
		if fName == name {
			if focus {
				tc.curtTree = f.TypeTree
				tc.curtFn = f
			}
			return f, true
		}
	}
	return nil, false
}

func (tc *TypeChecker) FindFunctionConsiderNest(name string, focus bool) (*Function, bool) {
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
			tc.curtFn = f
			return f, ok
		}
	}
	tc.curtTree = tmp
	return nil, false
}

func (tc *TypeChecker) FindVariableInCurrent(name string) (*Variable, bool) {
	for _, v := range tc.curtTree.V {
		if v.Name == name {
			return v, true
		}
	}
	return nil, false
}

func (tc *TypeChecker) FindVariableConsiderNest(name string) (*Variable, bool) {
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

var tc *TypeChecker

func TypeCheck(nodes []*parse.Node) (*TypeTree, error) {
	tc = NewTypeChecker()
	for _, node := range nodes {
		_, err := toplevel(node)
		if err != nil {
			return nil, err
		}
	}
	return tc.Tree, nil
}

func toplevel(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.FuncDef:
		return fn(node)
	case parse.VarDecl:
	default:
		return stmt(node)
	}
	return nil, fmt.Errorf("toplevel unimplemented")
}

func fn(node *parse.Node) ([]core.Type, error) {
	field := node.FuncDefField
	name := field.Identifier.IdentField.Ident

	// 引数
	var parameters []*Variable
	if field.Parameters != nil {
		for _, paramNode := range field.Parameters.PolynomialField.Values {
			param := paramNode.FuncParam
			parameters = append(parameters, NewVariable(
				param.Identifier.IdentField.Ident,
				param.DataType.DataTypeField.Type))
		}
	}

	// 戻り値
	var returns []*Variable
	if field.Returns != nil {
		for _, rtNode := range field.Returns.PolynomialField.Values {
			returns = append(returns, NewVariable("", rtNode.DataTypeField.Type))
		}
	}

	// グローバルにくっつけてあげる
	f := NewFunction(parameters, returns)
	tc.SetFunction(name, f, true)

	// 関数の中身を解析
	// 戻り値はReturnsにあるものを使うので捨てる
	_, err := stmt(field.Body)
	if err != nil {
		return nil, err
	}

	// グローバルに戻す
	tc.GoGlobal()
	return nil, nil
}

func stmt(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.Return:
		var rtTypes []core.Type
		for _, val := range node.PolynomialField.Values {
			t, err := expr(val)
			if err != nil {
				return nil, err
			}
			rtTypes = append(rtTypes, t...)
		}
		if !isSameType(variablesToTypes(tc.curtFn.Returns), rtTypes) {
			return nil, fmt.Errorf("戻り値の型が異なる: %v, %v", variablesToTypes(tc.curtFn.Returns), rtTypes)
		}
		return rtTypes, nil
	case parse.IfElse:
	case parse.For:
	case parse.Block:
		for _, statement := range node.BlockField.Statements {
			_, err := stmt(statement)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	default:
		return expr(node)
	}
	return nil, fmt.Errorf("stmt unimplemented")
}

func expr(node *parse.Node) ([]core.Type, error) {
	return assign(node)
}

func assign(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.VarDecl:
		name := node.VarDeclField.Identifier.IdentField.Ident
		typ := node.VarDeclField.Type.DataTypeField.Type
		tc.SetVariable(name, NewVariable(name, typ))
		return nil, nil
	case parse.ShortVarDecl:
	case parse.Assign:
		// 値から計算する
		valNode := node.AssignField.Value
		valType, err := expr(valNode)
		if err != nil {
			return nil, err
		}
		target := node.AssignField.To
		switch target.Kind {
		case parse.VarDecl:
			name := target.VarDeclField.Identifier.IdentField.Ident
			typ := target.VarDeclField.Type.DataTypeField.Type
			if !isSameType(wrap(typ), valType) {
				return nil, fmt.Errorf("代入先と代入元の型が異なる: %v, %v", wrap(typ), valType)
			}
			tc.SetVariable(name, NewVariable(name, typ))
			return nil, nil
		case parse.Ident:
			name := target.IdentField.Ident
			v, ok := tc.FindVariableConsiderNest(name)
			if !ok {
				return nil, fmt.Errorf("未定義変数: %s", name)
			}
			if !isSameType(wrap(v.Type), valType) {
				return nil, fmt.Errorf("代入先と代入元の型が異なる: %v, %v", wrap(v.Type), valType)
			}
			return nil, nil
		default:
			return nil, fmt.Errorf("未対応の代入先: %s", target.Kind.String())
		}
	default:
		return andor(node)
	}
	return nil, fmt.Errorf("assign unimplemented")
}

func andor(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.And, parse.Or:
	default:
		return equality(node)
	}
	return nil, fmt.Errorf("andor unimplemented")
}

func equality(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.Eq, parse.Ne:
	default:
		return relational(node)
	}
	return nil, fmt.Errorf("equality unimplemnted")
}

func relational(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.Lt, parse.Le, parse.Gt, parse.Ge:
	default:
		return add(node)
	}
	return nil, fmt.Errorf("relational unimplemented")
}

func add(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.Add:
	case parse.Sub:
	default:
		return mul(node)
	}
	return nil, fmt.Errorf("add unimplemented")
}

func mul(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.Mul, parse.Div, parse.Mod:
	default:
		return unary(node)
	}
	return nil, fmt.Errorf("mul unimplemented")
}

func unary(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.Not:
	default:
		return primary(node)
	}
	return nil, fmt.Errorf("unary unimplemented")
}

func primary(node *parse.Node) ([]core.Type, error) {
	return access(node)
}

func access(node *parse.Node) ([]core.Type, error) {
	return literal(node)
}

func literal(node *parse.Node) ([]core.Type, error) {
	switch node.Kind {
	case parse.Parenthesis:
	case parse.Ident:
		name := node.IdentField.Ident
		v, ok := tc.FindVariableConsiderNest(name)
		if !ok {
			return nil, fmt.Errorf("未定義変数: %s", name)
		}
		return wrap(v.Type), nil
	case parse.Call:
		// 呼び出し側データ
		name := node.CallField.Identifier.IdentField.Ident

		var args []*parse.Node
		var argVars []core.Type
		if node.CallField.Args == nil {
			args = []*parse.Node{}
		} else {
			args = node.CallField.Args.PolynomialField.Values
			// 引数の型取得
			for _, v := range args {
				argVariableType, err := expr(v)
				if err != nil {
					return nil, err
				}
				argVars = append(argVars, argVariableType...)
			}
		}

		// 呼び出し先の型データ
		var paramVars []core.Type
		called, ok := tc.FindFunctionConsiderNest(name, false)
		if !ok {
			return nil, fmt.Errorf("未定義関数: %s", name)
		}
		for _, param := range called.Params {
			paramVars = append(paramVars, param.Type)
		}

		// 引数とパラメータの型が一致していることを
		if !isSameType(argVars, paramVars) {
			return nil, fmt.Errorf("与えられた引数と予想されていたパラメータが異なる: %v, %v", argVars, paramVars)
		}

		// 戻り値の型を取り出す
		var returns []core.Type
		for _, rt := range called.Returns {
			returns = append(returns, rt.Type)
		}
		return returns, nil
	default:
		return wrap(node.LiteralField.Literal.GetKind()), nil
	}
	return nil, fmt.Errorf("literal unimplemented")
}
