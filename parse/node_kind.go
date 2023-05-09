package parse

type NodeKind int

const (
	Illegal NodeKind = iota

	Block
	Return
	ReturnTypes
	If
	IfElse
	While
	For

	Import

	Not // !
	Plus
	Minus

	And // &&
	Or  // ||
	Eq  // ==
	Ne  // !=
	Lt  // <
	Le  // <=
	Gt  // >
	Ge  // >=
	Add // +
	Sub // -
	Mul // *
	Div // /
	Mod // %

	FuncDef
	VarDecl
	ShortVarDecl
	Assign

	DataType

	Literal

	Ident
	Call

	List
	Dict
	KV
	Bool
	True
	False
	Void
	Nil
	Any

	Args
	Params
	Param

	Access
	Parenthesis

	Prefix

	Comment
)

var kind = []string{
	Illegal:      "Illegal",
	Block:        "Block",
	Return:       "Return",
	ReturnTypes:  "ReturnTypes",
	If:           "If",
	IfElse:       "IfElse",
	While:        "While",
	For:          "For",
	Import:       "Import",
	Not:          "Not",
	Plus:         "Plus",
	Minus:        "Minus",
	And:          "And",
	Or:           "Or",
	Eq:           "Eq",
	Ne:           "Ne",
	Lt:           "Lt",
	Le:           "Le",
	Gt:           "Gt",
	Ge:           "Ge",
	Add:          "Add",
	Sub:          "Sub",
	Mul:          "Mul",
	Div:          "Div",
	Mod:          "Mod",
	FuncDef:      "FuncDef",
	VarDecl:      "VarDecl",
	ShortVarDecl: "ShortVarDecl",
	Assign:       "Assign",
	DataType:     "DataType",
	Literal:      "Literal",
	Ident:        "Ident",
	Call:         "Call",
	List:         "List",
	Dict:         "Dict",
	KV:           "KV",
	Bool:         "Bool",
	True:         "True",
	False:        "False",
	Void:         "Void",
	Nil:          "Nil",
	Any:          "Any",
	Args:         "Args",
	Params:       "Params",
	Param:        "Param",
	Access:       "Access",
	Parenthesis:  "Parenthesis",
	Prefix:       "Prefix",
	Comment:      "Comment",
}

func (nk NodeKind) String() string {
	return kind[nk]
}
