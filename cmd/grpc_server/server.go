package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/testing/testpb"
	handler "github.com/k-akari/otel-example/internal/handler/grpchandler"
	"github.com/k-akari/otel-example/internal/handler/grpchandler/interceptor"
	"github.com/k-akari/otel-example/internal/infra/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type server struct {
	srv *grpc.Server
	l   net.Listener
}

func newServer(l net.Listener) *server {
	opts := []grpc.ServerOption{
		interceptor.UnaryServerInterceptors(),
	}
	return &server{
		srv: grpc.NewServer(opts...),
		l:   l,
	}
}

func (s *server) registerServices(db *database.Client) {
	reflection.Register(s.srv)
	hs := health.NewServer()
	healthpb.RegisterHealthServer(s.srv, hs)

	testpb.RegisterTestServiceServer(s.srv, handler.NewTestService(db))
}

func (s *server) run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := s.srv.Serve(s.l); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to serve: %w", err)
		}
	case <-ctx.Done():
		s.srv.GracefulStop()
	}

	return nil
}
