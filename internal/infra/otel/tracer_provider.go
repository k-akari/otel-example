package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newTracerProvider(exporter *otlptrace.Exporter, resource *resource.Resource) *sdktrace.TracerProvider {
	const samplingRate = 0.2
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(samplingRate)),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	return tp
}
