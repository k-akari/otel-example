package main

import (
	"context"
	"log"
	"os"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run(_ context.Context) error {
	return nil
}
