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

	tracer, close, err := otel.NewTracer(ctx, "grpc_server", env.EndpointJaeger)
	if err != nil {
		return fmt.Errorf("failed to create tracer: %w", err)
	}
	defer close()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := newServer(l, tracer)
	s.registerServices()
	return s.run(ctx)
}
