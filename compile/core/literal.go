package core

import "fmt"

type Literal struct {
	typ Type
	i   int
	f   float64
	s   string
}

func NewLiteral[T string | int | float64](v T) *Literal {
	switch any(v).(type) {
	case string:
		return &Literal{
			typ: String,
			s:   any(v).(string),
		}
	case int:
		return &Literal{
			typ: Int,
			i:   any(v).(int),
		}
	case float64:
		return &Literal{
			typ: Float,
			f:   any(v).(float64),
		}
	}
	return nil
}

func (l *Literal) String() string {
	var v any
	switch l.typ {
	case String:
		v = l.s
	case Int:
		v = l.i
	case Float:
		v = l.f
	}
	return fmt.Sprintf("Literal{ kind: %s, value: %v }", l.typ.String(), v)
}

func (l *Literal) GetKind() Type {
	return l.typ
}

func (l *Literal) GetString() string {
	return l.s
}
func (l *Literal) SetString(s string) {
	l.s = s
}

func (l *Literal) GetInt() int {
	return l.i
}
func (l *Literal) SetInt(i int) {
	l.i = i
}

func (l *Literal) GetFloat() float64 {
	return l.f
}
func (l *Literal) SetFloat(f float64) {
	l.f = f
}
