package typecheck

import "math/rand"

type Identifier string

const (
	idLength = 10
)

// 参考: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateId() Identifier {
	b := make([]rune, idLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return Identifier(b)
}
