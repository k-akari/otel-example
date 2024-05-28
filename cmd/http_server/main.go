package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	"github.com/k-akari/otel-example/internal/handler/httphandler"
	"github.com/k-akari/otel-example/internal/infra/database"
	internal_otel "github.com/k-akari/otel-example/internal/infra/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
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

	tracer, close, err := internal_otel.NewTracer(ctx, "http_server", env.EndpointJaeger)
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

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	}
	conn, err := grpc.NewClient(env.EndpointGRPCServer, opts...)
	if err != nil {
		return fmt.Errorf("failed to create grpc client: %w", err)
	}
	defer conn.Close()

	tsc := testpb.NewTestServiceClient(conn)
	tsh := httphandler.NewTestService(db, tsc)

	mux := newMux(tracer, tsh)
	s := newServer(l, mux)
	return s.run(ctx)
}
