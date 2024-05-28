package otel

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewTracer(ctx context.Context, serviceName, endpointJaeger string) (trace.Tracer, func(), error) {
	exporter, err := newExporter(ctx, endpointJaeger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	resource, err := newResource(ctx, serviceName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := newTracerProvider(exporter, resource)
	setPropagator()

	const moduleName = "github.com/k-akari/otel-example"
	tracer := otel.Tracer(filepath.Join(moduleName, serviceName))
	shutdown := func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown TracerProvider: %v", err)
		}
	}

	return tracer, shutdown, nil
}
