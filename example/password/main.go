package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/bitfield/qrand"
)

const chars = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+-={}|[]\\:\"<>?,./`

// Shows how to create a [rand.Rand] using qrand as the source.
func main() {
	apiKey := os.Getenv("AQN_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Please set AQN_API_KEY to a valid API key: see https://quantumnumbers.anu.edu.au/documentation for details.")
		os.Exit(1)
	}
	rnd := rand.New(qrand.NewSource(qrand.NewReader(apiKey)))
	password := make([]byte, 32)
	for i := range password {
		password[i] = chars[rnd.Intn(len(chars))]
	}
	fmt.Println(string(password))
}
