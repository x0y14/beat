package typecheck

import "github.com/x0y14/beat/core"

func variablesToTypes(variables []*Variable) []core.Type {
	var types []core.Type
	for _, v := range variables {
		types = append(types, v.Type)
	}
	return types
}

func wrap(t ...core.Type) []core.Type {
	return t
}

func isSameType(x1, x2 []core.Type) bool {
	// そもそも個数が一致しているか
	if len(x1) != len(x2) {
		return false
	}
	// 順番に型が一致しているか
	for i, x := range x1 {
		if x != x2[i] {
			return false
		}
	}
	return true
}

func isCalculable(x core.Type) bool {
	switch x {
	case core.Int, core.Float:
		return true
	default:
		return false
	}
}

func isComparable(x core.Type) bool {
	switch x {
	case core.Int, core.Float:
		return true
	default:
		return false
	}
}
