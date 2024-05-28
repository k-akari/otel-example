package otel

import (
	"context"
	"fmt"
	"log"
)

func Init(ctx context.Context, serviceName, endpointJaeger string) (func(), error) {
	exporter, err := newExporter(ctx, endpointJaeger)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	resource, err := newResource(ctx, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := newTracerProvider(exporter, resource)
	setPropagator()

	shutdown := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown TracerProvider: %v", err)
		}
	}

	return shutdown, nil
}
