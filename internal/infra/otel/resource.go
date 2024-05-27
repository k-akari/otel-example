package otel

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func newResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	res, err := resource.New(
		ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithTelemetrySDK(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithContainer(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		log.Printf("failed to create resource: %v", err)
	} else if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return res, nil
}
