package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitfield/qrand"
)

// Shows how to read random numbers directly from a qrand reader.
func main() {
	apiKey := os.Getenv("AQN_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Please set AQN_API_KEY to a valid API key: see https://quantumnumbers.anu.edu.au/documentation for details.")
		os.Exit(1)
	}
	q := qrand.NewReader(apiKey)
	numbers := make([]byte, 10)
	_, err := q.Read(numbers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(numbers)
}
