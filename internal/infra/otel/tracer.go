package otel

import (
	"context"
	"errors"
	"fmt"
	"log"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	samplingRate = 0.2
	moduleName   = "github.com/k-akari/otel-example"
)

func NewTracer(ctx context.Context, serviceName, gcpProjectID string) (trace.Tracer, func(), error) {
	var tp *sdktrace.TracerProvider
	if gcpProjectID != "" {
		exporter, err := texporter.New(texporter.WithProjectID(gcpProjectID))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}

		res, err := resource.New(ctx,
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
			return nil, nil, fmt.Errorf("failed to create resource: %w", err)
		}

		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithSampler(sdktrace.TraceIDRatioBased(samplingRate)),
			sdktrace.WithResource(res),
		)
	} else {
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)
	}
	otel.SetTracerProvider(tp)

	p := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(p)

	tr := otel.Tracer(moduleName)

	return tr, func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown TracerProvider: %v", err)
		}
	}, nil
}
