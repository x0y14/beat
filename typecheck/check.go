package typecheck

//
//import (
//	"fmt"
//	"github.com/x0y14/beat/core"
//	"github.com/x0y14/beat/parse"
//)
//
//var global TypeMap
//var currentFn *Function
//var currentBlock *Block
//var nest int
//
//func init() {
//	global = TypeMap{
//		Functions: map[string]*Function{},
//		Variables: map[string]*Variable{},
//		//UnknownFunctions: map[Identifier]*Function{},
//		//UnknownVariables: map[Identifier]*Variable{},
//	}
//	nest = 0
//}
//
//func wrap(t ...core.Type) []core.Type {
//	return t
//}
//
//func isSameType(x1, x2 []core.Type) bool {
//	// そもそも個数が一致しているか
//	if len(x1) != len(x2) {
//		return false
//	}
//	// 順番に型が一致しているか
//	for i, x := range x1 {
//		if x != x2[i] {
//			return false
//		}
//	}
//	return true
//}
//
//func isCalculable(x core.Type) bool {
//	switch x {
//	case core.Int, core.Float:
//		return true
//	default:
//		return false
//	}
//}
//
//func isComparable(x core.Type) bool {
//	switch x {
//	case core.Int, core.Float:
//		return true
//	default:
//		return false
//	}
//}
//
//func toplevel(node *parse.Node) ([]core.Type, error) {
//	return nil, nil
//}
//
//func function(node *parse.Node) ([]core.Type, error) {
//	fnBackUp := *currentFn
//
//	fn := &Function{
//		Id:      GenerateId(),
//		Params:  []*Variable{},
//		Returns: []*Variable{},
//		Block:   nil,
//	}
//	currentFn = fn
//
//	// 引数
//	if node.FuncDefField.Parameters != nil {
//		for _, param := range node.FuncDefField.Parameters.PolynomialField.Values {
//			paramData := param.FuncParam
//			fn.Params = append(fn.Params, &Variable{
//				Id:   GenerateId(),
//				Name: paramData.Identifier.IdentField.Ident,
//				Type: paramData.DataType.DataTypeField.Type,
//			})
//		}
//	}
//	// 戻り値
//	if node.FuncDefField.Returns != nil {
//		for _, ret := range node.FuncDefField.Returns.PolynomialField.Values {
//			fn.Returns = append(fn.Returns, &Variable{
//				Id:   GenerateId(),
//				Name: "",
//				Type: ret.DataTypeField.Type,
//			})
//		}
//	}
//
//	// 中身
//	_, err := block_(node.FuncDefField.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	currentFn = &fnBackUp
//	return nil, nil
//}
//
//func stmt(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Block:
//	}
//	return nil, nil
//}
//
//func for_(node *parse.Node) ([]core.Type, error) {
//	return nil, nil
//}
//
//func if_(node *parse.Node) ([]core.Type, error) {
//	return nil, nil
//}
//func block_(node *parse.Node) ([]core.Type, error) {
//	nest++
//	b := &Block{
//		Nest:  nest,
//		Vars:  []*Variable{},
//		Lower: map[int][]*Block{},
//	}
//	if nest == 1 {
//		// もしnest==1なら、関数のボディ
//		currentFn.Block = b
//	} else {
//		// そうでないなら、別ブロックの中のブロックなので追加してあげる
//		if currentBlock.Lower[nest] == nil {
//			currentBlock.Lower[nest] = []*Block{}
//		}
//		currentBlock.Lower[nest] = append(currentBlock.Lower[nest], b)
//	}
//	currentBlock = b
//	for _, statement := range node.BlockField.Statements {
//		_, err := stmt(statement)
//		if err != nil {
//			return nil, err
//		}
//	}
//	nest--
//	return nil, nil
//}
//
//func expr(node *parse.Node) ([]core.Type, error) {
//	return assign(node)
//}
//
//func assign(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.VarDecl: // var x type
//		name := node.VarDeclField.Identifier.IdentField.Ident
//		t := node.VarDeclField.Type.DataTypeField.Type
//		// 現在の関数内
//		_, ok := currentBlock.FindVarByName(name)
//		if ok {
//			return nil, fmt.Errorf("重複定義: %s", name)
//		}
//		// グローバル
//		_, ok = global.FindVarByName(name)
//		if ok {
//			return nil, fmt.Errorf("重複定義: %s", name)
//		}
//		// 未定義なので定義してあげる
//		_, err := currentFn.SetVar(name, t)
//		return nil, err
//	case parse.ShortVarDecl:
//		name := node.ShortVarDeclField.Identifier.IdentField.Ident
//		t, err := expr(node.ShortVarDeclField.Value)
//		if err != nil {
//			return nil, err
//		}
//		if len(t) != 1 {
//			// todo: パースで処理すべき
//			return nil, fmt.Errorf("複数の変数を同時に定義することはできない")
//		}
//		_, err = currentFn.SetVar(name, t[0])
//		if err != nil {
//			return nil, err
//		}
//		return nil, nil
//	case parse.Assign:
//		expectType, err := assign(node.AssignField.To)
//		if err != nil {
//			return nil, err
//		}
//		actualType, err := assign(node.AssignField.Value)
//		if err != nil {
//			return nil, err
//		}
//		if !isSameType(expectType, actualType) {
//			return nil, fmt.Errorf("一致しない代入: %v, %v", expectType, actualType)
//		}
//		return nil, nil
//	default:
//		return andor(node)
//	}
//}
//
//func check(node []*parse.Node) (*TypeMap, error) {
//	for _, top := range node {
//		_, err := toplevel(top)
//		if err != nil {
//			return nil, err
//		}
//	}
//	return &global, nil
//}
