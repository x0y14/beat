package tokenize

import (
	"fmt"
	"github.com/x0y14/beat/compile/core"
)

type Token struct {
	Kind TokenKind
	Pos  *core.Position
	Lit  *core.Literal
	Next *Token
}

func NewToken(kind TokenKind, pos *core.Position, lit *core.Literal) *Token {
	return &Token{
		Kind: kind,
		Pos:  pos,
		Lit:  lit,
		Next: nil,
	}
}

func NewIdentToken(pos *core.Position, ident string) *Token {
	return NewToken(Ident, pos, core.NewLiteral(ident))
}

//func NewCommentToken(Pos *core.Position, comment string) *Token {
//	return NewToken(Comment, Pos, core.NewLiteral(comment))
//}
//
//func NewWhiteToken(Pos *core.Position, white string) *Token {
//	return NewToken(White, Pos, core.NewLiteral(white))
//}
//
//func NewLineToken(Pos *core.Position, nl string) *Token {
//	return NewToken(Newline, Pos, core.NewLiteral(nl))
//}

func NewLiteralToken(pos *core.Position, lit *core.Literal) (*Token, error) {
	switch lit.GetKind() {
	case core.Int:
		return NewToken(Int, pos, lit), nil
	case core.Float:
		return NewToken(Float, pos, lit), nil
	case core.String:
		return NewToken(String, pos, lit), nil
	default:
		return nil, fmt.Errorf("unsuppoted Lit Kind: %s", lit.GetKind().String())
	}
}

func NewOpSymbolToken(pos *core.Position, opSymbol string) (*Token, error) {
	switch opSymbol {
	case "(":
		return NewToken(Lrb, pos, nil), nil
	case ")":
		return NewToken(Rrb, pos, nil), nil
	case "[":
		return NewToken(Lsb, pos, nil), nil
	case "]":
		return NewToken(Rsb, pos, nil), nil
	case "{":
		return NewToken(Lcb, pos, nil), nil
	case "}":
		return NewToken(Rcb, pos, nil), nil
	case ".":
		return NewToken(Dot, pos, nil), nil
	case ",":
		return NewToken(Comma, pos, nil), nil
	case ":":
		return NewToken(Colon, pos, nil), nil
	case ";":
		return NewToken(Semi, pos, nil), nil

	case "+":
		return NewToken(Add, pos, nil), nil
	case "-":
		return NewToken(Sub, pos, nil), nil
	case "*":
		return NewToken(Mul, pos, nil), nil
	case "/":
		return NewToken(Div, pos, nil), nil
	case "%":
		return NewToken(Mod, pos, nil), nil

	case "==":
		return NewToken(Eq, pos, nil), nil
	case "!=":
		return NewToken(Ne, pos, nil), nil
	case ">":
		return NewToken(Gt, pos, nil), nil
	case "<":
		return NewToken(Lt, pos, nil), nil
	case ">=":
		return NewToken(Ge, pos, nil), nil
	case "<=":
		return NewToken(Le, pos, nil), nil

	case "=":
		return NewToken(Assign, pos, nil), nil
	case "+=":
		return NewToken(AddAssign, pos, nil), nil
	case "-=":
		return NewToken(SubAssign, pos, nil), nil
	case "*=":
		return NewToken(MulAssign, pos, nil), nil
	case "/=":
		return NewToken(DivAssign, pos, nil), nil
	case "%=":
		return NewToken(ModAssign, pos, nil), nil
	case ":=":
		return NewToken(ColonAssign, pos, nil), nil

	case "&&":
		return NewToken(And, pos, nil), nil
	case "||":
		return NewToken(Or, pos, nil), nil
	case "!":
		return NewToken(Not, pos, nil), nil

	default:
		return nil, fmt.Errorf("unsupported op or symbol: %s", opSymbol)
	}
}

func NewEofToken(pos *core.Position) *Token {
	return NewToken(Eof, pos, nil)
}

func Chain(cur *Token, next *Token) *Token {
	cur.Next = next
	return next
}

//func IsReservedKeyword(ident string) bool {
//	switch ident {
//	case "true", "false", "nil":
//		return true
//	default:
//		return false
//	}
//}
