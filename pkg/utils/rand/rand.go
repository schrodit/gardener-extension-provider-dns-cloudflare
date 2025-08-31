package rand

import (
	"math/rand/v2"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

// RandStringRunes generates a random string of length n
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(b)
}
