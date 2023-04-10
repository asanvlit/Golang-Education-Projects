package generator

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, size)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:size]
}

func RandomStringArray(size int, minLength int, maxLength int) []string {
	strings := make([]string, size)

	for i := 0; i < size; i++ {
		strings[i] = RandomString(minLength + rand.Intn(maxLength-minLength+1))
	}

	return strings
}
