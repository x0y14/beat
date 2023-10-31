package tokenize

import (
	"fmt"
	"github.com/x0y14/beat/core"
	"strconv"
	"unicode"
)

var userInput []rune
var curtPos *core.Position

var singleOpSymbols []string
var compositeOpSymbols []string

func init() {
	singleOpSymbols = []string{
		"(", ")", "[", "]", "{", "}",
		".", ",", ":", ";",
		"+", "-", "*", "/", "%",
		">", "<",
		"=", "!",
	}
	compositeOpSymbols = []string{
		"==", "!=", ">=", "<=",
		"+=", "-=", "*=", "/=", "%=",
		"&&", "||", ":=",
	}
}

func startWith(s string) bool {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if len(userInput) <= curtPos.Wat+i || userInput[curtPos.Wat+i] != runes[i] {
			return false
		}
	}
	return true
}

func isIdentRune(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		('0' <= r && r <= '9') ||
		('_' == r || '.' == r)
}

func isEof() bool {
	return curtPos.Wat >= len(userInput)
}

func consumeComment() string {
	// "//"
	curtPos.Lat += 2
	curtPos.Wat += 2
	var s string
	for !isEof() {
		if userInput[curtPos.Wat] == '\n' {
			break
		}
		s += string(userInput[curtPos.Wat])
		curtPos.Lat++
		curtPos.Wat++
	}
	return s
}

func consumeIdent() string {
	var s string
	for !isEof() {
		if !isIdentRune(userInput[curtPos.Wat]) {
			break
		}
		s += string(userInput[curtPos.Wat])
		curtPos.Lat++
		curtPos.Wat++
	}
	// textは後で検査
	return s
}

func consumeString() string {
	var s string
	// "
	curtPos.Lat++
	curtPos.Wat++

	for !isEof() {
		if userInput[curtPos.Wat] == '"' {
			break
		}
		// escaped double quotation
		if userInput[curtPos.Wat] == '\\' && userInput[curtPos.Wat+1] == '"' {
			s += "\""
			curtPos.Lat += 2
			curtPos.Wat += 2
			continue
		}
		// newline
		if userInput[curtPos.Wat] == '\\' && userInput[curtPos.Wat+1] == 'n' {
			s += "\n"
			curtPos.Lat += 2
			curtPos.Wat += 2
			continue
		}
		// tab
		if userInput[curtPos.Wat] == '\\' && userInput[curtPos.Wat+1] == 't' {
			s += "\t"
			curtPos.Lat += 2
			curtPos.Wat += 2
			continue
		}
		// escaped single quotation
		if userInput[curtPos.Wat] == '\\' && userInput[curtPos.Wat+1] == '\'' {
			s += "'"
			curtPos.Lat += 2
			curtPos.Wat += 2
			continue
		}
		// escaped? slash
		if userInput[curtPos.Wat] == '\\' && userInput[curtPos.Wat+1] == '\\' {
			s += "\\"
			curtPos.Lat += 2
			curtPos.Wat += 2
			continue
		}

		s += string(userInput[curtPos.Wat])
		curtPos.Lat++
		curtPos.Wat++
	}

	// "
	curtPos.Lat++
	curtPos.Wat++
	return s
}

func consumeNumber() (string, bool) {
	isFloat := false
	var s string
	for !isEof() {
		if unicode.IsDigit(userInput[curtPos.Wat]) {
			s += string(userInput[curtPos.Wat])
			curtPos.Lat++
			curtPos.Wat++
			continue
		} else if userInput[curtPos.Wat] == '.' {
			// ポイントの次が、数字じゃなければ強制終了
			if len(userInput) <= curtPos.Wat+1 ||
				!unicode.IsDigit(userInput[curtPos.Wat+1]) {
				break
			}
			s += string(userInput[curtPos.Wat])
			curtPos.Lat++
			curtPos.Wat++
			isFloat = true
			continue
		} else {
			break
		}
	}
	return s, isFloat
}

func consumeWhite() string {
	var s string
	for !isEof() {
		if userInput[curtPos.Wat] == ' ' || userInput[curtPos.Wat] == '\t' {
			s += string(userInput[curtPos.Wat])
			curtPos.Lat++
			curtPos.Wat++
		} else {
			break
		}
	}
	return s
}

func Tokenize(input string) (*Token, error) {
	userInput = []rune(input)
	curtPos = core.NewPosition(1, 0, 0)
	var head Token
	cur := &head
Loop:
	for !isEof() {
		// white
		if userInput[curtPos.Wat] == ' ' || userInput[curtPos.Wat] == '\t' {
			// データを捨てるので入力の完全な復元はできなくなる
			_ = curtPos.Clone()
			_ = consumeWhite()
			continue
		}
		// nl
		if userInput[curtPos.Wat] == '\n' {
			curtPos.LineNo++
			curtPos.Lat = 0
			curtPos.Wat++
			continue
		}
		// comment
		if userInput[curtPos.Wat] == '/' && userInput[curtPos.Wat+1] == '/' {
			_ = consumeComment()
			continue
		}
		// op, symbol
		for _, r := range append(compositeOpSymbols, singleOpSymbols...) {
			if startWith(r) {
				opSm, err := NewOpSymbolToken(curtPos.Clone(), r)
				if err != nil {
					return nil, err
				}
				cur = Chain(cur, opSm)
				curtPos.Lat += len(r)
				curtPos.Wat += len(r)
				continue Loop // op, symbolではなく全体のループへ
			}
		}
		// ident
		if isIdentRune(userInput[curtPos.Wat]) && !unicode.IsDigit(userInput[curtPos.Wat]) {
			pos := curtPos.Clone()
			id := consumeIdent()
			cur = Chain(cur, NewIdentToken(pos, id))
			continue
		}
		// string
		if userInput[curtPos.Wat] == '"' {
			pos := curtPos.Clone()
			s := consumeString()
			lit, err := NewLiteralToken(pos, core.NewLiteral(s))
			if err != nil {
				return nil, err
			}
			cur = Chain(cur, lit)
			continue
		}
		// number
		if unicode.IsDigit(userInput[curtPos.Wat]) {
			pos := curtPos.Clone()
			numStr, isFloat := consumeNumber()
			if isFloat {
				n, err := strconv.ParseFloat(numStr, 64)
				if err != nil {
					return nil, err
				}
				lit, err := NewLiteralToken(pos, core.NewLiteral(n))
				if err != nil {
					return nil, err
				}
				cur = Chain(cur, lit)
				continue
			} else {
				n, err := strconv.ParseInt(numStr, 10, 0)
				if err != nil {
					return nil, err
				}
				lit, err := NewLiteralToken(pos, core.NewLiteral(int(n)))
				if err != nil {
					return nil, err
				}
				cur = Chain(cur, lit)
				continue
			}
		}
		return nil, fmt.Errorf("unexpected charactor: %s", string(userInput[curtPos.Wat]))
	}
	cur = Chain(cur, NewEofToken(curtPos.Clone()))
	return head.Next, nil
}
