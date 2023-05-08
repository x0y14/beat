package tokenize

import (
	"github.com/stretchr/testify/assert"
	"github.com/x0y14/beat/core"
	"testing"
	"unicode/utf8"
)

func GenPosForTest(str string) *core.Position {
	_wat := utf8.RuneCountInString(str)
	_lat := 0
	_ln := 1
	for _, r := range []rune(str) {
		if r == '\n' {
			_ln++
			_lat = 0
		} else {
			_lat++
		}
	}
	return core.NewPosition(_ln, _lat, _wat)
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name          string
		in            string
		expectedToken *Token
	}{
		{
			"add",
			"1 + 1",
			&Token{
				kind:    Int,
				pos:     GenPosForTest(""),
				literal: core.NewLiteral(1),
				next: &Token{
					kind:    Add,
					pos:     GenPosForTest("1 "),
					literal: nil,
					next: &Token{
						kind:    Int,
						pos:     GenPosForTest("1 + "),
						literal: core.NewLiteral(1),
						next: &Token{
							kind:    Eof,
							pos:     GenPosForTest("1 + 1"),
							literal: nil,
							next:    nil,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok, err := Tokenize(tt.in)
			assert.Equal(t, tt.expectedToken, tok)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
