package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/k-akari/otel-example/internal/infra/database"
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

	close, err := otel.Init(ctx, "grpc_server", env.EndpointJaeger)
	if err != nil {
		return fmt.Errorf("failed to create tracer: %w", err)
	}
	defer close()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	db, close, err := database.NewClient(env.DBUser, env.DBPass, env.DBHost, env.DBName, env.DBPort)
	if err != nil {
		return fmt.Errorf("failed to create database client: %w", err)
	}
	defer close()

	s := newServer(l)
	s.registerServices(db)
	return s.run(ctx)
}
