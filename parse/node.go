package parse

import "github.com/x0y14/beat/core"

type Node struct {
	Kind NodeKind
	Pos  *core.Position

	DataTypeField     *DataTypeField
	ImportField       *ImportField
	IdentField        *IdentField
	VarDeclField      *VarDeclField
	AssignField       *AssignField
	CommentField      *CommentField
	FuncDefField      *FuncDefField
	BlockField        *BlockField
	ReturnField       *ReturnField
	IfElseField       *IfElseField
	ForField          *ForField
	ShortVarDeclField *ShortVarDeclField
	BinaryField       *BinaryField
	UnaryField        *UnaryField
	CallField         *CallField
	PolynomialField   *PolynomialField
	FuncParam         *FuncParam
	PrefixField       *PrefixField
	LiteralField      *LiteralField
}

func NewNode(kind NodeKind, pos *core.Position) *Node {
	return &Node{Kind: kind, Pos: pos}
}

func NewDataTypeNode(pos *core.Position, typ core.Types) *Node {
	n := NewNode(DataType, pos)
	n.DataTypeField = &DataTypeField{Type: typ}
	return n
}

func NewImportNode(pos *core.Position, packageNameOrUrl string) *Node {
	n := NewNode(Import, pos)
	n.ImportField = &ImportField{Target: packageNameOrUrl}
	return n
}

func NewIdentNode(pos *core.Position, ident string) *Node {
	n := NewNode(Ident, pos)
	n.IdentField = &IdentField{Ident: ident}
	return n
}

func NewVarDeclNode(pos *core.Position, ident, typ *Node) *Node {
	n := NewNode(VarDecl, pos)
	n.VarDeclField = &VarDeclField{
		Identifier: ident,
		Type:       typ,
	}
	return n
}

func NewAssignNode(pos *core.Position, to, value *Node) *Node {
	n := NewNode(Assign, pos)
	n.AssignField = &AssignField{
		To:    to,
		Value: value,
	}
	return n
}

func NewCommentNode(pos *core.Position, comment string) *Node {
	n := NewNode(Comment, pos)
	n.CommentField = &CommentField{Comment: comment}
	return n
}

func NewFuncDefNode(pos *core.Position, ident, params, returns, body *Node) *Node {
	n := NewNode(FuncDef, pos)
	n.FuncDefField = &FuncDefField{
		Identifier: ident,
		Parameters: params,
		Returns:    returns,
		Body:       body,
	}
	return n
}

func NewBlockNode(pos *core.Position, stmts []*Node) *Node {
	n := NewNode(Block, pos)
	n.BlockField = &BlockField{Statements: stmts}
	return n
}

func NewIfElseNode(pos *core.Position, useElse bool, cond, ifBlock, elseBlock *Node) *Node {
	n := NewNode(IfElse, pos)
	n.IfElseField = &IfElseField{
		UseElse:   useElse,
		Cond:      cond,
		IfBlock:   ifBlock,
		ElseBlock: elseBlock,
	}
	return n
}

func NewForNode(pos *core.Position, init, cond, loop, body *Node) *Node {
	n := NewNode(For, pos)
	n.ForField = &ForField{
		Init: init,
		Cond: cond,
		Loop: loop,
		Body: body,
	}
	return n
}

func NewShortVarDeclNode(pos *core.Position, ident, value *Node) *Node {
	n := NewNode(ShortVarDecl, pos)
	n.ShortVarDeclField = &ShortVarDeclField{
		Identifier: ident,
		Value:      value,
	}
	return n
}

func NewBinaryNode(kind NodeKind, pos *core.Position, lhs, rhs *Node) *Node {
	n := NewNode(kind, pos)
	n.BinaryField = &BinaryField{
		Lhs: lhs,
		Rhs: rhs,
	}
	return n
}

func NewUnaryNode(kind NodeKind, pos *core.Position, value *Node) *Node {
	n := NewNode(kind, pos)
	n.UnaryField = &UnaryField{
		Value: value,
	}
	return n
}

func NewLiteralNode(pos *core.Position, literal *core.Literal) *Node {
	n := NewNode(Literal, pos)
	n.LiteralField = &LiteralField{Literal: literal}
	return n
}

func NewCallNode(pos *core.Position, ident, args *Node) *Node {
	n := NewNode(Call, pos)
	n.CallField = &CallField{
		Identifier: ident,
		Args:       args,
	}
	return n
}

func NewPolynomialNode(kind NodeKind, pos *core.Position, values []*Node) *Node {
	n := NewNode(kind, pos)
	n.PolynomialField = &PolynomialField{Values: values}
	return n
}

func NewFuncParamNode(pos *core.Position, ident, typ *Node) *Node {
	n := NewNode(Param, pos)
	n.FuncParam = &FuncParam{
		Identifier: ident,
		DataType:   typ,
	}
	return n
}

func NewPrefixNode(pos *core.Position, prefix string, child *Node) *Node {
	n := NewNode(Prefix, pos)
	n.PrefixField = &PrefixField{
		Prefix: prefix,
		Child:  child,
	}
	return n
}
