package utils

import (
	"fmt"
	"math/rand"
	"strconv"
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

func ID2String(id int64) string {
	return fmt.Sprintf("%d", id)
}

func String2ID(strid string) int64 {
	val, err := strconv.ParseInt(strid, 10, 64)
	if err != nil {
		return 0
	}
	return val
}
