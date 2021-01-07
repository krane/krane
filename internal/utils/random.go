package utils

import (
	crypto "crypto/rand"
	"fmt"
	"io"
	"math/rand"
)

// RandomString creates a random alpha-numeric string of size n
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// ShortID returns a short unique identifier of length 12
func ShortID() string {
	b := make([]byte, 12)
	_, err := io.ReadFull(crypto.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
