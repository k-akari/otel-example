package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/k-akari/otel-example/internal/infra/database"
	"github.com/k-akari/otel-example/internal/infra/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	env := mustNewConfig()

	connCollector, err := grpc.NewClient(env.EndpointCollector, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	close, err := otel.Init(ctx, "grpc_server", connCollector)
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
