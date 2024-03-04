package util

import (
	"crypto/rand"
	"io"
	"strings"
)

var (
	table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

func FormatPhone(phone string) string {
	return strings.ReplaceAll(phone, "+", "")
}

func GenerateCode(max int) (string, error) {
	b := make([]byte, max)

	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "", err
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b), nil
}
