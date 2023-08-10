package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/bitfield/qrand"
)

const Usage = `Usage: qrand NUM_BYTES

Queries the Australian National University Quantum Numbers (AQN) API to get the specified number of random bytes (maximum 1024), and prints the result as a hexadecimal string.

The program expects to read a valid API key from the environment variable AQN_API_KEY.

See https://quantumnumbers.anu.edu.au/documentation for details about the API.`

func main() {
	apiKey := os.Getenv("AQN_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Please set AQN_API_KEY to a valid API key: see https://quantumnumbers.anu.edu.au/documentation for details.")
		os.Exit(1)
	}
	q := qrand.NewReader(apiKey)
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, Usage)
		os.Exit(0)
	}
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, Usage)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	buf := make([]byte, n)
	_, err = q.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(hex.EncodeToString(buf))
}
