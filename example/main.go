package main

import (
	"fmt"
	"log"

	"github.com/bitfield/qrand"
)

func main() {
	numbers := make([]byte, 10)
	_, err := qrand.Read(numbers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(numbers)
}
