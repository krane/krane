package utils

import "math/rand"

func RandomString(chars int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, chars)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}