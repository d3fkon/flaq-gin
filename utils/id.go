package utils

import (
	"math/rand"
	"time"
)

func generateId(length uint8) string {
	rand.Seed(time.Now().UnixNano())
	charset := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	// var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, length)
	for i := range s {
		s[i] = charset[rand.Intn(len(charset))]
	}
	return string(s)
}

func GenerateReferral() string {
	return generateId(3) + "-" + generateId(4)
}
