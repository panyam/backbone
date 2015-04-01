package utils

import (
	"math/rand"
)

func RandSeq(chars string, n int) string {
	letters := []rune(chars)
	numLetters := len(letters)
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(numLetters)]
	}
	return string(b)
}

func RandDigits(n int) string {
	return RandSeq("0123456789", n)
}

func RandAlnum(n int) string {
	return RandSeq("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", n)
}

func RandAlpha(n int) string {
	return RandSeq("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", n)
}
