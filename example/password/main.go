package main

import (
	"fmt"
	"math/rand"

	"github.com/bitfield/qrand"
)

func main() {
	var random = rand.New(qrand.NewSource())
	// nchars stores the password length
	var nchars = 15
	var password []rune
	for i := 0; i < nchars; i++ {
		offset := random.Intn(26)
		password = append(password, rune('a' + offset))
	}
	fmt.Println(string(password))
}
