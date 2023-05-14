package typecheck

//
//import (
//	"fmt"
//	"github.com/x0y14/beat/core"
//	"github.com/x0y14/beat/parse"
//)
//
//var global *TypeMap
//var curtFn *Function
//
//func init() {
//	global = &TypeMap{
//		Functions:        nil,
//		Variables:        nil,
//		UnknownFunctions: nil,
//		UnknownVariables: nil,
//	}
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
//func fn(node *parse.Node) error {
//	field := node.FuncDefField
//	name := field.Identifier.IdentField.Ident
//	f := &Function{
//		Id:      GenerateId(),
//		Params:  []*Variable{},
//		Returns: []*Variable{},
//	}
//
//	// 引数
//	if field.Parameters != nil {
//		for _, paramNode := range field.Parameters.PolynomialField.Values {
//			param := paramNode.FuncParam
//			f.Params = append(f.Params, &Variable{
//				Id:   GenerateId(),
//				Name: param.Identifier.IdentField.Ident,
//				Type: param.DataType.DataTypeField.Type,
//			})
//		}
//	}
//
//	// 戻り値
//	if field.Returns != nil {
//		for _, returnTypeNode := range field.Returns.PolynomialField.Values {
//			f.Returns = append(f.Returns, &Variable{
//				Id:   GenerateId(),
//				Name: "",
//				Type: returnTypeNode.DataTypeField.Type,
//			})
//		}
//	}
//
//	// todo
//	curtFn = f
//
//	global.SetFn(name, f)
//	return nil
//}
//
//func for_(node *parse.Node) ([]core.Type, error) {
//	curtFn.currentNest++
//	ifBlock := &Block{
//		Vars:        []*Variable{},
//		Lower:       map[int][]*Block{},
//		currentNest: curtFn.currentNest,
//	}
//
//	curtFn.Lower[curtFn.currentNest] = append(curtFn.Lower[curtFn.currentNest], ifBlock)
//	curtFn.currentNest--
//	return nil, nil
//}
//
//func stmt(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Return:
//		var returnTypes []core.Type
//		for _, val := range node.PolynomialField.Values {
//			t, err := expr(val)
//			if err != nil {
//				return nil, err
//			}
//			returnTypes = append(returnTypes, t...)
//		}
//		return returnTypes, nil
//	case parse.IfElse:
//	case parse.For:
//	case parse.Block:
//		for _, s := range node.BlockField.Statements {
//			t, err := stmt(s)
//			if err != nil {
//				return nil, err
//			}
//			if s.Kind == parse.Return {
//				// 関数の戻り値の型と実際の戻り値の型を比較
//				var returnTypes []core.Type
//				for _, rt := range curtFn.Returns {
//					returnTypes = append(returnTypes, rt.Type)
//				}
//				if !isSameType(returnTypes, t) {
//					return nil, fmt.Errorf("定義されている戻り値と実際の戻り値の型が異なる: %v, %v", returnTypes, t)
//				}
//			}
//			return nil, nil
//		}
//	default:
//		return expr(node)
//	}
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
//		_, ok := curtFn.FindVarByName(name)
//		if ok {
//			return nil, fmt.Errorf("重複定義: %s", name)
//		}
//		// グローバル
//		_, ok = global.FindVarByName(name)
//		if ok {
//			return nil, fmt.Errorf("重複定義: %s", name)
//		}
//		// 未定義なので定義してあげる
//		_, err := curtFn.SetVar(name, t)
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
//		_, err = curtFn.SetVar(name, t[0])
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
//func andor(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.And, parse.Or:
//		lhs, err := andor(node.BinaryField.Lhs)
//		if err != nil {
//			return nil, err
//		}
//		rhs, err := andor(node.BinaryField.Rhs)
//		if err != nil {
//			return nil, err
//		}
//		if !isSameType(lhs, rhs) || !isSameType(wrap(core.Bool), lhs) {
//			return nil, fmt.Errorf("条件連結はBoolのみ: %v, %v", lhs, rhs)
//		}
//		return wrap(core.Bool), nil
//	default:
//		return equality(node)
//	}
//}
//
//func equality(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Eq, parse.Ne:
//		lhs, err := equality(node)
//		if err != nil {
//			return nil, err
//		}
//		rhs, err := equality(node)
//		if err != nil {
//			if !isSameType(lhs, rhs) {
//				return nil, fmt.Errorf("比較は同じ型のみ: %v, %v", lhs, rhs)
//			}
//		}
//		return wrap(core.Bool), nil
//	default:
//		return relational(node)
//	}
//}
//
//func relational(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Lt, parse.Le, parse.Gt, parse.Ge:
//		lhs, err := relational(node.BinaryField.Lhs)
//		if err != nil {
//			return nil, err
//		}
//		rhs, err := relational(node.BinaryField.Rhs)
//		if err != nil {
//			return nil, err
//		}
//		if !isSameType(lhs, rhs) {
//			return nil, fmt.Errorf("大小比較は同じ型のみで使用できます: %v, %v", lhs, rhs)
//		}
//		if !isComparable(lhs[0]) || !isComparable(rhs[0]) {
//			return nil, fmt.Errorf("大小比較は比較可能な型のみで使用できますInt,Float: %v, %v", lhs, rhs)
//		}
//		return wrap(core.Bool), nil
//	default:
//		return add(node)
//	}
//}
//
//func add(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Add:
//		lhs, err := add(node.BinaryField.Lhs)
//		if err != nil {
//			return nil, err
//		}
//		rhs, err := add(node.BinaryField.Rhs)
//		if err != nil {
//			return nil, err
//		}
//		if !isSameType(lhs, rhs) {
//			return nil, fmt.Errorf("計算はInt, Floatのみ: %v, %v", lhs, rhs)
//		}
//		// + だけは 文字列を許可
//		if (!isCalculable(lhs[0]) || !isCalculable(rhs[0])) && !isSameType(lhs, wrap(core.String)) {
//			return nil, fmt.Errorf("addはInt, Float, Stringのみ: %v, %v", lhs, rhs)
//		}
//		return lhs, nil
//	case parse.Sub:
//		lhs, err := add(node.BinaryField.Lhs)
//		if err != nil {
//			return nil, err
//		}
//		rhs, err := add(node.BinaryField.Rhs)
//		if err != nil {
//			return nil, err
//		}
//		if !isSameType(lhs, rhs) {
//			return nil, fmt.Errorf("計算はInt, Floatのみ: %v, %v", lhs, rhs)
//		}
//		if !isCalculable(lhs[0]) || !isCalculable(rhs[0]) {
//			return nil, fmt.Errorf("subはInt, Floatのみ: %v, %v", lhs, rhs)
//		}
//		return lhs, nil
//	default:
//		return mul(node)
//	}
//}
//
//func mul(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Mul, parse.Div, parse.Mod:
//		lhs, err := mul(node.BinaryField.Lhs)
//		if err != nil {
//			return nil, err
//		}
//		rhs, err := mul(node.BinaryField.Rhs)
//		if err != nil {
//			return nil, err
//		}
//		if !isSameType(lhs, rhs) {
//			return nil, fmt.Errorf("計算はInt, Floatのみ: %v, %v", lhs, rhs)
//		}
//		if !isCalculable(lhs[0]) || !isCalculable(rhs[0]) {
//			return nil, fmt.Errorf("mul,div, modはInt, Floatのみ: %v, %v", lhs, rhs)
//		}
//		return lhs, nil
//	default:
//		return unary(node)
//	}
//}
//
//func unary(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Not:
//		p, err := primary(node.UnaryField.Value)
//		if err != nil {
//			return nil, err
//		}
//		if !isSameType(p, []core.Type{core.Bool}) {
//			return nil, fmt.Errorf("boolではないものにnotをつけることはできません: %v", p)
//		}
//		return p, nil
//	default:
//		return primary(node)
//	}
//}
//
//func primary(node *parse.Node) ([]core.Type, error) {
//	return access(node)
//}
//
//func access(node *parse.Node) ([]core.Type, error) {
//	if node.Kind == parse.Prefix {
//		return literal(node.PrefixField.Child)
//	}
//	return literal(node)
//}
//
//func literal(node *parse.Node) ([]core.Type, error) {
//	switch node.Kind {
//	case parse.Parenthesis:
//		// todo
//		return nil, nil
//	case parse.Ident:
//		ident := node.IdentField.Ident
//		// キーワードなら即座に返す
//		if ident == "true" || ident == "false" {
//			return wrap(core.Bool), nil
//		}
//		// 関数内で検索, curtFn
//		v, ok := curtFn.FindVarByName(ident)
//		if ok {
//			return wrap(v.Type), nil
//		}
//		// グローバル内で検索
//		v, ok = global.FindVarByName(ident)
//		if ok {
//			return wrap(v.Type), nil
//		}
//		// unknownへ
//		return nil, fmt.Errorf("未定義変数: %s", ident)
//	case parse.Call:
//		// 型が複数の可能性あり, ex) f() (string, bool)
//		// todo
//		// なかったらunknownへ
//		return nil, nil
//	default:
//		return wrap(node.LiteralField.Literal.GetKind()), nil
//	}
//}
//
//func TypeCheck(nodes []*parse.Node) *TypeMap {
//	return global
//}
