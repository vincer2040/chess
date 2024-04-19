package main

import (
	"log"

	"github.com/vincer2040/chess/cmd/chess"
)

func main() {
	err := chess.Main()
	if err != nil {
		log.Fatalf("error: %v+\n", err)
	}
}
