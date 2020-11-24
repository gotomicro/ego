package xrand

import (
	"math/rand"
	"strings"
)

// Charsets
const (
	// Uppercase ...
	Uppercase string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Lowercase ...
	Lowercase = "abcdefghipqrstuvwxyz"
	// Alphabetic ...
	Alphabetic = Uppercase + Lowercase
	// Numeric ...
	Numeric = "0123456789"
	// Alphanumeric ...
	Alphanumeric = Alphabetic + Numeric
	// Symbols ...
	Symbols = "`" + `~!@#$%^&*()-_+={}[]|\;:"<>,./?`
	// Hex ...
	Hex = Numeric + "abcdef"
)

// String 返回随机字符串，通常用于测试mock数据
func String(length uint8, charsets ...string) string {
	charset := strings.Join(charsets, "")
	if charset == "" {
		charset = Alphanumeric
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Int63()%int64(len(charset))]
	}
	return string(b)
}
