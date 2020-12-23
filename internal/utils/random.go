package utils

import (
	crypto "crypto/rand"
	"fmt"
	"io"
	"math/rand"
)

// RandomString : create a random alpha-numeric string
func RandomString(chars int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, chars)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// ShortID : create a short unique identifier of length 12
func ShortID() string {
	b := make([]byte, 12)
	_, err := io.ReadFull(crypto.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
