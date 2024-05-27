package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/k-akari/otel-example/internal/infra/otel"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	env := mustNewConfig()

	_, close, err := otel.NewTracer(ctx, "client", env.GCPProjectID)
	if err != nil {
		return fmt.Errorf("failed to create tracer: %w", err)
	}
	defer close()

	_, err = net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return nil
}
