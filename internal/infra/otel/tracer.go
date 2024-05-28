package otel

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func Init(ctx context.Context, serviceName string, conn *grpc.ClientConn) (func(), error) {
	te, err := newTraceExporter(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	me, err := newMeterExporter(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	res, err := newResource(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := newTracerProvider(te, res)
	mp := newMeterProvider(me, res)
	setPropagator()

	shutdown := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown TracerProvider: %v", err)
		}
		if err := mp.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown MeterProvider: %v", err)
		}
	}

	return shutdown, nil
}
