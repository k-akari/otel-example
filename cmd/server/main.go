package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run(_ context.Context) error {
	env := mustNewConfig()

	_, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return nil
}
